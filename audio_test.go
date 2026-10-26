package audact_test

import (
	"encoding/json"
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestAudioService(t *testing.T) {
	yes := true
	tests := []svcCase{
		{name: "validate", status: 200, resp: `{"verdict":"ALLOW","segments":[{"segment_index":0,"parallel_asr_result":{"converged":true}}],"receipt":{"signature":"s"}}`,
			call: func(c *audact.Client) (any, error) {
				r, err := c.Audio.Validate(bg, audact.ValidateAudioRequest{
					CallID: "c1", AudioData: "AAA", AudioFormat: audact.AudioFormatWAV, ParallelASR: &yes,
					CallContext: &audact.AudioCallContext{Language: "de"}})
				if err == nil && (r.Receipt.Signature != "s" || !r.Segments[0].ParallelASRResult.Converged) {
					t.Errorf("resp = %+v", r)
				}
				return r, err
			},
			wantM: "POST", wantPath: "/v1/validate-audio", wantBody: `{"audio_data":"AAA","audio_format":"wav","call_context":{"language":"de"},"call_id":"c1","parallel_asr":true}`, wantAuth: "Bearer apikey"},
		{name: "reconcile with params", status: 200, resp: `{"call_id":"c1","has_more":true,"segment_results":[{"sequence_number":2}]}`,
			call: func(c *audact.Client) (any, error) {
				return c.Audio.Reconcile(bg, audact.AudioReconciliationRequest{CallID: "c1", AudioFormat: audact.AudioFormatMP3},
					&audact.ReconcileAudioParams{Limit: 50, AfterSequence: 1, IdempotencyKey: "idem"})
			},
			wantM: "POST", wantPath: "/v1/audio-reconciliation", wantQuery: "after_sequence=1&limit=50", wantBody: `{"audio_format":"mp3","call_id":"c1"}`, wantAuth: "Bearer apikey",
			wantHdr: map[string][]string{"Idempotency-Key": {"idem"}}},
		{name: "reconcile nil params", status: 200, resp: `{}`,
			call: func(c *audact.Client) (any, error) {
				return c.Audio.Reconcile(bg, audact.AudioReconciliationRequest{CallID: "c1"}, nil)
			},
			wantM: "POST", wantPath: "/v1/audio-reconciliation", wantBody: `{"call_id":"c1"}`, wantAuth: "Bearer apikey"},
	}
	runSvcCases(t, tests)
}

func TestAudioServiceNotFound(t *testing.T) {
	c := notFoundServer(t)
	if _, err := c.Audio.Validate(bg, audact.ValidateAudioRequest{}); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("validate err = %v", err)
	}
	if _, err := c.Audio.Reconcile(bg, audact.AudioReconciliationRequest{}, nil); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("reconcile err = %v", err)
	}
}

func TestAudioSegmentResult_TranscriptionStatus(t *testing.T) {
	tests := []struct {
		name string
		body string
		want audact.TranscriptionStatus
	}{
		{"transcribed", `{"transcription_status":"transcribed"}`, audact.TranscriptionStatusTranscribed},
		{"no speech", `{"transcription_status":"no_speech"}`, audact.TranscriptionStatusNoSpeech},
		{"absent", `{}`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r audact.AudioSegmentResult
			if err := json.Unmarshal([]byte(tt.body), &r); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if r.TranscriptionStatus != tt.want {
				t.Errorf("TranscriptionStatus = %q, want %q", r.TranscriptionStatus, tt.want)
			}
		})
	}
}
