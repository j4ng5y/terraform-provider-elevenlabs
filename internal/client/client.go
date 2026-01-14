package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
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

// Voices
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

func (c *Client) EditVoiceSettings(voiceID string, settings *models.VoiceSettings) error {
	body, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/voices/"+voiceID+"/settings/edit", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Models
func (c *Client) GetModels() ([]models.Model, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return nil, err
	}

	var models []models.Model
	err = c.doRequest(req, &models)
	return models, err
}

// Projects
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

// Pronunciation Dictionaries
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

func (c *Client) GetPronunciationDictionaries() ([]models.PronunciationDictionary, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/pronunciation-dictionaries", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Dictionaries []models.PronunciationDictionary `json:"pronunciation_dictionaries"`
	}
	err = c.doRequest(req, &wrapper)
	return wrapper.Dictionaries, err
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

// Audio Native
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

// Conversational AI Agents
func (c *Client) GetConvAIAgents() ([]models.ConvAIAgent, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/agents", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Agents []models.ConvAIAgent `json:"agents"`
	}
	err = c.doRequest(req, &wrapper)
	return wrapper.Agents, err
}

func (c *Client) CreateConvAIAgent(addReq *models.CreateConvAIAgentRequest) (*models.ConvAIAgent, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/agents/create", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var agent models.ConvAIAgent
	err = c.doRequest(req, &agent)
	return &agent, err
}

func (c *Client) GetConvAIAgent(agentID string) (*models.ConvAIAgent, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/agents/"+agentID, nil)
	if err != nil {
		return nil, err
	}

	var agent models.ConvAIAgent
	err = c.doRequest(req, &agent)
	return &agent, err
}

