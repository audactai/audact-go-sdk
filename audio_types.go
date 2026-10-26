package audact

// AudioFormat is the encoding of submitted audio.
type AudioFormat string

// Audio formats accepted by the audio validator.
const (
	// AudioFormatPCM is 16 kHz signed 16-bit little-endian PCM.
	AudioFormatPCM AudioFormat = "pcm_s16le_16000"
	// AudioFormatWAV is WAV.
	AudioFormatWAV AudioFormat = "wav"
	// AudioFormatMP3 is MP3.
	AudioFormatMP3 AudioFormat = "mp3"
)

// AudioCallContext carries per-call metadata.
type AudioCallContext struct {
	DisclosureDelivered *bool  `json:"disclosure_delivered,omitempty"`
	Language            string `json:"language,omitempty"`
}

// ValidateAudioRequest is the body of Audio.Validate.
type ValidateAudioRequest struct {
	ASRProvider string            `json:"asr_provider,omitempty"`
	AudioData   string            `json:"audio_data,omitempty"`
	AudioFormat AudioFormat       `json:"audio_format,omitempty"`
	CallContext *AudioCallContext `json:"call_context,omitempty"`
	CallID      string            `json:"call_id,omitempty"`
	// ConfigSealHash is the call's config seal. Optional: when empty the
	// call's pinned seal is used; it is required only for speech in a call
	// with no pinned seal.
	ConfigSealHash     string `json:"config_seal_hash,omitempty"`
	ConfigurationID    string `json:"configuration_id,omitempty"`
	ParallelASR        *bool  `json:"parallel_asr,omitempty"`
	SegmentationMethod string `json:"segmentation_method,omitempty"`
	TurnNumber         int    `json:"turn_number,omitempty"`
}

