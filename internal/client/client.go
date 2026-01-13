package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const baseURL = "https://api.elevenlabs.io/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	req.Header.Set("xi-api-key", c.apiKey)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

func (c *Client) GetVoices() ([]models.Voice, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/voices", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Voices []models.Voice `json:"voices"`
	}
	err = c.doRequest(req, &wrapper)
	return wrapper.Voices, err
}

func (c *Client) GetVoice(voiceID string) (*models.Voice, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/voices/"+voiceID, nil)
	if err != nil {
		return nil, err
	}

	var voice models.Voice
	err = c.doRequest(req, &voice)
	return &voice, err
}

func (c *Client) AddVoice(addReq *models.AddVoiceRequest) (*models.Voice, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	if addReq.Description != "" {
		_ = writer.WriteField("description", addReq.Description)
	}
	if len(addReq.Labels) > 0 {
		labelsJSON, _ := json.Marshal(addReq.Labels)
		_ = writer.WriteField("labels", string(labelsJSON))
	}

	for _, filePath := range addReq.Files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("files", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/voices/add", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var result struct {
		VoiceID string `json:"voice_id"`
	}
	err = c.doRequest(req, &result)
	if err != nil {
		return nil, err
	}

	return c.GetVoice(result.VoiceID)
}

func (c *Client) EditVoice(voiceID string, addReq *models.AddVoiceRequest) error {
	// Note: The API for edit is similar to add but uses a different endpoint
	// and voice_id in path.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	if addReq.Description != "" {
		_ = writer.WriteField("description", addReq.Description)
	}
	if len(addReq.Labels) > 0 {
		labelsJSON, _ := json.Marshal(addReq.Labels)
		_ = writer.WriteField("labels", string(labelsJSON))
	}

	// Optional: files can be empty if only metadata is changing
	for _, filePath := range addReq.Files {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("files", filepath.Base(filePath))
		if err != nil {
			return err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/voices/"+voiceID+"/edit", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(req, nil)
}

func (c *Client) DeleteVoice(voiceID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/voices/"+voiceID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) GetModels() ([]models.Model, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return nil, err
	}

	var models []models.Model
	err = c.doRequest(req, &models)
	return models, err
}

func (c *Client) GetProjects() ([]models.Project, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/projects", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Projects []models.Project `json:"projects"`
	}
	err = c.doRequest(req, &wrapper)
	return wrapper.Projects, err
}

func (c *Client) GetProject(projectID string) (*models.Project, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/projects/"+projectID, nil)
	if err != nil {
		return nil, err
	}

	var project models.Project
	err = c.doRequest(req, &project)
	return &project, err
}

func (c *Client) CreateProject(createReq *models.CreateProjectRequest) (*models.Project, error) {
	body, err := json.Marshal(createReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/projects", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var project models.Project
	err = c.doRequest(req, &project)
	return &project, err
}

func (c *Client) DeleteProject(projectID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/projects/"+projectID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) AddPronunciationDictionaryFromRules(addReq *models.AddPronunciationDictionaryFromRulesRequest) (*models.PronunciationDictionary, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/pronunciation-dictionaries/add-from-rules", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var dict models.PronunciationDictionary
	err = c.doRequest(req, &dict)
	return &dict, err
}

func (c *Client) AddPronunciationDictionaryFromFile(addReq *models.AddPronunciationDictionaryFromFileRequest) (*models.PronunciationDictionary, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	if addReq.Description != "" {
		_ = writer.WriteField("description", addReq.Description)
	}

	file, err := os.Open(addReq.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(addReq.FilePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/pronunciation-dictionaries/add-from-file", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var dict models.PronunciationDictionary
	err = c.doRequest(req, &dict)
	return &dict, err
}

func (c *Client) GetPronunciationDictionary(dictionaryID string) (*models.PronunciationDictionary, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/pronunciation-dictionaries/"+dictionaryID, nil)
	if err != nil {
		return nil, err
	}

	var dict models.PronunciationDictionary
	err = c.doRequest(req, &dict)
	return &dict, err
}

func (c *Client) ArchivePronunciationDictionary(dictionaryID string) error {
	body := map[string]bool{"archived": true}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/pronunciation-dictionaries/"+dictionaryID, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) CreateAudioNative(addReq *models.CreateAudioNativeRequest) (*models.AudioNativeProject, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	if addReq.Title != "" {
		_ = writer.WriteField("title", addReq.Title)
	}
	if addReq.Author != "" {
		_ = writer.WriteField("author", addReq.Author)
	}
	if addReq.VoiceID != "" {
		_ = writer.WriteField("voice_id", addReq.VoiceID)
	}
	if addReq.ModelID != "" {
		_ = writer.WriteField("model_id", addReq.ModelID)
	}
	if addReq.TextColor != "" {
		_ = writer.WriteField("text_color", addReq.TextColor)
	}
	if addReq.BackgroundColor != "" {
		_ = writer.WriteField("background_color", addReq.BackgroundColor)
	}
	if addReq.AutoConvert {
		_ = writer.WriteField("auto_convert", "true")
	}

	if addReq.FilePath != "" {
		file, err := os.Open(addReq.FilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(addReq.FilePath))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/audio-native", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var project models.AudioNativeProject
	err = c.doRequest(req, &project)
	return &project, err
}

func (c *Client) GetAudioNativeSettings(projectID string) (*models.AudioNativeSettings, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/audio-native/"+projectID+"/settings", nil)
	if err != nil {
		return nil, err
	}

	var settings models.AudioNativeSettings
	err = c.doRequest(req, &settings)
	return &settings, err
}

func (c *Client) CreateDubbingProject(addReq *models.CreateDubbingRequest) (*models.DubbingProject, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	_ = writer.WriteField("target_lang", addReq.TargetLang)
	if addReq.SourceURL != "" {
		_ = writer.WriteField("source_url", addReq.SourceURL)
	}
	if addReq.SourceLang != "" {
		_ = writer.WriteField("source_lang", addReq.SourceLang)
	}
	if addReq.NumSpeakers > 0 {
		_ = writer.WriteField("num_speakers", fmt.Sprintf("%d", addReq.NumSpeakers))
	}
	if addReq.Watermark {
		_ = writer.WriteField("watermark", "true")
	}
	if addReq.DubbingStudio {
		_ = writer.WriteField("dubbing_studio", "true")
	}
	if addReq.Mode != "" {
		_ = writer.WriteField("mode", addReq.Mode)
	}

	if addReq.FilePath != "" {
		file, err := os.Open(addReq.FilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(addReq.FilePath))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/dubbing", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var project models.DubbingProject
	err = c.doRequest(req, &project)
	return &project, err
}

func (c *Client) GetDubbingProject(dubbingID string) (*models.DubbingProject, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/dubbing/"+dubbingID, nil)
	if err != nil {
		return nil, err
	}

	var project models.DubbingProject
	err = c.doRequest(req, &project)
	return &project, err
}

func (c *Client) DeleteDubbingProject(dubbingID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/dubbing/"+dubbingID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}
