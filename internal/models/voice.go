package models

type Voice struct {
	VoiceID    string            `json:"voice_id"`
	Name       string            `json:"name"`
	Samples    []VoiceSample     `json:"samples"`
	Category   string            `json:"category"`
	Labels     map[string]string `json:"labels"`
	Settings   *VoiceSettings    `json:"settings"`
	PreviewURL string            `json:"preview_url"`
}

type VoiceSample struct {
	SampleID  string `json:"sample_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int    `json:"size_bytes"`
	Hash      string `json:"hash"`
}

type AddVoiceSampleRequest struct {
	FilePath string `json:"-"`
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
	Speed           float64 `json:"speed,omitempty"`
}

type AddVoiceRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Files       []string          `json:"-"` // Handled as multipart
}
