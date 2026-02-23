// Auto-generated tests for phase1-012

```go
// storage_test.go
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileMetricsStore_GetCurrent(t *testing.T) {
	tests := []struct {
		name        string
		setupData   *MetricsSnapshot
		expectError bool
	}{
		{
			name: "valid snapshot",
			setupData: &MetricsSnapshot{
				Timestamp: time.Now(),
				Agents: []AgentStatus{
					{Name: "agent1", Status: "active", Score: 85.5},
				},
				Consensus: ConsensusMetrics{ProposalID: "prop-1", AverageScore: 80.0},
				Quota:     QuotaStatus{TotalCalls: 100, RemainingCalls: 50},
				Loop:      LoopMetrics{LoopID: "loop-1"},
			},
			expectError: false,
		},
		{
			name:        "missing file",
			setupData:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			store := NewFileMetricsStore(tmpDir)

			if tt.setupData != nil {
				data, _ := json.Marshal(tt.setupData)
				os.WriteFile(filepath.Join(tmpDir, "metrics.json"), data, 0644)
			}

			result, err := store.GetCurrent()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Error("expected result but got nil")
			}
		})
	}
}

func TestFileMetricsStore_AddSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileMetricsStore(tmpDir)

	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Agents: []AgentStatus{
			{Name: "test-agent", Status: "active", Score: 90.0},
		},
		Consensus: ConsensusMetrics{ProposalID: "test-prop"},
		Quota:     QuotaStatus{TotalCalls: 200},
		Loop:      LoopMetrics{LoopID: "test-loop"},
	}

	err := store.AddSnapshot(snapshot)
	if err != nil {
		t.Fatalf("AddSnapshot failed: %v", err)
	}

	// Verify metrics.json was created
	if _, err := os.Stat(filepath.Join(tmpDir, "metrics.json")); os.IsNotExist(err) {
		t.Error("metrics.json was not created")
	}

	// Verify history file was created
	historyFile := filepath.Join(tmpDir, "history", time.Now().Format("2006-01-02")+".jsonl")
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Error("history file was not created")
	}

	// Verify data can be retrieved
	retrieved, err := store.GetCurrent()
	if err != nil {
		t.Fatalf("GetCurrent failed: %v", err)
	}
	if retrieved.Consensus.ProposalID != "test-prop" {
		t.Errorf("expected ProposalID 'test-prop', got '%s'", retrieved.Consensus.ProposalID)
	}
}

func TestFileMetricsStore_GetHistory(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileMetricsStore(tmpDir)

	// Create test snapshots across multiple days
	baseTime := time.Now().Add(-48 * time.Hour)
	snapshots := []MetricsSnapshot{
		{Timestamp: baseTime, Consensus: ConsensusMetrics{ProposalID: "prop-1"}},
		{Timestamp: baseTime.Add(24 * time.Hour), Consensus: ConsensusMetrics{ProposalID: "prop-2"}},
		{Timestamp: baseTime.Add(48 * time.Hour), Consensus: ConsensusMetrics{ProposalID: "prop-3"}},
	}

	for _, snap := range snapshots {
		store.AddSnapshot(snap)
	}

	tests := []struct {
		name          string
		start         time.Time
		end           time.Time
		expectedCount int
	}{
		{
			name:          "all snapshots",
			start:         baseTime.Add(-1 * time.Hour),
			end:           time.Now(),
			expectedCount: 3,
		},
		{
			name:          "first day only",
			start:         baseTime.Add(-1 * time.Hour),
			end:           baseTime.Add(12 * time.Hour),
			expectedCount: 1,
		},
		{
			name:          "no snapshots in range",
			start:         baseTime.Add(-72 * time.Hour),
			end:           baseTime.Add(-49 * time.Hour),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.GetHistory(tt.start, tt.end)
			if err != nil {
				t.Fatalf("GetHistory failed: %v", err)
			}
			if len(result) != tt.expectedCount {
				t.Errorf("expected %d snapshots, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestFileMetricsStore_GetQuota(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileMetricsStore(tmpDir)

	expectedQuota := QuotaStatus{
		TotalCalls:     100,
		RemainingCalls: 25,
		UsagePercent:   75.0,
	}

	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Quota:     expectedQuota,
	}

	store.AddSnapshot(snapshot)

	quota, err := store.GetQuota()
	if err != nil {
		t.Fatalf("GetQuota failed: %v", err)
	}
	if quota.TotalCalls != expectedQuota.TotalCalls {
		t.Errorf("expected TotalCalls %d, got %d", expectedQuota.TotalCalls, quota.TotalCalls)
	}
	if quota.RemainingCalls != expectedQuota.RemainingCalls {
		t.Errorf("expected RemainingCalls %d, got %d", expectedQuota.RemainingCalls, quota.RemainingCalls)
	}
}

func TestFileMetricsStore_Cleanup(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileMetricsStore(tmpDir)
	historyDir := filepath.Join(tmpDir, "history")
	os.MkdirAll(historyDir, 0755)

	// Create old and new files
	oldDate := time.Now().Add(-10 * 24 * time.Hour).Format("2006-01-02")
	newDate := time.Now().Format("2006-01-02")

	oldFile := filepath.Join(historyDir, oldDate+".jsonl")
	newFile := filepath.Join(historyDir, newDate+".jsonl")

	os.WriteFile(oldFile, []byte("test"), 0644)
	os.WriteFile(newFile, []byte("test"), 0644)

	// Cleanup with 7 day retention
	err := store.Cleanup(7)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Old file should be deleted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("old file was not deleted")
	}

	// New file should remain
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("new file was incorrectly deleted")
	}
}

func TestFileMetricsStore_Concurrency(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileMetricsStore(tmpDir)

	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Consensus: ConsensusMetrics{ProposalID: "concurrent-test"},
	}

	store.AddSnapshot(snapshot)

	// Concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := store.GetCurrent()
			if err != nil {
				t.Errorf("concurrent read failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
```