func (c *Client) UpdateConvAIAgent(agentID string, updateReq *models.CreateConvAIAgentRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/agents/"+agentID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteConvAIAgent(agentID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/agents/"+agentID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Knowledge Base
func (c *Client) CreateConvAIKnowledgeBase(addReq *models.CreateConvAIKnowledgeBaseRequest) (*models.ConvAIKnowledgeBase, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", addReq.Name)
	if addReq.URL != "" {
		_ = writer.WriteField("url", addReq.URL)
	}
	if addReq.Content != "" {
		_ = writer.WriteField("content", addReq.Content)
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

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/knowledge-base", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var kb models.ConvAIKnowledgeBase
	err = c.doRequest(req, &kb)
	return &kb, err
}

func (c *Client) GetConvAIKnowledgeBase(documentationID string) (*models.ConvAIKnowledgeBase, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/knowledge-base/"+documentationID, nil)
	if err != nil {
		return nil, err
	}

	var kb models.ConvAIKnowledgeBase
	err = c.doRequest(req, &kb)
	return &kb, err
}

func (c *Client) DeleteConvAIKnowledgeBase(documentationID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/knowledge-base/"+documentationID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) ListConvAIKnowledgeBaseDocuments(params *models.ListConvAIKnowledgeBaseDocumentsParams) (*models.ConvAIKnowledgeBaseListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/knowledge-base", nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := req.URL.Query()
		if params.PageSize != nil {
			query.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Search != "" {
			query.Set("search", params.Search)
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
		if params.ShowOnlyOwned != nil {
			query.Set("show_only_owned_documents", strconv.FormatBool(*params.ShowOnlyOwned))
		}
		for _, t := range params.Types {
			if t != "" {
				query.Add("types", t)
			}
		}
		req.URL.RawQuery = query.Encode()
	}

	var list models.ConvAIKnowledgeBaseListResponse
	err = c.doRequest(req, &list)
	return &list, err
}

// Conversational AI Tools
func (c *Client) CreateConvAITool(addReq *models.CreateConvAIToolRequest) (*models.ConvAITool, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/tools", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var tool models.ConvAITool
	err = c.doRequest(req, &tool)
	return &tool, err
}

func (c *Client) GetConvAITool(toolID string) (*models.ConvAITool, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/tools/"+toolID, nil)
	if err != nil {
		return nil, err
	}

	var tool models.ConvAITool
	err = c.doRequest(req, &tool)
	return &tool, err
}

func (c *Client) UpdateConvAITool(toolID string, updateReq *models.CreateConvAIToolRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/tools/"+toolID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteConvAITool(toolID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/tools/"+toolID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) GetConvAITools() ([]models.ConvAITool, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/tools", nil)
	if err != nil {
		return nil, err
	}

	var resp models.ConvAIToolsResponse
	err = c.doRequest(req, &resp)
	return resp.Tools, err
}

// Conversational AI Secrets
func (c *Client) CreateConvAISecret(addReq *models.CreateConvAISecretRequest) (*models.ConvAISecret, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/secrets", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var secret models.ConvAISecret
	err = c.doRequest(req, &secret)
	return &secret, err
}

func (c *Client) GetConvAISecret(secretID string) (*models.ConvAISecret, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/secrets/"+secretID, nil)
	if err != nil {
		return nil, err
	}

	var secret models.ConvAISecret
	err = c.doRequest(req, &secret)
	return &secret, err
}

func (c *Client) UpdateConvAISecret(secretID string, updateReq *models.CreateConvAISecretRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/secrets/"+secretID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteConvAISecret(secretID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/secrets/"+secretID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Agent Testing
func (c *Client) CreateConvAIAgentTest(addReq *models.CreateConvAIAgentTestRequest) (*models.ConvAIAgentTest, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/agent-testing/create", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var test models.ConvAIAgentTest
	err = c.doRequest(req, &test)
	return &test, err
}

func (c *Client) GetConvAIAgentTest(testID string) (*models.ConvAIAgentTest, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/agent-testing/"+testID, nil)
	if err != nil {
		return nil, err
	}

	var test models.ConvAIAgentTest
	err = c.doRequest(req, &test)
	return &test, err
}

func (c *Client) DeleteConvAIAgentTest(testID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/agent-testing/"+testID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI MCP Servers
func (c *Client) GetConvAIMCPServers() ([]models.ConvAIMCPServer, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/mcp-servers", nil)
	if err != nil {
		return nil, err
	}

	var servers []models.ConvAIMCPServer
	err = c.doRequest(req, &servers)
	return servers, err
}

func (c *Client) CreateConvAIMCPServer(addReq *models.CreateConvAIMCPServerRequest) (*models.ConvAIMCPServer, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/mcp-servers", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var server models.ConvAIMCPServer
	err = c.doRequest(req, &server)
	return &server, err
}

func (c *Client) UpdateConvAIMCPServer(mcpServerID string, updateReq *models.CreateConvAIMCPServerRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/mcp-servers/"+mcpServerID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteConvAIMCPServer(mcpServerID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/mcp-servers/"+mcpServerID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Phone Numbers
func (c *Client) GetConvAIPhoneNumbers() ([]models.ConvAIPhoneNumber, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/phone-numbers", nil)
	if err != nil {
		return nil, err
	}

	var numbers []models.ConvAIPhoneNumber
	err = c.doRequest(req, &numbers)
	return numbers, err
}

func (c *Client) ImportConvAIPhoneNumber(addReq *models.ImportPhoneNumberRequest) (*models.ConvAIPhoneNumber, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/phone-numbers", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var phone models.ConvAIPhoneNumber
	err = c.doRequest(req, &phone)
	return &phone, err
}

func (c *Client) UpdateConvAIPhoneNumber(phoneNumberID string, updateReq *models.ImportPhoneNumberRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/phone-numbers/"+phoneNumberID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteConvAIPhoneNumber(phoneNumberID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/phone-numbers/"+phoneNumberID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI WhatsApp Accounts
func (c *Client) ListConvAIWhatsAppAccounts() ([]models.ConvAIWhatsAppAccount, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/whatsapp-accounts", nil)
	if err != nil {
		return nil, err
	}

	var resp models.ConvAIWhatsAppAccountListResponse
	err = c.doRequest(req, &resp)
	return resp.Items, err
}

func (c *Client) ImportConvAIWhatsAppAccount(addReq *models.ImportWhatsAppAccountRequest) (*models.ConvAIWhatsAppAccount, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/whatsapp-accounts", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var account models.ConvAIWhatsAppAccount
	err = c.doRequest(req, &account)
	return &account, err
}

func (c *Client) GetConvAIWhatsAppAccount(phoneNumberID string) (*models.ConvAIWhatsAppAccount, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/whatsapp-accounts/"+phoneNumberID, nil)
	if err != nil {
		return nil, err
	}

	var account models.ConvAIWhatsAppAccount
	err = c.doRequest(req, &account)
	return &account, err
}

func (c *Client) UpdateConvAIWhatsAppAccount(phoneNumberID string, updateReq *models.UpdateWhatsAppAccountRequest) (*models.ConvAIWhatsAppAccount, error) {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/whatsapp-accounts/"+phoneNumberID, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var account models.ConvAIWhatsAppAccount
	err = c.doRequest(req, &account)
	return &account, err
}

func (c *Client) DeleteConvAIWhatsAppAccount(phoneNumberID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/whatsapp-accounts/"+phoneNumberID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Settings
func (c *Client) GetConvAISettings() (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/settings", nil)
	if err != nil {
		return nil, err
	}

	var settings map[string]interface{}
	err = c.doRequest(req, &settings)
	return settings, err
}

func (c *Client) UpdateConvAISettings(settings map[string]interface{}) error {
	body, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/convai/settings", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) GetConvAISecrets() ([]models.ConvAISecret, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/secrets", nil)
	if err != nil {
		return nil, err
	}

	var secrets []models.ConvAISecret
	err = c.doRequest(req, &secrets)
	return secrets, err
}

func (c *Client) GetConvAIConversation(conversationID string) (*map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/conversations/"+conversationID, nil)
	if err != nil {
		return nil, err
	}

	var conversation map[string]interface{}
	err = c.doRequest(req, &conversation)
	return &conversation, err
}

func (c *Client) DeleteConvAIConversation(conversationID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/convai/conversations/"+conversationID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Conversations
func (c *Client) GetConvAIConversations() ([]map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/conversations", nil)
	if err != nil {
		return nil, err
	}

	var conversations []map[string]interface{}
	err = c.doRequest(req, &conversations)
	return conversations, err
}

// Dubbing
type DubbingListResponse struct {
	Dubs       []DubbingMetadata `json:"dubs"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}

type DubbingMetadata struct {
	DubbingID     string `json:"dubbing_id"`
	InputFileName string `json:"input_file_name"`
	InputURL      string `json:"input_url"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

func (c *Client) GetDubs() (*DubbingListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/dubbing", nil)
	if err != nil {
		return nil, err
	}

	var resp DubbingListResponse
	err = c.doRequest(req, &resp)
	return &resp, err
}

func (c *Client) GetDubbing(dubbingID string) (*DubbingMetadata, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/dubbing/"+dubbingID, nil)
	if err != nil {
		return nil, err
	}

	var dub DubbingMetadata
	err = c.doRequest(req, &dub)
	return &dub, err
}

func (c *Client) DeleteDubbing(dubbingID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/dubbing/"+dubbingID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Conversational AI Dashboard Settings
func (c *Client) GetConvAIDashboardSettings() (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/dashboard/settings", nil)
	if err != nil {
		return nil, err
	}

	var settings map[string]interface{}
	err = c.doRequest(req, &settings)
	return settings, err
}

// Conversational AI Batch Calling
func (c *Client) GetConvAIBatchCalls() ([]map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/batch-calling/workspace", nil)
	if err != nil {
		return nil, err
	}

	var batches []map[string]interface{}
	err = c.doRequest(req, &batches)
	return batches, err
}

func (c *Client) GetConvAIBatchCall(batchID string) (*map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/convai/batch-calling/"+batchID, nil)
	if err != nil {
		return nil, err
	}

	var batch map[string]interface{}
	err = c.doRequest(req, &batch)
	return &batch, err
}

func (c *Client) CancelConvAIBatchCall(batchID string) error {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/convai/batch-calling/"+batchID+"/cancel", nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Workspace Webhooks
func (c *Client) ListWorkspaceWebhooks() ([]models.WorkspaceWebhook, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var webhooks []models.WorkspaceWebhook
	err = c.doRequest(req, &webhooks)
	return webhooks, err
}

func (c *Client) CreateWorkspaceWebhook(addReq *models.CreateWorkspaceWebhookRequest) (*models.WorkspaceWebhook, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/webhooks", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var webhook models.WorkspaceWebhook
	err = c.doRequest(req, &webhook)
	return &webhook, err
}

func (c *Client) GetWorkspaceWebhook(webhookID string) (*models.WorkspaceWebhook, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/webhooks/"+webhookID, nil)
	if err != nil {
		return nil, err
	}

	var webhook models.WorkspaceWebhook
	err = c.doRequest(req, &webhook)
	return &webhook, err
}

func (c *Client) UpdateWorkspaceWebhook(webhookID string, updateReq *models.CreateWorkspaceWebhookRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/workspace/webhooks/"+webhookID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteWorkspaceWebhook(webhookID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/workspace/webhooks/"+webhookID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Workspace Members
func (c *Client) GetWorkspaceMembers() ([]models.WorkspaceMember, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/members", nil)
	if err != nil {
		return nil, err
	}

	var members []models.WorkspaceMember
	err = c.doRequest(req, &members)
	return members, err
}

func (c *Client) UpdateWorkspaceMember(userID string, updateReq *models.UpdateWorkspaceMemberRequest) error {
	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/members/"+userID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Workspace Invites
func (c *Client) GetWorkspaceInvites() ([]models.WorkspaceInvite, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/invites", nil)
	if err != nil {
		return nil, err
	}

	var invites []models.WorkspaceInvite
	err = c.doRequest(req, &invites)
	return invites, err
}

func (c *Client) CreateWorkspaceInvite(addReq *models.CreateWorkspaceInviteRequest) error {
	body, err := json.Marshal(addReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/invites/add", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) DeleteWorkspaceInvite(email string) error {
	body := map[string]string{"email": email}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/workspace/invites", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Workspace Groups
func (c *Client) SearchWorkspaceGroups(query string) ([]models.WorkspaceGroup, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/groups/search?search="+query, nil)
	if err != nil {
		return nil, err
	}

	var groups []models.WorkspaceGroup
	err = c.doRequest(req, &groups)
	return groups, err
}

func (c *Client) AddWorkspaceGroupMember(groupID, email string) error {
	body := map[string]string{"email": email}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/groups/"+groupID+"/members", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) RemoveWorkspaceGroupMember(groupID, email string) error {
	body := map[string]string{"email": email}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/groups/"+groupID+"/members/remove", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Service Accounts
func (c *Client) GetWorkspaceServiceAccounts() ([]models.WorkspaceServiceAccount, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/service-accounts", nil)
	if err != nil {
		return nil, err
	}

	var accounts []models.WorkspaceServiceAccount
	err = c.doRequest(req, &accounts)
	return accounts, err
}

func (c *Client) GetServiceAccountAPIKeys(userID string) ([]models.ServiceAccountKey, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/service-accounts/"+userID+"/api-keys", nil)
	if err != nil {
		return nil, err
	}

	var keys []models.ServiceAccountKey
	err = c.doRequest(req, &keys)
	return keys, err
}

func (c *Client) CreateServiceAccountKey(userID string, addReq *models.CreateServiceAccountKeyRequest) (*models.ServiceAccountKey, error) {
	body, err := json.Marshal(addReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/service-accounts/"+userID+"/api-keys", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var key models.ServiceAccountKey
	err = c.doRequest(req, &key)
	return &key, err
}

func (c *Client) DeleteServiceAccountKey(userID, keyID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/service-accounts/"+userID+"/api-keys/"+keyID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Workspace Resources
func (c *Client) GetWorkspaceResources() ([]models.WorkspaceResource, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/resources", nil)
	if err != nil {
		return nil, err
	}

	var resources []models.WorkspaceResource
	err = c.doRequest(req, &resources)
	return resources, err
}

func (c *Client) GetWorkspaceResource(resourceID string) (*models.WorkspaceResource, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/workspace/resources/"+resourceID, nil)
	if err != nil {
		return nil, err
	}

	var resource models.WorkspaceResource
	err = c.doRequest(req, &resource)
	return &resource, err
}

func (c *Client) ShareResource(resourceID, resourceType, email, role string) error {
	body := map[string]interface{}{
		"email":         email,
		"resource_type": resourceType,
		"role":          role,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/resources/"+resourceID+"/share", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) UnshareResource(resourceID, resourceType, email string) error {
	body := map[string]interface{}{
		"email":         email,
		"resource_type": resourceType,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/resources/"+resourceID+"/unshare", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) CopyResourceToWorkspace(resourceID, resourceType, targetWorkspaceID string) error {
	body := map[string]interface{}{
		"resource_type":       resourceType,
		"target_workspace_id": targetWorkspaceID,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/workspace/resources/"+resourceID+"/copy-to-workspace", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Shared Voices
func (c *Client) AddSharedVoice(publicUserID, voiceID, newName string) (string, error) {
	body := map[string]string{"new_name": newName}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/voices/add/"+publicUserID+"/"+voiceID, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	var result struct {
		VoiceID string `json:"voice_id"`
	}
	err = c.doRequest(req, &result)
	return result.VoiceID, err
}

// Voice Samples (Standalone)
func (c *Client) AddVoiceSample(voiceID string, addReq *models.AddVoiceSampleRequest) (*models.VoiceSample, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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

	req, err := http.NewRequest(http.MethodPost, baseURL+"/voices/pvc/"+voiceID+"/samples", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var sample models.VoiceSample
	err = c.doRequest(req, &sample)
	return &sample, err
}

func (c *Client) DeleteVoiceSample(voiceID, sampleID string) error {
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/voices/"+voiceID+"/samples/"+sampleID, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}
