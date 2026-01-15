package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

func TestClient_CreatePVCVoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/voices/pvc" {
			t.Errorf("Expected POST /voices/pvc, got %s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"voice_id": "test-voice-id"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	req := &models.CreatePVCVoiceRequest{
		Name:     "Test Voice",
		Language: "en",
	}

	voice, err := client.CreatePVCVoice(req)
	if err != nil {
		t.Fatalf("CreatePVCVoice failed: %v", err)
	}

	if voice.VoiceID != "test-voice-id" {
		t.Errorf("Expected voice ID 'test-voice-id', got '%s'", voice.VoiceID)
	}
}

func TestClient_GetPVCVoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/voices/pvc/test-voice-id" {
			t.Errorf("Expected GET /voices/pvc/test-voice-id, got %s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"voice_id": "test-voice-id",
			"name": "Test Voice",
			"language": "en",
			"description": "Test description",
			"state": "ready",
			"verification": "verified"
		}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	voice, err := client.GetPVCVoice("test-voice-id")
	if err != nil {
		t.Fatalf("GetPVCVoice failed: %v", err)
	}

	if voice.VoiceID != "test-voice-id" {
		t.Errorf("Expected voice ID 'test-voice-id', got '%s'", voice.VoiceID)
	}
	if voice.Name != "Test Voice" {
		t.Errorf("Expected name 'Test Voice', got '%s'", voice.Name)
	}
}

func TestClient_ListPVCVoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/voices/pvc" {
			t.Errorf("Expected GET /voices/pvc, got %s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"voices": [
				{
					"voice_id": "voice-1",
					"name": "Voice 1",
					"language": "en"
				},
				{
					"voice_id": "voice-2",
					"name": "Voice 2",
					"language": "es"
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	response, err := client.ListPVCVoices()
	if err != nil {
		t.Fatalf("ListPVCVoices failed: %v", err)
	}

	if len(response.Voices) != 2 {
		t.Errorf("Expected 2 voices, got %d", len(response.Voices))
	}
	if response.Voices[0].VoiceID != "voice-1" {
		t.Errorf("Expected first voice ID 'voice-1', got '%s'", response.Voices[0].VoiceID)
	}
}
