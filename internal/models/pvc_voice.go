package models

// PVCVoice represents a Professional Voice Cloning voice
type PVCVoice struct {
	VoiceID      string            `json:"voice_id"`
	Name         string            `json:"name"`
	Language     string            `json:"language"`
	Description  string            `json:"description,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	State        string            `json:"state,omitempty"`        // training state
	Verification string            `json:"verification,omitempty"` // verification status
	Samples      []PVCVoiceSample  `json:"samples,omitempty"`
	Settings     *VoiceSettings    `json:"settings,omitempty"`
	CreatedAt    string            `json:"created_at,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty"`
}

// PVCVoiceSample represents a training sample for a PVC voice
type PVCVoiceSample struct {
	SampleID      string  `json:"sample_id"`
	FileName      string  `json:"file_name"`
	MimeType      string  `json:"mime_type"`
	SizeBytes     int     `json:"size_bytes"`
	Hash          string  `json:"hash"`
	State         string  `json:"state,omitempty"` // processing state
	Transcription string  `json:"transcription,omitempty"`
	Duration      float64 `json:"duration,omitempty"`
	SampleRate    int     `json:"sample_rate,omitempty"`
	Channels      int     `json:"channels,omitempty"`
}

// CreatePVCVoiceRequest represents the request to create a new PVC voice
type CreatePVCVoiceRequest struct {
	Name        string            `json:"name"`
	Language    string            `json:"language"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// UpdatePVCVoiceRequest represents the request to update a PVC voice
type UpdatePVCVoiceRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// AddPVCVoiceSampleRequest represents the request to add a sample to a PVC voice
type AddPVCVoiceSampleRequest struct {
	FilePath string `json:"-"` // Handled as multipart form data
}

// UpdatePVCVoiceSampleRequest represents the request to update a PVC voice sample
type UpdatePVCVoiceSampleRequest struct {
	Transcription string `json:"transcription,omitempty"`
}

// PVCVoiceTrainingRequest represents the request to start training a PVC voice
type PVCVoiceTrainingRequest struct {
	// Training parameters would go here if any
}

// PVCVoiceVerificationRequest represents the request for manual verification
type PVCVoiceVerificationRequest struct {
	// Verification parameters would go here if any
}

// PVCVoiceCaptchaResponse represents the captcha response for verification
type PVCVoiceCaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	ImageURL  string `json:"image_url"`
}

// PVCVoiceCaptchaRequest represents the request to submit captcha solution
type PVCVoiceCaptchaRequest struct {
	CaptchaID string `json:"captcha_id"`
	Solution  string `json:"solution"`
}

// PVCVoiceListResponse represents the response when listing PVC voices
type PVCVoiceListResponse struct {
	Voices []PVCVoice `json:"voices"`
}

// PVCVoiceSampleListResponse represents the response when listing PVC voice samples
type PVCVoiceSampleListResponse struct {
	Samples []PVCVoiceSample `json:"samples"`
}