```go
// alerts_test.go
package main

import (
	"testing"
	"time"
)

func TestSimpleAlertEngine_CheckAlerts_QuotaUsage(t *testing.T) {
	tests := []struct {
		name          string
		usagePercent  float64
		expectedLevel string
		expectedCount int
	}{
		{
			name:          "critical quota usage",
			usagePercent:  95.0,
			expectedLevel: "critical",
			expectedCount: 1,
		},
		{
			name:          "warning quota usage",
			usagePercent:  85.0,
			expectedLevel: "warning",
			expectedCount: 1,
		},
		{
			name:          "normal quota usage",
			usagePercent:  50.0,
			expectedLevel: "",
			expectedCount: 0,
		},
		{
			name:          "exactly 90 percent",
			usagePercent:  90.0,
			expectedLevel: "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewAlertEngine()
			snapshot := MetricsSnapshot{
				Quota: QuotaStatus{UsagePercent: tt.usagePercent},
			}

			alerts := engine.CheckAlerts(snapshot)

			if len(alerts) != tt.expectedCount {
				t.Errorf("expected %d alerts, got %d", tt.expectedCount, len(alerts))
			}

			if tt.expectedCount > 0 && alerts[0].Level != tt.expectedLevel {
				t.Errorf("expected level '%s', got '%s'", tt.expectedLevel, alerts[0].Level)
			}
		})
	}
}

func TestSimpleAlertEngine_CheckAlerts_AgentHealth(t *testing.T) {
	tests := []struct {
		name          string
		healthStatus  string
		expectedLevel string
		shouldAlert   bool
	}{
		{
			name:          "agent down",
			healthStatus:  "down",
			expectedLevel: "critical",
			shouldAlert:   true,
		},
		{
			name:          "agent degraded",
			healthStatus:  "degraded",
			expectedLevel: "warning",
			shouldAlert:   true,
		},
		{
			name:          "agent healthy",
			healthStatus:  "healthy",
			expectedLevel: "",
			shouldAlert:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewAlertEngine()
			snapshot := MetricsSnapshot{
				Agents: []AgentStatus{
					{
						Name:         "test-agent",
						HealthStatus: tt.healthStatus,
						LastActivity: time.Now(),
					},
				},
			}

			alerts := engine.CheckAlerts(snapshot)

			foundAlert := false
			for _, alert := range alerts {
				if alert.AgentName == "test-agent" {
					foundAlert = true
					if alert.Level != tt.expectedLevel {
						t.Errorf("expected level '%s', got '%s'", tt.expectedLevel, alert.Level)
					}
				}
			}

			if tt.shouldAlert && !foundAlert {
				t.Error("expected alert but none found")
			}
			if !tt.shouldAlert && foundAlert {
				t.Error("unexpected alert generated")
			}
		})
	}
}

func TestSimpleAlertEngine_CheckAlerts_StaleAgent(t *testing.T) {
	engine := NewAlertEngine()

	tests := []struct {
		name          string
		lastActivity  time.Time
		shouldAlert   bool
	}{
		{
			name:         "active agent",
			lastActivity: time.Now().Add(-5 * time.Minute),
			shouldAlert:  false,
		},
		{
			name:         "stale agent",
			lastActivity: time.Now().Add(-15 * time.Minute),
			shouldAlert:  true,
		},
		{
			name:         "exactly 10 minutes",
			lastActivity: time.Now().Add(-10 * time.Minute),
			shouldAlert:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := MetricsSnapshot{
				Agents: []AgentStatus{
					{
						Name:         "test-agent",
						HealthStatus: "healthy",
						LastActivity: tt.lastActivity,
					},
				},
			}

			alerts := engine.CheckAlerts(snapshot)

			foundStaleAlert := false
			for _, alert := range alerts {
				if alert.AgentName == "test-agent" && alert.Message == "test-agent agent inactive for >10 minutes" {
					foundStaleAlert = true
				}
			}

			if tt.shouldAlert && !foundStaleAlert {
				t.Error("expected stale agent alert but none found")
			}
			if !tt.shouldAlert && foundStaleAlert {
				t.Error("unexpected stale agent alert generated")
			}
		})
	}
}

func TestSimpleAlertEngine_CheckAlerts_ConsensusScore(t *testing.T) {
	tests := []struct {
		name        string
		avgScore    float64
		shouldAlert bool
	}{
		{
			name:        "low consensus score",
			avgScore:    45.0,
			shouldAlert: true,
		},
		{
			name:        "borderline score",
			avgScore:    59.9,
			shouldAlert: true,
		},
		{
			name:        "acceptable score",
			avgScore:    60.0,
			shouldAlert: false,
		},
		{
			name:        "high score",
			avgScore:    95.0,
			shouldAlert: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewAlertEngine()
			snapshot := MetricsSnapshot{
				Consensus: ConsensusMetrics{
					AverageScore: tt.avgScore,
					Decision:     "APPROVE",
				},
			}

			alerts := engine.CheckAlerts(snapshot)

			foundConsensusAlert := false
			for _, alert := range alerts {
				if alert.Message == "Low consensus score: APPROVE" {
					foundConsensusAlert = true
				}
			}

			if tt.shouldAlert && !foundConsensusAlert {
				t.Error("expected consensus alert but none found")
			}
			if !tt.shouldAlert && foundConsensusAlert {
				t.Error("unexpected consensus alert generated")
			}
		})
	}
}

func TestSimpleAlertEngine_GetActiveAlerts(t *testing.T) {
	engine := NewAlertEngine()

	// Initially no alerts
	alerts := engine.GetActiveAlerts()
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}

	// Check alerts to generate some
	snapshot := MetricsSnapshot{
		Quota: QuotaStatus{UsagePercent: 95.0},
	}
	engine.CheckAlerts(snapshot)

	// Should now have active alerts
	alerts = engine.GetActiveAlerts()
	if len(alerts) == 0 {
		t.Error("expected active alerts but got none")
	}

	// Verify returned slice is a copy (mutation shouldn't affect internal state)
	alerts[0].Level = "modified"
	originalAlerts := engine.GetActiveAlerts()
	if originalAlerts[0].Level == "modified" {
		t.Error("GetActiveAlerts should return a copy, not the internal slice")
	}
}

func TestSimpleAlertEngine_MultipleAlertTypes(t *testing.T) {
	engine := NewAlertEngine()

	snapshot := MetricsSnapshot{
		Quota: QuotaStatus{UsagePercent: 95.0},
		Agents: []AgentStatus{
			{Name: "agent1", HealthStatus: "down", LastActivity: time.Now()},
			{Name: "agent2", HealthStatus: "degraded", LastActivity: time.Now()},
		},
		Consensus: ConsensusMetrics{AverageScore: 50.0, Decision: "REJECT"},
	}

	alerts := engine.CheckAlerts(snapshot)

	// Should have 4 alerts: 1 quota + 2 agent health + 1 consensus
	if len(alerts) != 4 {
		t.Errorf("expected 4 alerts, got %d", len(alerts))
	}

	// Verify alert types
	hasCritical := false
	hasWarning := false
	for _, alert := range alerts {
		if alert.Level == "critical" {
			hasCritical = true
		}
		if alert.Level == "warning" {
			hasWarning = true
		}
	}

	if !hasCritical {
		t.Error("expected at least one critical alert")
	}
	if !hasWarning {
		t.Error("expected at least one warning alert")
	}
}
```

