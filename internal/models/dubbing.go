package models

type DubbingProject struct {
	DubbingID       string   `json:"dubbing_id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	TargetLanguages []string `json:"target_languages"`
}

type CreateDubbingRequest struct {
	Name          string `json:"name"`
	SourceURL     string `json:"source_url,omitempty"`
	SourceLang    string `json:"source_lang,omitempty"`
	TargetLang    string `json:"target_lang"`
	NumSpeakers   int    `json:"num_speakers,omitempty"`
	Watermark     bool   `json:"watermark,omitempty"`
	DubbingStudio bool   `json:"dubbing_studio,omitempty"`
	Mode          string `json:"mode,omitempty"`
	FilePath      string `json:"-"`
}
