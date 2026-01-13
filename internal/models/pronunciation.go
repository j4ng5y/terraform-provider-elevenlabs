package models

type PronunciationDictionary struct {
	ID               string `json:"id"`
	LatestVersionID  string `json:"latest_version_id"`
	Name             string `json:"name"`
	CreatedBy        string `json:"created_by"`
	CreationTimeUnix int64  `json:"creation_time_unix"`
	ArchivedTimeUnix int64  `json:"archived_time_unix,omitempty"`
}

type PronunciationRule struct {
	Type            string `json:"type"`
	StringToReplace string `json:"string_to_replace"`
	Alias           string `json:"alias,omitempty"`
	Phoneme         string `json:"phoneme,omitempty"`
	Alphabet        string `json:"alphabet,omitempty"`
}

type AddPronunciationDictionaryFromRulesRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Rules       []PronunciationRule `json:"rules"`
}

type AddPronunciationDictionaryFromFileRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	FilePath    string `json:"-"`
}