```go
// handlers_test.go
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockMetricsStore struct {
	currentSnapshot *MetricsSnapshot
	historyData     []MetricsSnapshot
	quotaData       *QuotaStatus
	err             error
}

func (m *mockMetricsStore) GetCurrent() (*MetricsSnapshot, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.currentSnapshot, nil
}

func (m *mockMetricsStore) GetHistory(start, end time.Time) ([]MetricsSnapshot, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.historyData, nil
}

func (m *mockMetricsStore) AddSnapshot(snapshot MetricsSnapshot) error {
	return m.err
}

func (m *mockMetricsStore) GetQuota() (*QuotaStatus, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.quotaData, nil
}

func (m *mockMetricsStore) Cleanup(retentionDays int) error {
	return m.err
}

type mockAlertEngine struct {
	alerts []Alert
}

func (m *mockAlertEngine) CheckAlerts(snapshot MetricsSnapshot) []Alert {
	return m.alerts
}

func (m *mockAlertEngine) GetActiveAlerts() []Alert {
	return m.alerts
}

type mockWebSocketHub struct{}

func (m *mockWebSocketHub) Broadcast(msg WSMessage)              {}
func (m *mockWebSocketHub) RegisterClient(client *Client)        {}
func (m *mockWebSocketHub) UnregisterClient(client *Client)      {}
func (m *mockWebSocketHub) Run()                                 {}

func TestHandleGetCurrent_Success(t *testing.T) {
	now := time.Now()
	mockStore := &mockMetricsStore{
		currentSnapshot: &MetricsSnapshot{
			Timestamp: now,
			Loop:      LoopMetrics{StartTime: now.Add(-1 * time.Hour)},
		},
	}
	mockAlerts := &mockAlertEngine{
		alerts: []Alert{{Level: "warning", Message: "test alert"}},
	}
	mockHub := &mockWebSocketHub{}

	server := NewAPIServer(mockStore, mockAlerts, mockHub)

	req := httptest.NewRequest("GET", "/api/current", nil)
	w := httptest.NewRecorder()

	server.HandleGetCurrent(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response CurrentStateResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Alerts) != 1 {
		t.Errorf("expected 1 alert, got %d", len(response.Alerts))
	}

	if response.Uptime == 0 {
		t.Error("expected uptime > 0")
	}
}

func TestHandleGetCurrent_StoreError(t *testing.T) {
	mockStore := &mockMetricsStore{
		err: errors.New("store error"),
	}
	mockAlerts := &mockAlertEngine{}
	mockHub := &mockWebSocketHub{}

	server := NewAPIServer(mockStore, mockAlerts, mockHub)

	req := httptest.NewRequest("GET", "/api/current", nil)
	w := httptest.NewRecorder()

	server.HandleGetCurrent(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestHandleGetHistory_Success(t *testing.T) {
	mockSnapshots := []MetricsSnapshot{
		{Timestamp: time.Now(), Consensus: ConsensusMetrics{Decision: "APPROVE"}},
		{Timestamp: time.Now(), Consensus: ConsensusMetrics{Decision: "REJECT"}},
	}

	mockStore := &mockMetricsStore{
		historyData: mockSnapshots,
	}
	mockAlerts := &mockAlertEngine{}
	mockHub := &mockWebSocketHub{}

	server := NewAPIServer(mockStore, mockAlerts, mockHub)

	tests := []struct {
		name          string
		queryParam    string
		expectedHours int
	}{
		{
			name:          "default hours",
			queryParam:    "",
			expectedHours: 24,
		},
		{
			name:          "custom hours",
			queryParam:    "?hours=48",
			expectedHours: 48,
		},
		{
			name:          "invalid hours",
			queryParam:    "?hours=invalid",
			expectedHours: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/history"+tt.queryParam, nil)
			w := httptest.NewRecorder()

			server.HandleGetHistory(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			var response HistoryResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(response.Snapshots) != len(mockSnapshots) {
				t.Errorf("expected %d snapshots, got %d", len(mockSnapshots), len(response.Snapshots))
			}
		})
	}
}

func TestHandleGetQuota_Success(t *testing.T) {
	mockQuota := &QuotaStatus{
		TotalCalls:     100,
		RemainingCalls: 50,
		UsagePercent:   50.0,
	}

	mockStore := &mockMetricsStore{
		quotaData: mockQuota,
	}
	mockAlerts := &mockAlertEngine{}
	mockHub := &mockWebSocketHub{}

	server := NewAPIServer(mockStore, mockAlerts, mockHub)

	req := httptest.NewRequest("GET", "/api/quota", nil)
	w := httptest.NewRecorder()

	server.HandleGetQuota(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response QuotaStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.TotalCalls != mockQuota.TotalCalls {
		t.Errorf("expected TotalCalls %d, got %d", mockQuota.TotalCalls, response.TotalCalls)
	}
}

func TestHandleExport_Success(t *testing.T) {
	mockSnapshots := []MetricsSnapshot{
		{Timestamp: time.Now()},
		{Timestamp: time.Now()},
	}

	mockStore := &mockMetricsStore{
		historyData: mockSnapshots,
	}
	mockAlerts := &mockAlertEngine{}
	mockHub := &mockWebSocketHub{}

	server := NewAPIServer(mockStore, mockAlerts, mockHub)

	req := httptest.NewRequest("GET", "/api/export", nil)
	w := httptest.NewRecorder()

	server.HandleExport(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentDisposition := w.Header().Get("Content-Disposition")
	if contentDisposition != "attachment; filename=metrics-export.json" {
		t.Errorf("unexpected Content-Disposition: %s", contentDisposition)
	}

	var response []MetricsSnapshot
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != len(mockSnapshots) {
		t.Errorf("expected %d snapshots, got %d", len(mockSnapshots), len(response))
	}
}

func TestCalculateSummary(t *testing.T) {
	tests := []struct {
		name              string
		snapshots         []MetricsSnapshot
		expectedProposals int
		expectedApproval  float64
		expectedActive    string
	}{
		{
			name:              "empty snapshots",
			snapshots:         []MetricsSnapshot{},
			expectedProposals: 0,
			expectedApproval:  0.0,
			expectedActive:    "",
		},
		{
			name: "all approved",
			snapshots: []MetricsSnapshot{
				{Consensus: ConsensusMetrics{Decision: "APPROVE", ExecutionTime: 100}},
				{Consensus: ConsensusMetrics{Decision: "APPROVE", ExecutionTime: 200}},
			},
			expectedProposals: 2,
			expectedApproval:  100.0,
		},
		{
			name: "mixed decisions",
			snapshots: []MetricsSnapshot{
				{
					Consensus: ConsensusMetrics{Decision: "APPROVE", ExecutionTime: 100},
					Agents:    []AgentStatus{{Name: "agent1"}},
				},
				{
					Consensus: ConsensusMetrics{Decision: "REJECT", ExecutionTime: 200},
					Agents:    []AgentStatus{{Name: "agent1"}, {Name: "agent2"}},
				},
			},
			expectedProposals: 2,
			expectedApproval:  50.0,
			expectedActive:    "agent1",
		},
		{
			name: "zero execution time",
			snapshots: []MetricsSnapshot{
				{Consensus: ConsensusMetrics{Decision: "APPROVE", ExecutionTime: 0}},
			},
			expectedProposals: 1,
			expectedApproval:  100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := calculateSummary(tt.snapshots)

			if summary.TotalProposals != tt.expectedProposals {
				t.Errorf("expected %d proposals, got %d", tt.expectedProposals, summary.TotalProposals)
			}

			if summary.ApprovalRate != tt.expectedApproval {
				t.Errorf("expected approval rate %.2f, got %.2f", tt.expectedApproval, summary.ApprovalRate)
			}

			if tt.expectedActive != "" && summary.MostActiveAgent != tt.expectedActive {
				t.Errorf("expected most active agent '%s', got '%s'", tt.expectedActive, summary.MostActiveAgent)
			}
		})
	}
}

func TestCalculateSummary_AverageConsensusTime(t *testing.T) {
	snapshots := []MetricsSnapshot{
		{Consensus: ConsensusMetrics{ExecutionTime: 100}},
		{Consensus: ConsensusMetrics{ExecutionTime: 200}},
		{Consensus: ConsensusMetrics{ExecutionTime: 300}},
	}

	summary := calculateSummary(snapshots)

	expectedAvg := int64(200)
	if summary.AvgConsensusTime != expectedAvg {
		t.Errorf("expected avg consensus time %d, got %d", expectedAvg, summary.AvgConsensusTime)
	}
}
```

