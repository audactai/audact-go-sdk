package audact

import "encoding/json"

// ToolCallFunction is the function name and arguments of a tool call.
type ToolCallFunction struct {
	Arguments string `json:"arguments,omitempty"`
	Name      string `json:"name,omitempty"`
}

// ToolCall is a tool call made by the model.
type ToolCall struct {
	Function ToolCallFunction `json:"function"`
	ID       string           `json:"id,omitempty"`
	Index    int              `json:"index,omitempty"`
	Type     string           `json:"type,omitempty"`
}

// ChatMessage is one message of a chat conversation.
type ChatMessage struct {
	Content    string     `json:"content,omitempty"`
	Name       string     `json:"name,omitempty"`
	Role       string     `json:"role,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// ToolFunctionSchema is the schema of a callable function.
type ToolFunctionSchema struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	// Parameters is the JSON Schema of the function arguments.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// ToolDefinition is a function tool offered to the model.
type ToolDefinition struct {
	Function ToolFunctionSchema `json:"function"`
	Type     string             `json:"type,omitempty"`
}

// ChatCompletionRequest is an OpenAI-compatible chat completion request.
type ChatCompletionRequest struct {
	MaxCompletionTokens *int             `json:"max_completion_tokens,omitempty"`
	MaxTokens           *int             `json:"max_tokens,omitempty"`
	Messages            []ChatMessage    `json:"messages"`
	Model               string           `json:"model,omitempty"`
	ParallelToolCalls   *bool            `json:"parallel_tool_calls,omitempty"`
	ReasoningEffort     string           `json:"reasoning_effort,omitempty"`
	Stream              *bool            `json:"stream,omitempty"`
	Temperature         *float64         `json:"temperature,omitempty"`
	ToolChoice          json.RawMessage  `json:"tool_choice,omitempty"`
	Tools               []ToolDefinition `json:"tools,omitempty"`
}

// ProxyParams holds the optional query parameters and call-context headers
// of an LLM proxy request. Empty fields are not sent.
type ProxyParams struct {
	// ConfigurationID selects a tenant configuration (query configuration_id).
	ConfigurationID string
	// CallID is sent as X-Audact-Call-ID.
	CallID string
	// TenantID is sent as X-Audact-Tenant-ID.
	TenantID string
	// ConfigID is sent as X-Audact-Config-ID.
	ConfigID string
	// Turn is sent as X-Audact-Turn when non-zero.
	Turn int
}

// UsageStats reports token consumption.
type UsageStats struct {
	CompletionTokens int `json:"completion_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

// AudactValidationMeta is the compliance metadata injected by the proxy.
type AudactValidationMeta struct {
	CallID   string `json:"call_id,omitempty"`
	FailOpen bool   `json:"fail_open,omitempty"`
	Modified bool   `json:"modified,omitempty"`
	Verdict  string `json:"verdict,omitempty"`
}

// ChunkDelta is the incremental content of one streamed choice.
type ChunkDelta struct {
	Content   string     `json:"content,omitempty"`
	Role      string     `json:"role,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChunkChoice is one choice delta of a streamed chunk.
type ChunkChoice struct {
	Delta        ChunkDelta `json:"delta"`
	FinishReason string     `json:"finish_reason,omitempty"`
	Index        int        `json:"index,omitempty"`
}

// ChatCompletionChunk is one event of a streaming chat completion
// (object "chat.completion.chunk").
type ChatCompletionChunk struct {
	AudactValidation *AudactValidationMeta `json:"audact_validation,omitempty"`
	Choices          []ChunkChoice         `json:"choices,omitempty"`
	Created          int64                 `json:"created,omitempty"`
	ID               string                `json:"id,omitempty"`
	Model            string                `json:"model,omitempty"`
	Object           string                `json:"object,omitempty"`
}

// ChatCompletionChoice is one completion choice.
type ChatCompletionChoice struct {
	FinishReason string      `json:"finish_reason,omitempty"`
	Index        int         `json:"index,omitempty"`
	Message      ChatMessage `json:"message"`
}

// ChatCompletionResponse is an OpenAI-compatible chat completion response.
type ChatCompletionResponse struct {
	AudactValidation *AudactValidationMeta  `json:"audact_validation,omitempty"`
	Choices          []ChatCompletionChoice `json:"choices,omitempty"`
	Created          int64                  `json:"created,omitempty"`
	ID               string                 `json:"id,omitempty"`
	Model            string                 `json:"model,omitempty"`
	Object           string                 `json:"object,omitempty"`
	Usage            *UsageStats            `json:"usage,omitempty"`
}
