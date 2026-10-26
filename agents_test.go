package audact_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

type seen struct {
	method, path, rawQuery, auth, ctype string
	body                                []byte
}

func newAgentsServer(t *testing.T, status int, resp string) (*audact.Client, *seen) {
	t.Helper()
	var got seen
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = seen{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Authorization"), r.Header.Get("Content-Type"), b}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(resp))
	}))
	t.Cleanup(srv.Close)
	c, err := audact.NewClient("key123", audact.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return c, &got
}

func TestAgentsService(t *testing.T) {
	temp := 0.5
	f := false
	tests := []struct {
		name      string
		status    int
		resp      string
		call      func(c *audact.Client) (any, error)
		wantM     string
		wantPath  string
		wantQuery string
		wantBody  string // JSON, "" = no body
		wantOut   string // JSON of result, "" = none
	}{
		{"list nil params", 200, `{"agents":[{"agent_id":"a1","name":"N"}],"has_more":true,"next_cursor":"c2"}`,
			func(c *audact.Client) (any, error) { return c.Agents.List(context.Background(), nil) },
			"GET", "/v1/agents", "", "", `{"agents":[{"agent_id":"a1","name":"N"}],"has_more":true,"next_cursor":"c2"}`},
		{"list params", 200, `{}`,
			func(c *audact.Client) (any, error) {
				return c.Agents.List(context.Background(), &audact.ListAgentsParams{Limit: 5, Cursor: "c1", ConfigurationID: "cfg"})
			},
			"GET", "/v1/agents", "configuration_id=cfg&cursor=c1&limit=5", "", `{}`},
		{"create", 201, `{"agent_id":"a1"}`,
			func(c *audact.Client) (any, error) {
				return c.Agents.Create(context.Background(), &audact.CreateAgentRequest{Name: "N", Status: audact.AgentStatusDraft, Temperature: &temp, GreetOnConnect: &f})
			},
			"POST", "/v1/agents", "", `{"name":"N","status":"draft","temperature":0.5,"greet_on_connect":false}`, `{"agent_id":"a1"}`},
		{"create nil", 201, `{}`,
			func(c *audact.Client) (any, error) { return c.Agents.Create(context.Background(), nil) },
			"POST", "/v1/agents", "", `{}`, `{}`},
		{"get escapes", 200, `{"agent_id":"a/1"}`,
			func(c *audact.Client) (any, error) { return c.Agents.Get(context.Background(), "a/1") },
			"GET", "/v1/agents/a%2F1", "", "", `{"agent_id":"a/1"}`},
		{"update", 200, `{"agent_id":"a1"}`,
			func(c *audact.Client) (any, error) {
				return c.Agents.Update(context.Background(), "a1", &audact.UpdateAgentRequest{Status: audact.AgentStatusArchived})
			},
			"PUT", "/v1/agents/a1", "", `{"status":"archived"}`, `{"agent_id":"a1"}`},
		{"patch sends only set fields", 200, `{"agent_id":"a1"}`,
			func(c *audact.Client) (any, error) {
				return c.Agents.Patch(context.Background(), "a1", &audact.PatchAgentRequest{Prompt: "hi", GreetOnConnect: &f})
			},
			"PATCH", "/v1/agents/a1", "", `{"prompt":"hi","greet_on_connect":false}`, `{"agent_id":"a1"}`},
		{"delete", 200, `{"deleted":true}`,
			func(c *audact.Client) (any, error) { return c.Agents.Delete(context.Background(), "a1") },
			"DELETE", "/v1/agents/a1", "", "", `{"deleted":true}`},
		{"duplicate", 201, `{"agent_id":"a2"}`,
			func(c *audact.Client) (any, error) { return c.Agents.Duplicate(context.Background(), "a1") },
			"POST", "/v1/agents/a1/duplicate", "", "", `{"agent_id":"a2"}`},
		{"list voices", 200, `{"voxtral":{"display_name":"Voxtral","voices":[{"voice_id":"v1"}]}}`,
			func(c *audact.Client) (any, error) { return c.Agents.ListVoices(context.Background()) },
			"GET", "/v1/voices", "", "", `{"voxtral":{"display_name":"Voxtral","voices":[{"voice_id":"v1"}]}}`},
		{"get voice", 200, `{"voice_id":"v 1"}`,
			func(c *audact.Client) (any, error) { return c.Agents.GetVoice(context.Background(), "v 1") },
			"GET", "/v1/voices/v%201", "", "", `{"voice_id":"v 1"}`},
		{"delete voice 204", 204, ``,
			func(c *audact.Client) (any, error) { return nil, c.Agents.DeleteVoice(context.Background(), "v1") },
			"DELETE", "/v1/voices/v1", "", "", ""},
		{"pipeline status", 200, `{"active_pipelines":3,"timestamp":"t"}`,
			func(c *audact.Client) (any, error) { return c.Agents.PipelineStatus(context.Background()) },
			"GET", "/v1/pipeline/status", "", "", `{"active_pipelines":3,"timestamp":"t"}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, got := newAgentsServer(t, tc.status, tc.resp)
			out, err := tc.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if got.method != tc.wantM || got.path != tc.wantPath || got.rawQuery != tc.wantQuery {
				t.Errorf("request = %s %s ?%s", got.method, got.path, got.rawQuery)
			}
			if got.auth != "Bearer key123" {
				t.Errorf("auth = %q", got.auth)
			}
			if tc.wantBody == "" {
				if len(got.body) != 0 {
					t.Errorf("unexpected body %s", got.body)
				}
			} else {
				if got.ctype != "application/json" {
					t.Errorf("content-type = %q", got.ctype)
				}
				if !jsonEqual(t, got.body, tc.wantBody) {
					t.Errorf("body = %s, want %s", got.body, tc.wantBody)
				}
			}
			if tc.wantOut != "" {
				b, _ := json.Marshal(out)
				if !jsonEqual(t, b, tc.wantOut) {
					t.Errorf("out = %s, want %s", b, tc.wantOut)
				}
			}
		})
	}
}

func jsonEqual(t *testing.T, a []byte, b string) bool {
	t.Helper()
	var x, y any
	if err := json.Unmarshal(a, &x); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &y); err != nil {
		t.Fatal(err)
	}
	xb, _ := json.Marshal(x)
	yb, _ := json.Marshal(y)
	return string(xb) == string(yb)
}

func TestAgentsPatchOmitsUnsetFields(t *testing.T) {
	c, got := newAgentsServer(t, 200, `{}`)
	if _, err := c.Agents.Patch(context.Background(), "a1", &audact.PatchAgentRequest{Name: "x"}); err != nil {
		t.Fatal(err)
	}
	if string(got.body) != `{"name":"x"}` {
		t.Errorf("body = %s", got.body)
	}
}

func TestAgentsCreateVoice(t *testing.T) {
	c, got := newAgentsServer(t, 201, `{"voice_id":"v9","provider":"voxtral"}`)
	out, err := c.Agents.CreateVoice(context.Background(), &audact.CreateVoiceRequest{
		Name: "Me", Audio: strings.NewReader("RIFFdata"), Filename: "clip.wav", ConsentConfirmed: true, Language: "en",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.VoiceID != "v9" || got.method != "POST" || got.path != "/v1/voices" || got.auth != "Bearer key123" {
		t.Fatalf("got %+v out %+v", got, out)
	}
	if !strings.HasPrefix(got.ctype, "multipart/form-data; boundary=") {
		t.Fatalf("content-type = %q", got.ctype)
	}
	r, _ := http.NewRequest("POST", "/", strings.NewReader(string(got.body)))
	r.Header.Set("Content-Type", got.ctype)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"name": "Me", "consent_confirmed": "true", "language": "en"}
	for k, v := range want {
		if r.FormValue(k) != v {
			t.Errorf("field %s = %q", k, r.FormValue(k))
		}
	}
	if _, ok := r.MultipartForm.Value["accent"]; ok {
		t.Error("accent should be omitted when empty")
	}
	fh := r.MultipartForm.File["audio"]
	if len(fh) != 1 || fh[0].Filename != "clip.wav" {
		t.Fatalf("audio part = %+v", fh)
	}
	f, _ := fh[0].Open()
	b, _ := io.ReadAll(f)
	if string(b) != "RIFFdata" {
		t.Errorf("audio = %q", b)
	}
}

func TestAgentsCreateVoiceValidation(t *testing.T) {
	c, _ := newAgentsServer(t, 201, `{}`)
	for name, req := range map[string]*audact.CreateVoiceRequest{"nil": nil, "no audio": {Name: "x"}} {
		t.Run(name, func(t *testing.T) {
			if _, err := c.Agents.CreateVoice(context.Background(), req); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestAgentsErrors(t *testing.T) {
	c, _ := newAgentsServer(t, 404, `{"error":{"code":"not_found","message":"agent not found","request_id":"req_123"}}`)
	_, err := c.Agents.Get(context.Background(), "nope")
	if !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	var apiErr *audact.APIError
	if !errors.As(err, &apiErr) || apiErr.RequestID != "req_123" {
		t.Fatalf("apiErr = %+v", apiErr)
	}
	if err := c.Agents.DeleteVoice(context.Background(), "v"); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("delete voice err = %v", err)
	}
}

// TestCreateAgent_ElevenLabs: tts_model is sent when set, omitted when
// empty (so the server applies the provider default), and parsed back.
func TestCreateAgent_ElevenLabs(t *testing.T) {
	c, got := newAgentsServer(t, 201, `{"agent_id":"a1","tts_provider":"elevenlabs","tts_model":"eleven_turbo_v2_5","voice_id":"EXAVITQu4vr4xnSDxMaL"}`)
	resp, err := c.Agents.Create(t.Context(), &audact.CreateAgentRequest{
		Name:        "Concierge",
		TTSProvider: audact.TTSProviderElevenLabs,
		TTSModel:    audact.TTSModelElevenTurboV25,
		VoiceID:     "EXAVITQu4vr4xnSDxMaL",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(got.body, &body); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	if body["tts_model"] != "eleven_turbo_v2_5" || body["tts_provider"] != "elevenlabs" {
		t.Fatalf("request body = %s", got.body)
	}
	if resp.TTSModel != "eleven_turbo_v2_5" {
		t.Fatalf("response TTSModel = %q", resp.TTSModel)
	}

	// Default model: the field is omitted, never sent as "".
	c2, got2 := newAgentsServer(t, 201, `{"agent_id":"a2"}`)
	if _, err := c2.Agents.Create(t.Context(), &audact.CreateAgentRequest{Name: "Default", TTSProvider: audact.TTSProviderElevenLabs, TTSModel: audact.TTSModelDefault, VoiceID: "EXAVITQu4vr4xnSDxMaL"}); err != nil {
		t.Fatal(err)
	}
	var body2 map[string]any
	_ = json.Unmarshal(got2.body, &body2)
	if _, present := body2["tts_model"]; present {
		t.Fatalf("tts_model must be omitted for the provider default; body = %s", got2.body)
	}
}

// TestListVoicesFor_LanguageAndProvider: filters reach the query; languages and
// models are parsed.
func TestListVoicesFor_LanguageAndProvider(t *testing.T) {
	c, got := newAgentsServer(t, 200, `{"elevenlabs":{"display_name":"ElevenLabs","models":[{"id":"eleven_flash_v2_5","default":true,"languages":["de","en"]}],"voices":[{"voice_id":"DeDeDeDeDeDeDeDeDeDe","name":"Lena","language":"de","languages":["de","en"]}]}}`)
	out, err := c.Agents.ListVoicesFor(t.Context(), &audact.ListVoicesParams{Language: "de-DE", Provider: audact.TTSProviderElevenLabs})
	if err != nil {
		t.Fatal(err)
	}
	if got.rawQuery != "language=de-DE&provider=elevenlabs" {
		t.Fatalf("query = %q", got.rawQuery)
	}
	el := out["elevenlabs"]
	if len(el.Models) != 1 || !el.Models[0].Default || len(el.Voices) != 1 || el.Voices[0].Languages[0] != "de" {
		t.Fatalf("parsed = %+v", el)
	}
	// ListVoices (no filters) sends no query.
	c2, got2 := newAgentsServer(t, 200, `{}`)
	if _, err := c2.Agents.ListVoices(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got2.rawQuery != "" {
		t.Fatalf("unfiltered query = %q", got2.rawQuery)
	}
}
