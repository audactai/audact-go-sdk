package audact

// ToolHeader is a static header sent with a tool's HTTP request.
type ToolHeader struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// ToolQueryParam is a query parameter of a tool; Mode is "const" or "description".
type ToolQueryParam struct {
	Key   string `json:"key,omitempty"`
	Mode  string `json:"mode,omitempty"`
	Value string `json:"value,omitempty"`
}

// ToolResponseVar extracts a dynamic variable from a tool response.
type ToolResponseVar struct {
	JSONPath string `json:"json_path,omitempty"`
	Name     string `json:"name,omitempty"`
}

// ToolSpeakDuring controls hold audio while a tool executes.
type ToolSpeakDuring struct {
	Enabled bool   `json:"enabled,omitempty"`
	Text    string `json:"text,omitempty"`
}

// ToolExecution holds runtime dispatch settings of a tool.
type ToolExecution struct {
	ModelOverride string           `json:"model_override,omitempty"`
	SpeakDuring   *ToolSpeakDuring `json:"speak_during,omitempty"`
	TimeoutMs     int              `json:"timeout_ms,omitempty"`
}

// Tool is a tool configuration.
type Tool struct {
	ArgsOnly        bool              `json:"args_only,omitempty"`
	ConnectionID    string            `json:"connection_id,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	DeletedAt       string            `json:"deleted_at,omitempty"`
	Description     string            `json:"description,omitempty"`
	DestinationHint string            `json:"destination_hint,omitempty"`
	Execution       *ToolExecution    `json:"execution,omitempty"`
	Headers         []ToolHeader      `json:"headers,omitempty"`
	Method          string            `json:"method,omitempty"`
	Name            string            `json:"name,omitempty"`
	Parameters      map[string]any    `json:"parameters,omitempty"`
	Provider        string            `json:"provider,omitempty"`
	ProviderConfig  map[string]any    `json:"provider_config,omitempty"`
	ProviderTool    string            `json:"provider_tool,omitempty"`
	QueryParams     []ToolQueryParam  `json:"query_params,omitempty"`
	ResponseVars    []ToolResponseVar `json:"response_vars,omitempty"`
	TenantID        string            `json:"tenant_id,omitempty"`
	ToolID          string            `json:"tool_id,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty"`
	URL             string            `json:"url,omitempty"`
}

// CreateToolRequest is the body of Tools.Create and Tools.Replace.
type CreateToolRequest struct {
	ArgsOnly       bool              `json:"args_only,omitempty"`
	ConnectionID   string            `json:"connection_id,omitempty"`
	Description    string            `json:"description,omitempty"`
	Execution      *ToolExecution    `json:"execution,omitempty"`
	Headers        []ToolHeader      `json:"headers,omitempty"`
	Method         string            `json:"method,omitempty"`
	Name           string            `json:"name,omitempty"`
	Parameters     map[string]any    `json:"parameters,omitempty"`
	Provider       string            `json:"provider,omitempty"`
	ProviderConfig map[string]any    `json:"provider_config,omitempty"`
	ProviderTool   string            `json:"provider_tool,omitempty"`
	QueryParams    []ToolQueryParam  `json:"query_params,omitempty"`
	ResponseVars   []ToolResponseVar `json:"response_vars,omitempty"`
	URL            string            `json:"url,omitempty"`
}

// UpdateToolRequest is the body of Tools.Update; unset fields are not sent.
type UpdateToolRequest struct {
	ArgsOnly       *bool             `json:"args_only,omitempty"`
	ConnectionID   string            `json:"connection_id,omitempty"`
	Description    string            `json:"description,omitempty"`
	Execution      *ToolExecution    `json:"execution,omitempty"`
	Headers        []ToolHeader      `json:"headers,omitempty"`
	Method         string            `json:"method,omitempty"`
	Name           string            `json:"name,omitempty"`
	Parameters     map[string]any    `json:"parameters,omitempty"`
	Provider       string            `json:"provider,omitempty"`
	ProviderConfig map[string]any    `json:"provider_config,omitempty"`
	ProviderTool   string            `json:"provider_tool,omitempty"`
	QueryParams    []ToolQueryParam  `json:"query_params,omitempty"`
	ResponseVars   []ToolResponseVar `json:"response_vars,omitempty"`
	URL            string            `json:"url,omitempty"`
}

// ListToolsParams holds optional query parameters for Tools.List.
type ListToolsParams struct {
	Name   string
	Cursor string
	Limit  int
}

// ListToolsResponse is a page of tools.
type ListToolsResponse struct {
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more,omitempty"`
	Tools   []Tool `json:"tools,omitempty"`
}

// DeleteToolResponse reports a tool deletion.
type DeleteToolResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

// TestToolResponse is the result of a tool health check.
type TestToolResponse struct {
	ChecksPassed     []string `json:"checks_passed,omitempty"`
	DynamicVarFields []string `json:"dynamic_var_fields,omitempty"`
	Error            string   `json:"error,omitempty"`
	Healthy          bool     `json:"healthy,omitempty"`
	ToolID           string   `json:"tool_id,omitempty"`
}

// ListAgentToolsParams holds optional query parameters for Tools.ListForAgent.
type ListAgentToolsParams struct {
	Cursor string
	Limit  int
}

// ListAgentToolsResponse is a page of tools attached to an agent.
type ListAgentToolsResponse struct {
	HasMore    bool   `json:"has_more,omitempty"`
	Items      []Tool `json:"items,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// AttachToolRequest is the body of Tools.AttachToAgent.
type AttachToolRequest struct {
	ToolID string `json:"tool_id"`
}

// AttachToolResponse reports a tool attachment.
type AttachToolResponse struct {
	AgentID  string `json:"agent_id,omitempty"`
	Attached bool   `json:"attached,omitempty"`
	ToolID   string `json:"tool_id,omitempty"`
}

// DetachToolResponse reports a tool detachment.
type DetachToolResponse struct {
	Detached bool `json:"detached,omitempty"`
}