// AudioTextCheckResult is one text check outcome.
type AudioTextCheckResult struct {
	CheckID string `json:"check_id,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Verdict string `json:"verdict,omitempty"`
}

// AudioParallelASRSummary summarises a parallel ASR run.
type AudioParallelASRSummary struct {
	ConfidenceSpread float64 `json:"confidence_spread,omitempty"`
	Converged        bool    `json:"converged,omitempty"`
	ModelCount       int     `json:"model_count,omitempty"`
}

// TranscriptionStatus reports whether the ASR produced text for an audio segment.
type TranscriptionStatus string

const (
	// TranscriptionStatusTranscribed means the ASR produced text and the segment was checked.
	TranscriptionStatusTranscribed TranscriptionStatus = "transcribed"
	// TranscriptionStatusNoSpeech means the ASR produced no text: the segment is ALLOW but nothing was checked.
	TranscriptionStatusNoSpeech TranscriptionStatus = "no_speech"
)

// AudioSegmentResult is the outcome for one audio segment.
type AudioSegmentResult struct {
	Action              string                   `json:"action,omitempty"`
	ASRConfidence       float64                  `json:"asr_confidence,omitempty"`
	ASRModelID          string                   `json:"asr_model_id,omitempty"`
	AudioHash           string                   `json:"audio_hash,omitempty"`
	CheckResults        []AudioTextCheckResult   `json:"check_results,omitempty"`
	ParallelASRResult   *AudioParallelASRSummary `json:"parallel_asr_result,omitempty"`
	SegmentIndex        int                      `json:"segment_index,omitempty"`
	SegmentationMethod  string                   `json:"segmentation_method,omitempty"`
	TextHash            string                   `json:"text_hash,omitempty"`
	TranscribedText     string                   `json:"transcribed_text,omitempty"`
	TranscriptionStatus TranscriptionStatus      `json:"transcription_status,omitempty"`
	Verdict             string                   `json:"verdict,omitempty"`
}

// AudioComplianceReceipt is the compliance receipt of an audio validation.
//
// CallID, TurnNumber, Timestamp, ConfigSealHash, AuthMechanism and
// ReceiptFormatVersion are absent on receipts issued by older servers.
type AudioComplianceReceipt struct {
	ASRConfidence        float64 `json:"asr_confidence,omitempty"`
	AuthMechanism        string  `json:"auth_mechanism,omitempty"`
	CallID               string  `json:"call_id,omitempty"`
	ConfigSealHash       string  `json:"config_seal_hash,omitempty"`
	ConfigVersion        int     `json:"config_version,omitempty"`
	EvidenceHash         string  `json:"evidence_hash,omitempty"`
	InputHash            string  `json:"input_hash,omitempty"`
	OutputHash           string  `json:"output_hash,omitempty"`
	PolicyVersion        string  `json:"policy_version,omitempty"`
	ReceiptFormatVersion string  `json:"receipt_format_version,omitempty"`
	SegmentationMethod   string  `json:"segmentation_method,omitempty"`
	Signature            string  `json:"signature,omitempty"`
	TextHash             string  `json:"text_hash,omitempty"`
	Timestamp            string  `json:"timestamp,omitempty"`
	TurnNumber           int     `json:"turn_number,omitempty"`
	Verdict              string  `json:"verdict,omitempty"`
}

// ValidateAudioResponse is the result of Audio.Validate.
type ValidateAudioResponse struct {
	EvidenceEventID     string                  `json:"evidence_event_id,omitempty"`
	EvidenceHash        string                  `json:"evidence_hash,omitempty"`
	EvidenceWriteStatus string                  `json:"evidence_write_status,omitempty"`
	FailOpen            bool                    `json:"fail_open,omitempty"`
	Receipt             *AudioComplianceReceipt `json:"receipt,omitempty"`
	Segments            []AudioSegmentResult    `json:"segments,omitempty"`
	Verdict             string                  `json:"verdict,omitempty"`
}

// AudioReconciliationRequest is the body of Audio.Reconcile.
type AudioReconciliationRequest struct {
	ASRProvider     string      `json:"asr_provider,omitempty"`
	AudioData       string      `json:"audio_data,omitempty"`
	AudioFormat     AudioFormat `json:"audio_format,omitempty"`
	CallID          string      `json:"call_id,omitempty"`
	ConfigurationID string      `json:"configuration_id,omitempty"`
	Language        string      `json:"language,omitempty"`
}

// ReconcileAudioParams are the options of Audio.Reconcile.
type ReconcileAudioParams struct {
	// Limit caps segment_results (default 100, max 200).
	Limit int
	// AfterSequence returns only segments with a greater sequence number.
	AfterSequence int
	// IdempotencyKey is sent as the Idempotency-Key header when non-empty.
	IdempotencyKey string
}

// AudioSegmentReconciliation is the reconciliation of one segment.
type AudioSegmentReconciliation struct {
	MismatchClassification string `json:"mismatch_classification,omitempty"`
	SequenceNumber         int    `json:"sequence_number,omitempty"`
	Status                 string `json:"status,omitempty"`
	TranscriptHash         string `json:"transcript_hash,omitempty"`
	ValidatedHash          string `json:"validated_hash,omitempty"`
}

// AudioReconciliationResponse is the result of Audio.Reconcile.
type AudioReconciliationResponse struct {
	ASRConfidence       float64                      `json:"asr_confidence,omitempty"`
	ASRModelID          string                       `json:"asr_model_id,omitempty"`
	CallID              string                       `json:"call_id,omitempty"`
	HasMore             bool                         `json:"has_more,omitempty"`
	MatchedSegments     int                          `json:"matched_segments,omitempty"`
	MismatchedSegments  int                          `json:"mismatched_segments,omitempty"`
	SegmentResults      []AudioSegmentReconciliation `json:"segment_results,omitempty"`
	TotalSegments       int                          `json:"total_segments,omitempty"`
	UnvalidatedSegments int                          `json:"unvalidated_segments,omitempty"`
}
