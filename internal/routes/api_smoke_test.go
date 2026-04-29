package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bakasub-backend/internal/db"
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/routes"
	"bakasub-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type apiResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type languagesResponse struct {
	Languages []models.Language `json:"languages"`
}

type presetsResponse struct {
	Presets []models.TranslationPreset `json:"presets"`
}

type logsResponse struct {
	Logs  []models.LogEntry `json:"logs"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type jobsResponse struct {
	Jobs  []models.TranslationJob `json:"jobs"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
}

func TestSQLiteAPISmoke(t *testing.T) {
	tempDir := t.TempDir()
	database := openSQLiteTestDB(t, filepath.Join(tempDir, "smoke.db"))
	server := newSQLiteTestServer(t, database, "smoke-secret")
	defer server.Close()

	mediaDir := filepath.Join(tempDir, "media")
	mustMkdirAll(t, filepath.Join(mediaDir, "season-01"))
	mustWriteFile(t, filepath.Join(mediaDir, "episode01.mkv"), "fake video")
	mustWriteFile(t, filepath.Join(mediaDir, "episode01.srt"), "1\n00:00:00,000 --> 00:00:01,000\nhello\n")
	mustWriteFile(t, filepath.Join(mediaDir, "season-01", "nested.txt"), "nested")

	client := server.Client()
	baseURL := server.URL + "/api/v1"

	t.Run("health", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("health request failed: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("reading health response failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected health status 200, got %d with body %s", resp.StatusCode, string(body))
		}

		if strings.TrimSpace(string(body)) != "OK" {
			t.Fatalf("expected health body OK, got %q", string(body))
		}
	})

	t.Run("config and favorites", func(t *testing.T) {
		configResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/config/", nil, http.StatusOK)
		config := decodeResponseData[models.UserConfig](t, configResp)
		if config.DefaultPreset != "anime" {
			t.Fatalf("expected default preset anime, got %q", config.DefaultPreset)
		}

		mustJSONRequest(t, client, http.MethodPut, baseURL+"/config/", map[string]any{
			"default_model":           "smoke/model",
			"default_preset":          "anime",
			"default_language":        "en",
			"remove_sdh_default":      true,
			"video_timeout_minutes":   45,
			"log_retention_days":      14,
			"openrouter_api_key":      "openrouter-secret",
			"tmdb_access_token":       "tmdb-secret",
			"tmdb_metadata_enabled":   true,
			"concurrent_translations": 2,
			"max_retries":             4,
			"base_retry_delay":        3,
		}, http.StatusOK)

		updatedResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/config/", nil, http.StatusOK)
		updatedConfig := decodeResponseData[models.UserConfig](t, updatedResp)
		if updatedConfig.OpenRouterApiKey != "openrouter-secret" {
			t.Fatalf("expected decrypted openrouter key, got %q", updatedConfig.OpenRouterApiKey)
		}
		if updatedConfig.DefaultLanguage != "en" || !updatedConfig.RemoveSdhDefault || updatedConfig.ConcurrentTranslations != 2 {
			t.Fatalf("unexpected updated config payload: %+v", updatedConfig)
		}

		var storedKey string
		if err := database.QueryRow("SELECT openrouter_api_key FROM user_config LIMIT 1").Scan(&storedKey); err != nil {
			t.Fatalf("querying stored openrouter key failed: %v", err)
		}
		if storedKey == "openrouter-secret" {
			t.Fatal("expected encrypted openrouter key in SQLite, got plaintext")
		}

		mustJSONRequest(t, client, http.MethodPut, baseURL+"/favorites/", map[string]any{
			"favorite_models": []string{"model/a", "model/b"},
		}, http.StatusOK)

		favoritesResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/favorites/", nil, http.StatusOK)
		favorites := decodeResponseData[[]string](t, favoritesResp)
		if len(favorites) != 2 || favorites[0] != "model/a" || favorites[1] != "model/b" {
			t.Fatalf("unexpected favorites payload: %+v", favorites)
		}
	})

	t.Run("languages", func(t *testing.T) {
		languagesResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/languages/", nil, http.StatusOK)
		languages := decodeResponseData[languagesResponse](t, languagesResp)
		if len(languages.Languages) == 0 {
			t.Fatal("expected seeded languages from migrations")
		}

		mustJSONRequest(t, client, http.MethodPost, baseURL+"/languages/", map[string]any{
			"code": "zz-smoke",
			"name": "Smoke Language",
		}, http.StatusOK)

		mustJSONRequest(t, client, http.MethodPut, baseURL+"/languages/", map[string]any{
			"code": "zz-smoke",
			"name": "Smoke Language Updated",
		}, http.StatusOK)

		languagesResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/languages/", nil, http.StatusOK)
		languages = decodeResponseData[languagesResponse](t, languagesResp)
		if !hasLanguage(languages.Languages, "zz-smoke", "Smoke Language Updated") {
			t.Fatalf("expected updated smoke language in response, got %+v", languages.Languages)
		}

		mustJSONRequest(t, client, http.MethodDelete, baseURL+"/languages/", map[string]any{
			"code": "zz-smoke",
		}, http.StatusOK)

		languagesResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/languages/", nil, http.StatusOK)
		languages = decodeResponseData[languagesResponse](t, languagesResp)
		if hasLanguage(languages.Languages, "zz-smoke", "Smoke Language Updated") {
			t.Fatalf("expected smoke language to be deleted, got %+v", languages.Languages)
		}
	})

	t.Run("presets", func(t *testing.T) {
		presetsResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/presets/", nil, http.StatusOK)
		presets := decodeResponseData[presetsResponse](t, presetsResp)
		if len(presets.Presets) == 0 {
			t.Fatal("expected seeded presets from migrations")
		}

		mustJSONRequest(t, client, http.MethodPost, baseURL+"/presets/", map[string]any{
			"alias":         "smoke-preset",
			"name":          "Smoke Preset",
			"system_prompt": "translate like a smoke test",
			"batch_size":    1234,
			"temperature":   0.4,
		}, http.StatusOK)

		presetsResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/presets/", nil, http.StatusOK)
		presets = decodeResponseData[presetsResponse](t, presetsResp)
		preset, found := findPresetByAlias(presets.Presets, "smoke-preset")
		if !found {
			t.Fatalf("expected smoke preset in response, got %+v", presets.Presets)
		}

		mustJSONRequest(t, client, http.MethodPut, baseURL+"/presets/", map[string]any{
			"id":            preset.ID,
			"name":          "Smoke Preset Updated",
			"system_prompt": "updated prompt",
			"temperature":   0.6,
		}, http.StatusOK)

		presetsResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/presets/", nil, http.StatusOK)
		presets = decodeResponseData[presetsResponse](t, presetsResp)
		preset, found = findPresetByAlias(presets.Presets, "smoke-preset")
		if !found || preset.Name != "Smoke Preset Updated" || preset.SystemPrompt != "updated prompt" {
			t.Fatalf("unexpected updated preset payload: %+v", preset)
		}

		mustJSONRequest(t, client, http.MethodDelete, baseURL+"/presets/", map[string]any{
			"id": preset.ID,
		}, http.StatusOK)

		presetsResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/presets/", nil, http.StatusOK)
		presets = decodeResponseData[presetsResponse](t, presetsResp)
		if _, found := findPresetByAlias(presets.Presets, "smoke-preset"); found {
			t.Fatalf("expected smoke preset to be deleted, got %+v", presets.Presets)
		}
	})

	t.Run("folders explore and scans", func(t *testing.T) {
		mustJSONRequest(t, client, http.MethodPost, baseURL+"/folders/", map[string]any{
			"alias": "Smoke Folder",
			"path":  mediaDir,
		}, http.StatusOK)

		foldersResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/folders/", nil, http.StatusOK)
		folders := decodeResponseData[[]models.FolderConfig](t, foldersResp)
		if len(folders) != 1 || folders[0].Path != mediaDir {
			t.Fatalf("unexpected folders payload: %+v", folders)
		}

		videosResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/folders/scan/videos?path="+url.QueryEscape(mediaDir), nil, http.StatusOK)
		videos := decodeResponseData[[]string](t, videosResp)
		if len(videos) != 1 || videos[0] != "episode01.mkv" {
			t.Fatalf("unexpected video scan payload: %+v", videos)
		}

		subtitlesResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/folders/scan/subtitles?path="+url.QueryEscape(mediaDir), nil, http.StatusOK)
		subtitles := decodeResponseData[[]string](t, subtitlesResp)
		if len(subtitles) != 1 || subtitles[0] != "episode01.srt" {
			t.Fatalf("unexpected subtitle scan payload: %+v", subtitles)
		}

		exploreResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/folders/explore?path="+url.QueryEscape(mediaDir), nil, http.StatusOK)
		explore := decodeResponseData[models.ExploreResponse](t, exploreResp)
		if explore.FolderName != filepath.Base(mediaDir) || len(explore.Items) < 3 {
			t.Fatalf("unexpected explore payload: %+v", explore)
		}

		mustJSONRequest(t, client, http.MethodDelete, baseURL+"/folders/", map[string]any{
			"id": folders[0].ID,
		}, http.StatusOK)

		foldersResp = mustJSONRequest(t, client, http.MethodGet, baseURL+"/folders/", nil, http.StatusOK)
		folders = decodeResponseData[[]models.FolderConfig](t, foldersResp)
		if len(folders) != 0 {
			t.Fatalf("expected no folders after delete, got %+v", folders)
		}
	})

	t.Run("logs", func(t *testing.T) {
		logService := services.NewLogService(database)
		if err := logService.CreateLog("INFO", "smoke", "smoke log entry", map[string]any{"origin": "test"}); err != nil {
			t.Fatalf("seeding log failed: %v", err)
		}

		logsResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/logs/?module=smoke&limit=10&page=1", nil, http.StatusOK)
		logs := decodeResponseData[logsResponse](t, logsResp)
		if logs.Total == 0 || len(logs.Logs) == 0 {
			t.Fatalf("expected smoke logs in response, got %+v", logs)
		}
		if logs.Logs[0].Module != "smoke" || logs.Logs[0].Metadata["origin"] != "test" {
			t.Fatalf("unexpected log payload: %+v", logs.Logs[0])
		}
	})

	t.Run("jobs", func(t *testing.T) {
		jobService := services.NewJobService(database)
		if err := jobService.CreateJob("smoke-job", "/tmp/episode01.srt", "en", "anime", "smoke/model"); err != nil {
			t.Fatalf("creating smoke job failed: %v", err)
		}
		if err := jobService.UpdateTotalLines("smoke-job", 10); err != nil {
			t.Fatalf("updating total lines failed: %v", err)
		}
		if err := jobService.SetCachedLines("smoke-job", 3); err != nil {
			t.Fatalf("setting cached lines failed: %v", err)
		}
		if err := jobService.IncrementProgress("smoke-job", 7, 120, 45, 0.1234); err != nil {
			t.Fatalf("incrementing job progress failed: %v", err)
		}
		if err := jobService.UpdateStatus("smoke-job", "completed", ""); err != nil {
			t.Fatalf("updating job status failed: %v", err)
		}

		jobsResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/jobs/?limit=10&page=1", nil, http.StatusOK)
		jobs := decodeResponseData[jobsResponse](t, jobsResp)
		if jobs.Total == 0 || len(jobs.Jobs) == 0 {
			t.Fatalf("expected smoke jobs in response, got %+v", jobs)
		}

		jobResp := mustJSONRequest(t, client, http.MethodGet, baseURL+"/jobs/smoke-job", nil, http.StatusOK)
		job := decodeResponseData[models.TranslationJob](t, jobResp)
		if job.Status != "completed" || job.TotalLines != 10 || job.CachedLines != 3 || job.ProcessedLines != 7 {
			t.Fatalf("unexpected job payload: %+v", job)
		}
	})
}

func openSQLiteTestDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	database, err := db.InitializeSQLite(dsn)
	if err != nil {
		t.Fatalf("failed to initialize sqlite test database: %v", err)
	}

	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Fatalf("failed to close sqlite test database: %v", err)
		}
	})

	return database
}

func newSQLiteTestServer(t *testing.T, database *sql.DB, secretKey string) *httptest.Server {
	t.Helper()

	router := chi.NewRouter()
	router.Mount("/api/v1/", routes.APIRoutes(database, secretKey))

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server
}

func mustJSONRequest(t *testing.T, client *http.Client, method, rawURL string, payload any, expectedStatus int) apiResponse {
	t.Helper()

	var body io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal request payload: %v", err)
		}
		body = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status %d, got %d with body %s", expectedStatus, resp.StatusCode, string(responseBody))
	}

	var apiResp apiResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		t.Fatalf("failed to decode api response: %v; body=%s", err, string(responseBody))
	}

	if apiResp.Status != "success" {
		t.Fatalf("expected success response, got %+v", apiResp)
	}

	return apiResp
}

func decodeResponseData[T any](t *testing.T, response apiResponse) T {
	t.Helper()

	var data T
	if len(response.Data) == 0 || string(response.Data) == "null" {
		return data
	}

	if err := json.Unmarshal(response.Data, &data); err != nil {
		t.Fatalf("failed to decode response data: %v; raw=%s", err, string(response.Data))
	}

	return data
}

func mustMkdirAll(t *testing.T, dirPath string) {
	t.Helper()
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		t.Fatalf("failed to create directory %q: %v", dirPath, err)
	}
}

func mustWriteFile(t *testing.T, filePath, content string) {
	t.Helper()
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file %q: %v", filePath, err)
	}
}

func hasLanguage(languages []models.Language, code, name string) bool {
	for _, language := range languages {
		if language.Code == code && language.Name == name {
			return true
		}
	}

	return false
}

func findPresetByAlias(presets []models.TranslationPreset, alias string) (models.TranslationPreset, bool) {
	for _, preset := range presets {
		if preset.Alias == alias {
			return preset, true
		}
	}

	return models.TranslationPreset{}, false
}
