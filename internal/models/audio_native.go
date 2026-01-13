package models

type AudioNativeProject struct {
	ProjectID   string `json:"project_id"`
	HTMLSnippet string `json:"html_snippet,omitempty"`
}

type AudioNativeSettings struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	TextColor       string `json:"text_color"`
	BackgroundColor string `json:"background_color"`
	Status          string `json:"status"`
}

type CreateAudioNativeRequest struct {
	Name                       string   `json:"name"`
	Title                      string   `json:"title,omitempty"`
	Author                     string   `json:"author,omitempty"`
	VoiceID                    string   `json:"voice_id,omitempty"`
	ModelID                    string   `json:"model_id,omitempty"`
	TextColor                  string   `json:"text_color,omitempty"`
	BackgroundColor            string   `json:"background_color,omitempty"`
	AutoConvert                bool     `json:"auto_convert,omitempty"`
	ApplyTextNormalization     string   `json:"apply_text_normalization,omitempty"`
	PronunciationDictionaryIDs []string `json:"pronunciation_dictionary_locators,omitempty"`
	FilePath                   string   `json:"-"`
}
