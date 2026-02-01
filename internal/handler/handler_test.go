package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/willianbsanches13/pack-calculator/internal/storage"
)

func setupTestRouter() (*gin.Engine, *Handler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := storage.NewMemoryStorage()
	h := New(store)
	h.RegisterRoutes(r)
	return r, h
}

func TestHealth(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", resp["status"])
	}
}

func TestGetPackSizes(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/pack-sizes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp PackSizesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.PackSizes) != 5 {
		t.Errorf("expected 5 default pack sizes, got %d", len(resp.PackSizes))
	}
}

func TestSetPackSizes(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"pack_sizes": [100, 200, 300]}`
	req := httptest.NewRequest(http.MethodPut, "/api/pack-sizes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp PackSizesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.PackSizes) != 3 {
		t.Errorf("expected 3 pack sizes, got %d", len(resp.PackSizes))
	}
}

func TestSetPackSizesEmpty(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"pack_sizes": []}`
	req := httptest.NewRequest(http.MethodPut, "/api/pack-sizes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSetPackSizesInvalid(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"pack_sizes": [100, -50, 300]}`
	req := httptest.NewRequest(http.MethodPut, "/api/pack-sizes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCalculatePost(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"amount": 251}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.OrderAmount != 251 {
		t.Errorf("expected order_amount 251, got %d", resp.OrderAmount)
	}

	if resp.TotalItems != 500 {
		t.Errorf("expected total_items 500, got %d", resp.TotalItems)
	}

	if resp.TotalPacks != 1 {
		t.Errorf("expected total_packs 1, got %d", resp.TotalPacks)
	}
}

func TestCalculateGet(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/calculate?amount=501", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.TotalItems != 750 {
		t.Errorf("expected total_items 750, got %d", resp.TotalItems)
	}

	if resp.TotalPacks != 2 {
		t.Errorf("expected total_packs 2, got %d", resp.TotalPacks)
	}
}

func TestCalculateGetMissingAmount(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/calculate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCalculateGetInvalidAmount(t *testing.T) {
	r, _ := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/calculate?amount=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCalculatePostZeroAmount(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"amount": 0}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCalculateWithCustomPackSizes(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"amount": 100, "pack_sizes": [23, 31, 53]}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.TotalItems < 100 {
		t.Errorf("expected total_items >= 100, got %d", resp.TotalItems)
	}
}

func TestAddPackSize(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"size": 750}`
	req := httptest.NewRequest(http.MethodPost, "/api/pack-sizes/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp PackSizesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	found := false
	for _, s := range resp.PackSizes {
		if s == 750 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected pack size 750 to be added")
	}
}

func TestAddPackSizeDuplicate(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"size": 250}`
	req := httptest.NewRequest(http.MethodPost, "/api/pack-sizes/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

func TestAddPackSizeInvalid(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"size": -100}`
	req := httptest.NewRequest(http.MethodPost, "/api/pack-sizes/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRemovePackSize(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"size": 250}`
	req := httptest.NewRequest(http.MethodPost, "/api/pack-sizes/remove", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp PackSizesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	for _, s := range resp.PackSizes {
		if s == 250 {
			t.Error("expected pack size 250 to be removed")
		}
	}
}

func TestRemovePackSizeNotFound(t *testing.T) {
	r, _ := setupTestRouter()

	body := `{"size": 999}`
	req := httptest.NewRequest(http.MethodPost, "/api/pack-sizes/remove", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