```go
// websocket_test.go
package main

import (
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()

	if hub.clients == nil {
		t.Error("clients map not initialized")
	}
	if hub.broadcast == nil {
		t.Error("broadcast channel not initialized")
	}
	if hub.register == nil {
		t.Error("register channel not initialized")
	}
	if hub.unregister == nil {
		t.Error("unregister channel not initialized")
	}
}

func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		hub:  hub,
		send: make(chan WSMessage, 256),
	}

	hub.RegisterClient(client)

	// Give time for registration
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	if !hub.clients[client] {
		t.Error("client was not registered")
	}
	hub.mu.RUnlock()
}

func TestHub_UnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		hub:  hub,
		send: make(chan WSMessage, 256),
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	hub.UnregisterClient(client)
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	if hub.clients[client] {
		t.Error("client was not unregistered")
	}
	hub.mu.RUnlock()
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{
		hub:  hub,
		send: make(chan WSMessage, 256),
	}
	client2 := &Client{
		hub:  hub,
		send: make(chan WSMessage, 256),
	}

	hub.RegisterClient(client1)
	hub.RegisterClient(client2)
	time.Sleep(10 * time.Millisecond)

	testMsg := WSMessage{
		Type:    "test",
		Payload: "test data",
	}

	hub.Broadcast(testMsg)
	time.Sleep(10 * time.Millisecond)

	// Check both clients received the message
	select {
	case msg := <-client1.send:
		if msg.Type != "test" {
			t.Errorf("client1 received wrong message type: %s", msg.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 did not receive broadcast message")
	}

	select {
	case msg := <-client2.send:
		if msg.Type != "test" {
			t.Errorf("client2 received wrong message type: %s", msg.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client2 did not receive broadcast message")
	}
}

func TestHub_BroadcastToBlockedClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create client with zero buffer
	client := &Client{
		hub:  hub,
		send: make(chan WSMessage),
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	// Fill the channel
	testMsg := WSMessage{Type: "test"}

	// This should not block even if client can't receive
	done := make(chan bool)
	go func() {
		hub.Broadcast(testMsg)
		done <- true
	}()

	select {
	case <-done:
		// Success - broadcast didn't block
	case <-time.After(100 * time.Millisecond):
		t.Error("broadcast blocked on full client channel")
	}
}

func TestHub_MultipleClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	numClients := 5
	clients := make([]*Client, numClients)

	for i := 0; i < numClients; i++ {
		clients[i] = &Client{
			hub:  hub,
			send: make(chan WSMessage, 256),
		}
		hub.RegisterClient(clients[i])
	}

	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	if len(hub.clients) != numClients {
		t.Errorf("expected %d clients, got %d", numClients, len(hub.clients))
	}
	hub.mu.RUnlock()

	// Unregister one client
	hub.UnregisterClient(clients[0])
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	if len(hub.clients) != numClients-1 {
		t.Errorf("expected %d clients after unregister, got %d", numClients-1, len(hub.clients))
	}
	hub.mu.RUnlock()
}

func TestHub_ConcurrentOperations(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	done := make(chan bool)
	numGoroutines := 10

	// Concurrent registrations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			client := &Client{
				hub:  hub,
				send: make(chan WSMessage, 256),
			}
			hub.RegisterClient(client)
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	clientCount := len(hub.clients)
	hub.mu.RUnlock()

	if clientCount != numGoroutines {
		t.Errorf("expected %d clients, got %d", numGoroutines, clientCount)
	}
}
```

