package audact

import "io"

// AgentStatus is the lifecycle status of an agent.
type AgentStatus string

// Agent lifecycle statuses. Create accepts draft and published only.
const (
	AgentStatusDraft     AgentStatus = "draft"
	AgentStatusPublished AgentStatus = "published"
	AgentStatusArchived  AgentStatus = "archived"
)

// LLMProvider is a language-model provider.
type LLMProvider string

// Supported LLM providers.
const (
	LLMProviderMistral LLMProvider = "mistral"
	LLMProviderOpenAI  LLMProvider = "openai"
)

// LLMModel is a language-model name.
type LLMModel string

// Supported LLM models.
const (
	LLMModelGPT41Mini        LLMModel = "gpt-4.1-mini"
	LLMModelMistralSmall2603 LLMModel = "mistral-small-2603"
)

// STTProvider is a speech-to-text provider.
type STTProvider string

// STTProviderDeepgram is the Deepgram speech-to-text provider.
const STTProviderDeepgram STTProvider = "deepgram"

// TTSProvider is a text-to-speech provider.
type TTSProvider string

// Supported text-to-speech providers. Voxtral (Mistral, EU) is the default;
// ElevenLabs (US, global endpoint) is opt-in per agent.
const (
	TTSProviderElevenLabs TTSProvider = "elevenlabs"
	TTSProviderVoxtral    TTSProvider = "voxtral"
)

// TTSModel is a provider-specific text-to-speech model. The allowed values
// depend on the agent's TTSProvider; TTSModelDefault ("") selects the
// provider default.
type TTSModel string

// Supported TTS models.
const (
	TTSModelDefault                TTSModel = ""
	TTSModelElevenFlashV25         TTSModel = "eleven_flash_v2_5"
	TTSModelElevenMultilingualV2   TTSModel = "eleven_multilingual_v2"
	TTSModelElevenTurboV25         TTSModel = "eleven_turbo_v2_5"
	TTSModelElevenV3Conversational TTSModel = "eleven_v3_conversational"
	TTSModelElevenV4Turbo          TTSModel = "eleven_v4_turbo"
	TTSModelVoxtralMiniTTS2603     TTSModel = "voxtral-mini-tts-2603"
)

// StartSpeaker says who speaks first on a call.
type StartSpeaker string

// Start speaker values.
const (
	StartSpeakerAgent  StartSpeaker = "agent"
	StartSpeakerCaller StartSpeaker = "caller"
)

// VoicemailBranch is the action taken when voicemail is detected.
type VoicemailBranch string

// Voicemail branches.
const (
	VoicemailBranchLeaveMessage VoicemailBranch = "leave_message"
	VoicemailBranchHangup       VoicemailBranch = "hangup"
)

// VoicemailConfig configures voicemail detection and handling.
type VoicemailConfig struct {
	BeepHz           int               `json:"beep_hz,omitempty"`
	Branch           VoicemailBranch   `json:"branch,omitempty"`
	DecisionWindowMS int               `json:"decision_window_ms,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	Message          map[string]string `json:"message,omitempty"`
}

// EndCallSilenceConfig configures silence handling before ending a call.
type EndCallSilenceConfig struct {
	MaxReprompts int               `json:"max_reprompts,omitempty"`
	RepromptText map[string]string `json:"reprompt_text,omitempty"`
	WindowMS     int               `json:"window_ms,omitempty"`
}

// EndCallConfig configures how the agent ends calls.
type EndCallConfig struct {
	ClosingLine    map[string]string     `json:"closing_line,omitempty"`
	Silence        *EndCallSilenceConfig `json:"silence,omitempty"`
	TriggerPhrases map[string][]string   `json:"trigger_phrases,omitempty"`
}

// ExtractionField describes one field of post-call extraction.
type ExtractionField struct {
	Description string   `json:"description,omitempty"`
	Format      string   `json:"format,omitempty"`
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type,omitempty"`
	Values      []string `json:"values,omitempty"`
}

// PostCallExtractionConfig configures post-call data extraction.
type PostCallExtractionConfig struct {
	Enabled *bool             `json:"enabled,omitempty"`
	Fields  []ExtractionField `json:"fields,omitempty"`
}

// PostCallSummaryConfig configures post-call summaries.
type PostCallSummaryConfig struct {
	Dispositions []string `json:"dispositions,omitempty"`
	Enabled      *bool    `json:"enabled,omitempty"`
}

