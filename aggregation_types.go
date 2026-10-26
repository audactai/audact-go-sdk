package audact

// ListProofsParams are the query parameters of Aggregation.ListProofs.
type ListProofsParams struct {
	Limit  int
	Cursor string
}

// MetachainParams are the query parameters of Aggregation.Metachain.
type MetachainParams struct {
	Limit  int
	Cursor string
	// Order is the sort order (the order query parameter).
	Order string
}

// AdminAggregateRequest selects the window to aggregate.
type AdminAggregateRequest struct {
	WindowEnd   string `json:"window_end,omitempty"`
	WindowStart string `json:"window_start,omitempty"`
}

// AdminAggregateResponse is the result of a triggered aggregation.
type AdminAggregateResponse struct {
	AggregationID string `json:"aggregation_id,omitempty"`
	BatchNumber   int    `json:"batch_number,omitempty"`
	LeafCount     int    `json:"leaf_count,omitempty"`
	MerkleRoot    string `json:"merkle_root,omitempty"`
	TenantCount   int    `json:"tenant_count,omitempty"`
	TSTimestamp   string `json:"ts_timestamp,omitempty"`
	TSAName       string `json:"tsa_name,omitempty"`
	WindowEnd     string `json:"window_end,omitempty"`
	WindowStart   string `json:"window_start,omitempty"`
}

// MetaChainEntry is one entry of the aggregation meta-chain.
type MetaChainEntry struct {
	AggregationID string `json:"aggregation_id,omitempty"`
	ChainHash     string `json:"chain_hash,omitempty"`
	MerkleRoot    string `json:"merkle_root,omitempty"`
	MetaID        string `json:"meta_id,omitempty"`
	PreviousHash  string `json:"previous_hash,omitempty"`
	Sequence      int    `json:"sequence,omitempty"`
}

// MetaChainResponse is a page of the meta-chain.
type MetaChainResponse struct {
	Count      int              `json:"count,omitempty"`
	Entries    []MetaChainEntry `json:"entries,omitempty"`
	HasMore    bool             `json:"has_more,omitempty"`
	NextCursor string           `json:"next_cursor,omitempty"`
	Total      int              `json:"total,omitempty"`
}

// ProofSummary summarises one aggregation window proof.
type ProofSummary struct {
	AggregationID       string `json:"aggregation_id,omitempty"`
	BatchNumber         int    `json:"batch_number,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	MerkleRoot          string `json:"merkle_root,omitempty"`
	Status              string `json:"status,omitempty"`
	TenantAggregateHash string `json:"tenant_aggregate_hash,omitempty"`
	TimestampTime       string `json:"timestamp_time,omitempty"`
	TSAURL              string `json:"tsa_url,omitempty"`
	WindowEnd           string `json:"window_end,omitempty"`
	WindowID            string `json:"window_id,omitempty"`
	WindowStart         string `json:"window_start,omitempty"`
}

// ProofListResponse is a page of proof summaries.
type ProofListResponse struct {
	Count      int            `json:"count,omitempty"`
	HasMore    bool           `json:"has_more,omitempty"`
	NextCursor string         `json:"next_cursor,omitempty"`
	Proofs     []ProofSummary `json:"proofs,omitempty"`
	Total      int            `json:"total,omitempty"`
}

// ProofStep is one sibling hash on a Merkle authentication path.
type ProofStep struct {
	// Hash is the hex-encoded SHA-256 sibling hash.
	Hash string `json:"hash,omitempty"`
	// Position is "left" or "right".
	Position string `json:"position,omitempty"`
}

// BranchProof is the Merkle authentication path of a tenant aggregate.
type BranchProof struct {
	LeafCount int         `json:"leaf_count,omitempty"`
	LeafHash  string      `json:"leaf_hash,omitempty"`
	LeafIndex int         `json:"leaf_index,omitempty"`
	PadCount  int         `json:"pad_count,omitempty"`
	Root      string      `json:"root,omitempty"`
	Steps     []ProofStep `json:"steps,omitempty"`
	TreeDepth int         `json:"tree_depth,omitempty"`
}

// ProofResponse is the full proof of one aggregation window.
type ProofResponse struct {
	AggregationID       string       `json:"aggregation_id,omitempty"`
	BatchNumber         int          `json:"batch_number,omitempty"`
	BranchProof         *BranchProof `json:"branch_proof,omitempty"`
	MerkleRoot          string       `json:"merkle_root,omitempty"`
	TenantAggregateHash string       `json:"tenant_aggregate_hash,omitempty"`
	TimestampTime       string       `json:"timestamp_time,omitempty"`
	TimestampToken      string       `json:"timestamp_token,omitempty"`
	TSAURL              string       `json:"tsa_url,omitempty"`
	WindowID            string       `json:"window_id,omitempty"`
}

// VerifyInclusionRequest carries a Merkle inclusion proof to verify.
type VerifyInclusionRequest struct {
	LeafCount int         `json:"leaf_count,omitempty"`
	LeafHash  string      `json:"leaf_hash,omitempty"`
	LeafIndex *int        `json:"leaf_index,omitempty"`
	PadCount  *int        `json:"pad_count,omitempty"`
	Root      string      `json:"root,omitempty"`
	Steps     []ProofStep `json:"steps,omitempty"`
	TreeDepth int         `json:"tree_depth,omitempty"`
}

// VerifyInclusionResponse is the verification result.
type VerifyInclusionResponse struct {
	Valid bool `json:"valid,omitempty"`
}