```go
// watcher_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFileWatcher(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.json")
	hub := &mockWebSocketHub{}
	store := &mockMetricsStore{}

	watcher := NewFileWatcher(testPath, hub, store)

	if watcher.path != testPath {
		t.Errorf("expected path %s, got %s", testPath, watcher.path)
	}
	if watcher.pollInterval != 5*time.Second {
		t.Errorf("expected poll interval 5s, got %v", watcher.pollInterval)
	}
	if watcher.hub == nil {
		t.Error("hub not set")
	}
	if watcher.store == nil {
		t.Error("store not set")
	}
}

func TestFileWatcher_DetectsFileChange(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.json")

	// Create initial file
	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Consensus: ConsensusMetrics{ProposalID: "test"},
	}

	store := &mockMetricsStore{
		currentSnapshot: &snapshot,
	}

	broadcastReceived := make(chan bool, 1)
	hub := &mockWebSocketHubWithCallback{
		onBroadcast: func(msg WSMessage) {
			broadcastReceived <- true
		},
	}

	watcher := NewFileWatcher(testPath, hub, store)
	watcher.pollInterval = 100 * time.Millisecond

	// Create initial file
	os.WriteFile(testPath, []byte("initial"), 0644)
	time.Sleep(50 * time.Millisecond)

	// Start watching in background
	go watcher.Watch()

	// Modify file
	time.Sleep(150 * time.Millisecond)
	os.WriteFile(testPath, []byte("modified"), 0644)

	// Wait for broadcast
	select {
	case <-broadcastReceived:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Error("file change was not detected")
	}
}

func TestFileWatcher_IgnoresUnchangedFile(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.json")

	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
	}

	store := &mockMetricsStore{
		currentSnapshot: &snapshot,
	}

	broadcastCount := 0
	hub := &mockWebSocketHubWithCallback{
		onBroadcast: func(msg WSMessage) {
			broadcastCount++
		},
	}

	// Create file
	os.WriteFile(testPath, []byte("test"), 0644)

	watcher := NewFileWatcher(testPath, hub, store)
	watcher.pollInterval = 50 * time.Millisecond

	go watcher.Watch()

	// Wait for several poll cycles
	time.Sleep(200 * time.Millisecond)

	if broadcastCount > 0 {
		t.Errorf("expected 0 broadcasts for unchanged file, got %d", broadcastCount)
	}
}

func TestFileWatcher_HandlesNonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "nonexistent.json")

	store := &mockMetricsStore{}
	hub := &mockWebSocketHub{}

	watcher := NewFileWatcher(testPath, hub, store)
	watcher.pollInterval = 50 * time.Millisecond

	// Should not panic
	done := make(chan bool)
	go func() {
		watcher.Watch()
		done <- true
	}()

	time.Sleep(150 * time.Millisecond)

	// Watcher should continue running even with missing file
	select {
	case <-done:
		t.Error("watcher stopped unexpectedly")
	default:
		// Success - still running
	}
}

type mockWebSocketHubWithCallback struct {
	onBroadcast func(msg WSMessage)
}

func (m *mockWebSocketHubWithCallback) Broadcast(msg WSMessage) {
	if m.onBroadcast != nil {
		m.onBroadcast(msg)
	}
}

func (m *mockWebSocketHubWithCallback) RegisterClient(client *Client)   {}
func (m *mockWebSocketHubWithCallback) UnregisterClient(client *Client) {}
func (m *mockWebSocketHubWithCallback) Run()                            {}
```