// Agent is a voice agent.
type Agent struct {
	AgentID            string                    `json:"agent_id,omitempty"`
	AssignedPhone      string                    `json:"assigned_phone,omitempty"`
	AutoHandoverConfig map[string]any            `json:"auto_handover_config,omitempty"`
	BeginMessage       string                    `json:"begin_message,omitempty"`
	ConfigurationID    string                    `json:"configuration_id,omitempty"`
	CreatedAt          string                    `json:"created_at,omitempty"`
	DeletedAt          string                    `json:"deleted_at,omitempty"`
	EndCall            *EndCallConfig            `json:"end_call,omitempty"`
	GreetOnConnect     bool                      `json:"greet_on_connect,omitempty"`
	Jurisdiction       string                    `json:"jurisdiction,omitempty"`
	Language           string                    `json:"language,omitempty"`
	LLMModel           string                    `json:"llm_model,omitempty"`
	LLMProvider        string                    `json:"llm_provider,omitempty"`
	MaxDurationSeconds int                       `json:"max_duration_seconds,omitempty"`
	Name               string                    `json:"name,omitempty"`
	PostCallExtraction *PostCallExtractionConfig `json:"post_call_extraction,omitempty"`
	PostCallSummary    *PostCallSummaryConfig    `json:"post_call_summary,omitempty"`
	Prompt             string                    `json:"prompt,omitempty"`
	StartSpeaker       StartSpeaker              `json:"start_speaker,omitempty"`
	Status             string                    `json:"status,omitempty"`
	STTProvider        string                    `json:"stt_provider,omitempty"`
	Temperature        float64                   `json:"temperature,omitempty"`
	TenantID           string                    `json:"tenant_id,omitempty"`
	ToolIDs            []string                  `json:"tool_ids,omitempty"`
	TTSModel           string                    `json:"tts_model,omitempty"`
	TTSProvider        string                    `json:"tts_provider,omitempty"`
	UpdatedAt          string                    `json:"updated_at,omitempty"`
	Version            int                       `json:"version,omitempty"`
	VoiceID            string                    `json:"voice_id,omitempty"`
	Voicemail          *VoicemailConfig          `json:"voicemail,omitempty"`
}

// ListAgentsParams are the optional query parameters of Agents.List.
type ListAgentsParams struct {
	Limit           int
	Cursor          string
	ConfigurationID string
}

// ListAgentsResponse is a page of agents.
type ListAgentsResponse struct {
	Agents     []Agent `json:"agents,omitempty"`
	HasMore    bool    `json:"has_more,omitempty"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

// CreateAgentRequest is the body of Agents.Create. Unset fields are not sent.
type CreateAgentRequest struct {
	AutoHandoverConfig map[string]any            `json:"auto_handover_config,omitempty"`
	BeginMessage       string                    `json:"begin_message,omitempty"`
	EndCall            *EndCallConfig            `json:"end_call,omitempty"`
	GreetOnConnect     *bool                     `json:"greet_on_connect,omitempty"`
	Jurisdiction       string                    `json:"jurisdiction,omitempty"`
	Language           string                    `json:"language,omitempty"`
	LLMModel           LLMModel                  `json:"llm_model,omitempty"`
	LLMProvider        LLMProvider               `json:"llm_provider,omitempty"`
	MaxDurationSeconds *int                      `json:"max_duration_seconds,omitempty"`
	Name               string                    `json:"name,omitempty"`
	PostCallExtraction *PostCallExtractionConfig `json:"post_call_extraction,omitempty"`
	PostCallSummary    *PostCallSummaryConfig    `json:"post_call_summary,omitempty"`
	Prompt             string                    `json:"prompt,omitempty"`
	StartSpeaker       StartSpeaker              `json:"start_speaker,omitempty"`
	Status             AgentStatus               `json:"status,omitempty"`
	STTProvider        STTProvider               `json:"stt_provider,omitempty"`
	Temperature        *float64                  `json:"temperature,omitempty"`
	ToolIDs            []string                  `json:"tool_ids,omitempty"`
	TTSModel           TTSModel                  `json:"tts_model,omitempty"`
	TTSProvider        TTSProvider               `json:"tts_provider,omitempty"`
	VoiceID            string                    `json:"voice_id,omitempty"`
	Voicemail          *VoicemailConfig          `json:"voicemail,omitempty"`
}

// UpdateAgentRequest is the body of Agents.Update (PUT). It has the same
// fields as CreateAgentRequest; status may also be "archived".
type UpdateAgentRequest CreateAgentRequest

// PatchAgentRequest is the body of Agents.Patch (PATCH). Only the fields that
// are set are sent.
//
// Because unset means "omitted", a zero value can never be sent: Patch cannot
// clear a field. In particular it cannot set TTSModel back to TTSModelDefault
// ("") while keeping TTSProvider. To do that, use Agents.Update (a full PUT, where
// an omitted tts_model is the provider default), or change TTSProvider, which
// resets the model server-side.
type PatchAgentRequest CreateAgentRequest

// CreateAgentResponse is the result of Agents.Create.
type CreateAgentResponse struct {
	AgentID         string   `json:"agent_id,omitempty"`
	ConfigurationID string   `json:"configuration_id,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	Jurisdiction    string   `json:"jurisdiction,omitempty"`
	Language        string   `json:"language,omitempty"`
	LLMModel        string   `json:"llm_model,omitempty"`
	LLMProvider     string   `json:"llm_provider,omitempty"`
	Name            string   `json:"name,omitempty"`
	Prompt          string   `json:"prompt,omitempty"`
	Status          string   `json:"status,omitempty"`
	STTProvider     string   `json:"stt_provider,omitempty"`
	Temperature     float64  `json:"temperature,omitempty"`
	TenantID        string   `json:"tenant_id,omitempty"`
	ToolIDs         []string `json:"tool_ids,omitempty"`
	TTSModel        string   `json:"tts_model,omitempty"`
	TTSProvider     string   `json:"tts_provider,omitempty"`
	VoiceID         string   `json:"voice_id,omitempty"`
}

