package models

type Model struct {
	ModelID                  string          `json:"model_id"`
	Name                     string          `json:"name"`
	Description              string          `json:"description"`
	Languages                []ModelLanguage `json:"languages"`
	CanDoTextToSpeech        bool            `json:"can_do_text_to_speech"`
	CanDoVoiceConversion     bool            `json:"can_do_voice_conversion"`
	CanUseSpeakerBoost       bool            `json:"can_use_speaker_boost"`
	CanUseStyle              bool            `json:"can_use_style"`
	RequiresGuidance         bool            `json:"requires_guidance"`
	MaxCharactersRequestFree int             `json:"max_characters_request_free"`
	MaxCharactersRequestPaid int             `json:"max_characters_request_paid"`
	TokenCostFactor          float64         `json:"token_cost_factor"`
}

type ModelLanguage struct {
	LanguageID string `json:"language_id"`
	Name       string `json:"name"`
}

type Project struct {
	ProjectID               string `json:"project_id"`
	Name                    string `json:"name"`
	CreateDate              int64  `json:"create_date"`
	LastApplyDate           int64  `json:"last_apply_date"`
	DefaultModelID          string `json:"default_model_id"`
	DefaultParagraphVoiceID string `json:"default_paragraph_voice_id"`
	DefaultTitleVoiceID     string `json:"default_title_voice_id"`
	CanBeDownloaded         bool   `json:"can_be_downloaded"`
	State                   string `json:"state"`
}

type CreateProjectRequest struct {
	Name                    string `json:"name"`
	DefaultModelID          string `json:"default_model_id,omitempty"`
	DefaultParagraphVoiceID string `json:"default_paragraph_voice_id,omitempty"`
	DefaultTitleVoiceID     string `json:"default_title_voice_id,omitempty"`
}