// DeleteAgentResponse is the result of Agents.Delete.
type DeleteAgentResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

// Voice is a TTS voice.
type Voice struct {
	Accent   string `json:"accent,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Language string `json:"language,omitempty"`
	// Languages are the language codes the voice can speak (Language is its
	// native one). ElevenLabs voices are multilingual.
	Languages  []string `json:"languages,omitempty"`
	Name       string   `json:"name,omitempty"`
	PreviewURL string   `json:"preview_url,omitempty"`
	Style      string   `json:"style,omitempty"`
	VoiceID    string   `json:"voice_id,omitempty"`
}

// TTSModelInfo is one TTS model of a provider and the languages it speaks.
type TTSModelInfo struct {
	// ID is the tts_model value ("" on an agent selects the default).
	ID      string `json:"id"`
	Default bool   `json:"default"`
	// Languages the model speaks; empty when not published (Voxtral).
	Languages []string `json:"languages,omitempty"`
}

// ListVoicesParams filters Agents.ListVoicesFor.
type ListVoicesParams struct {
	// Language keeps voices (and models) that speak it, e.g. "de" or "de-DE".
	Language string
	// Provider keeps one TTS provider.
	Provider TTSProvider
	// Limit is the maximum voices per provider, 1-1000 (0 = server default 100).
	Limit int
}

// VoiceProvider is a TTS provider and its voices.
type VoiceProvider struct {
	DisplayName string         `json:"display_name,omitempty"`
	Models      []TTSModelInfo `json:"models,omitempty"`
	Voices      []Voice        `json:"voices,omitempty"`
	// Total is how many voices matched before the per-provider limit.
	Total int `json:"total,omitempty"`
	// Truncated is true when Voices was cut to the limit; narrow with Language
	// or raise Limit.
	Truncated bool `json:"truncated,omitempty"`
}

// ListVoicesResponse maps provider name to its voices.
type ListVoicesResponse map[string]VoiceProvider

// CreateVoiceRequest is the input of Agents.CreateVoice, sent as
// multipart/form-data. Audio is a ~10s reference clip (audio/mpeg or
// audio/wav, max 2 MiB). ConsentConfirmed must be true.
type CreateVoiceRequest struct {
	Name             string
	Audio            io.Reader
	Filename         string
	ConsentConfirmed bool
	Language         string
	Accent           string
	Gender           string
	Style            string
}

// CreateVoiceResponse is the result of Agents.CreateVoice.
type CreateVoiceResponse struct {
	Accent     string `json:"accent,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Language   string `json:"language,omitempty"`
	Name       string `json:"name,omitempty"`
	PreviewURL string `json:"preview_url,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Style      string `json:"style,omitempty"`
	VoiceID    string `json:"voice_id,omitempty"`
}

// PipelineStatusResponse is the result of Agents.PipelineStatus.
type PipelineStatusResponse struct {
	ActivePipelines int    `json:"active_pipelines,omitempty"`
	Timestamp       string `json:"timestamp,omitempty"`
}
