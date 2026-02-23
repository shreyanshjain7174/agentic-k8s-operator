// Auto-generated for phase1-012
// Task: Add monitoring dashboard showing agent activity and consensus metrics

```go
// types.go
package main

import "time"

// AgentStatus represents the current state of an agent
type AgentStatus struct {
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CurrentTask  string    `json:"current_task"`
	LastActivity time.Time `json:"last_activity"`
	Score        float64   `json:"score"`
	Duration     int64     `json:"duration"`
	HealthStatus string    `json:"health_status"`
}

// ConsensusMetrics represents voting results
type ConsensusMetrics struct {
	ProposalID    string             `json:"proposal_id"`
	AgentScores   map[string]float64 `json:"agent_scores"`
	AverageScore  float64            `json:"average_score"`
	Decision      string             `json:"decision"`
	Timestamp     time.Time          `json:"timestamp"`
	ExecutionTime int64              `json:"execution_time"`
}

// QuotaStatus represents API quota usage
type QuotaStatus struct {
	TotalCalls        int            `json:"total_calls"`
	RemainingCalls    int            `json:"remaining_calls"`
	WindowResetAt     time.Time      `json:"window_reset_at"`
	ModelDistribution map[string]int `json:"model_distribution"`
	UsagePercent      float64        `json:"usage_percent"`
}

// LoopMetrics represents a single execution loop
type LoopMetrics struct {
	LoopID             string    `json:"loop_id"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Duration           int64     `json:"duration"`
	ProposalsGenerated int       `json:"proposals_generated"`
	SuccessCount       int       `json:"success_count"`
	FailureCount       int       `json:"failure_count"`
	BottleneckAgent    string    `json:"bottleneck_agent"`
}

// MetricsSnapshot represents a point-in-time snapshot
type MetricsSnapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Agents    []AgentStatus    `json:"agents"`
	Consensus ConsensusMetrics `json:"consensus"`
	Quota     QuotaStatus      `json:"quota"`
	Loop      LoopMetrics      `json:"loop"`
}

// Alert represents a system alert
type Alert struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	AgentName string    `json:"agent_name,omitempty"`
}

// CurrentStateResponse is the API response for current state
type CurrentStateResponse struct {
	Snapshot MetricsSnapshot `json:"snapshot"`
	Alerts   []Alert         `json:"alerts"`
	Uptime   int64           `json:"uptime"`
}

// TimeRange represents a time window
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// HistorySummary provides aggregated statistics
type HistorySummary struct {
	TotalProposals   int     `json:"total_proposals"`
	ApprovalRate     float64 `json:"approval_rate"`
	AvgConsensusTime int64   `json:"avg_consensus_time"`
	MostActiveAgent  string  `json:"most_active_agent"`
}

// HistoryResponse is the API response for historical data
type HistoryResponse struct {
	TimeRange TimeRange         `json:"time_range"`
	Snapshots []MetricsSnapshot `json:"snapshots"`
	Summary   HistorySummary    `json:"summary"`
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// MetricsStore interface for storage operations
type MetricsStore interface {
	GetCurrent() (*MetricsSnapshot, error)
	GetHistory(start, end time.Time) ([]MetricsSnapshot, error)
	AddSnapshot(snapshot MetricsSnapshot) error
	GetQuota() (*QuotaStatus, error)
	Cleanup(retentionDays int) error
}

// AlertEngine interface for alert generation
type AlertEngine interface {
	CheckAlerts(snapshot MetricsSnapshot) []Alert
	GetActiveAlerts() []Alert
}

// WebSocketHub interface for real-time communication
type WebSocketHub interface {
	Broadcast(msg WSMessage)
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	Run()
}
```

```go
// storage.go
package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileMetricsStore implements MetricsStore using local files
type FileMetricsStore struct {
	monitoringDir string
	cache         *MetricsSnapshot
	mu            sync.RWMutex
}

// NewFileMetricsStore creates a new file-based metrics store
func NewFileMetricsStore(dir string) *FileMetricsStore {
	return &FileMetricsStore{
		monitoringDir: dir,
	}
}

// GetCurrent returns the latest metrics snapshot
func (s *FileMetricsStore) GetCurrent() (*MetricsSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(filepath.Join(s.monitoringDir, "metrics.json"))
	if err != nil {
		return nil, err
	}

	var snapshot MetricsSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// GetHistory returns snapshots within a time range
func (s *FileMetricsStore) GetHistory(start, end time.Time) ([]MetricsSnapshot, error) {
	var snapshots []MetricsSnapshot

	current := start
	for current.Before(end) || current.Equal(end) {
		filename := filepath.Join(
			s.monitoringDir,
			"history",
			current.Format("2006-01-02")+".jsonl",
		)

		if data, err := os.ReadFile(filename); err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}

				var snap MetricsSnapshot
				if err := json.Unmarshal([]byte(line), &snap); err == nil {
					if (snap.Timestamp.After(start) || snap.Timestamp.Equal(start)) &&
						(snap.Timestamp.Before(end) || snap.Timestamp.Equal(end)) {
						snapshots = append(snapshots, snap)
					}
				}
			}
		}

		current = current.Add(24 * time.Hour)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.Before(snapshots[j].Timestamp)
	})

	return snapshots, nil
}

// AddSnapshot adds a new metrics snapshot
func (s *FileMetricsStore) AddSnapshot(snapshot MetricsSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Write to current metrics file
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(s.monitoringDir, "metrics.json"), data, 0644); err != nil {
		return err
	}

	// Append to history
	historyDir := filepath.Join(s.monitoringDir, "history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return err
	}

	historyFile := filepath.Join(historyDir, time.Now().Format("2006-01-02")+".jsonl")
	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, _ := json.Marshal(snapshot)
	if _, err := f.Write(append(line, '\n')); err != nil {
		return err
	}

	s.cache = &snapshot
	return nil
}

// GetQuota returns current quota status
func (s *FileMetricsStore) GetQuota() (*QuotaStatus, error) {
	snapshot, err := s.GetCurrent()
	if err != nil {
		return nil, err
	}
	return &snapshot.Quota, nil
}

// Cleanup removes old history files
func (s *FileMetricsStore) Cleanup(retentionDays int) error {
	historyDir := filepath.Join(s.monitoringDir, "history")
	cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		dateStr := strings.TrimSuffix(entry.Name(), ".jsonl")
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			os.Remove(filepath.Join(historyDir, entry.Name()))
		}
	}

	return nil
}
```

```go
// websocket.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client represents a WebSocket client
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan WSMessage
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WSMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan WSMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msg WSMessage) {
	h.broadcast <- msg
}

// RegisterClient adds a new client
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes a client
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, _ := json.Marshal(message)
			w.Write(data)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from clients
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan WSMessage, 256),
	}

	hub.RegisterClient(client)

	go client.writePump()
	go client.readPump()
}
```

```go
// alerts.go
package main

import (
	"sync"
	"time"
)

// SimpleAlertEngine implements AlertEngine
type SimpleAlertEngine struct {
	alerts []Alert
	mu     sync.RWMutex
}

// NewAlertEngine creates a new alert engine
func NewAlertEngine() *SimpleAlertEngine {
	return &SimpleAlertEngine{
		alerts: make([]Alert, 0),
	}
}

// CheckAlerts analyzes a snapshot and generates alerts
func (e *SimpleAlertEngine) CheckAlerts(snapshot MetricsSnapshot) []Alert {
	e.mu.Lock()
	defer e.mu.Unlock()

	var newAlerts []Alert

	// Check quota usage
	if snapshot.Quota.UsagePercent > 90 {
		newAlerts = append(newAlerts, Alert{
			Level:     "critical",
			Message:   "API quota usage above 90%",
			Timestamp: time.Now(),
		})
	} else if snapshot.Quota.UsagePercent > 80 {
		newAlerts = append(newAlerts, Alert{
			Level:     "warning",
			Message:   "API quota usage above 80%",
			Timestamp: time.Now(),
		})
	}

	// Check agent health
	for _, agent := range snapshot.Agents {
		if agent.HealthStatus == "down" {
			newAlerts = append(newAlerts, Alert{
				Level:     "critical",
				Message:   agent.Name + " agent is down",
				Timestamp: time.Now(),
				AgentName: agent.Name,
			})
		} else if agent.HealthStatus == "degraded" {
			newAlerts = append(newAlerts, Alert{
				Level:     "warning",
				Message:   agent.Name + " agent performance degraded",
				Timestamp: time.Now(),
				AgentName: agent.Name,
			})
		}

		// Check for stale agents (no activity in 10 minutes)
		if time.Since(agent.LastActivity) > 10*time.Minute {
			newAlerts = append(newAlerts, Alert{
				Level:     "warning",
				Message:   agent.Name + " agent inactive for >10 minutes",
				Timestamp: time.Now(),
				AgentName: agent.Name,
			})
		}
	}

	// Check consensus score
	if snapshot.Consensus.AverageScore < 60 {
		newAlerts = append(newAlerts, Alert{
			Level:     "warning",
			Message:   "Low consensus score: " + snapshot.Consensus.Decision,
			Timestamp: time.Now(),
		})
	}

	e.alerts = newAlerts
	return newAlerts
}

// GetActiveAlerts returns current active alerts
func (e *SimpleAlertEngine) GetActiveAlerts() []Alert {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return append([]Alert{}, e.alerts...)
}
```

```go
// handlers.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// APIServer handles HTTP requests
type APIServer struct {
	store  MetricsStore
	alerts AlertEngine
	hub    WebSocketHub
}

// NewAPIServer creates a new API server
func NewAPIServer(store MetricsStore, alerts AlertEngine, hub WebSocketHub) *APIServer {
	return &APIServer{
		store:  store,
		alerts: alerts,
		hub:    hub,
	}
}

// HandleGetCurrent handles GET /api/current
func (s *APIServer) HandleGetCurrent(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.store.GetCurrent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alerts := s.alerts.CheckAlerts(*snapshot)
	uptime := int64(time.Since(snapshot.Loop.StartTime).Seconds())

	response := CurrentStateResponse{
		Snapshot: *snapshot,
		Alerts:   alerts,
		Uptime:   uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetHistory handles GET /api/history
func (s *APIServer) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	hoursStr := r.URL.Query().Get("hours")
	hours, _ := strconv.Atoi(hoursStr)
	if hours == 0 {
		hours = 24
	}

	end := time.Now()
	start := end.Add(-time.Duration(hours) * time.Hour)

	snapshots, err := s.store.GetHistory(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	summary := calculateSummary(snapshots)

	response := HistoryResponse{
		TimeRange: TimeRange{Start: start, End: end},
		Snapshots: snapshots,
		Summary:   summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetQuota handles GET /api/quota
func (s *APIServer) HandleGetQuota(w http.ResponseWriter, r *http.Request) {
	quota, err := s.store.GetQuota()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quota)
}

// HandleExport handles GET /api/export
func (s *APIServer) HandleExport(w http.ResponseWriter, r *http.Request) {
	hoursStr := r.URL.Query().Get("hours")
	hours, _ := strconv.Atoi(hoursStr)
	if hours == 0 {
		hours = 24
	}

	end := time.Now()
	start := end.Add(-time.Duration(hours) * time.Hour)

	snapshots, err := s.store.GetHistory(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=metrics-export.json")
	json.NewEncoder(w).Encode(snapshots)
}

// HandleWebSocket handles WebSocket connections
func (s *APIServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ServeWs(s.hub.(*Hub), w, r)
}

// calculateSummary computes aggregate statistics
func calculateSummary(snapshots []MetricsSnapshot) HistorySummary {
	if len(snapshots) == 0 {
		return HistorySummary{}
	}

	totalProposals := 0
	approvedCount := 0
	totalConsensusTime := int64(0)
	agentActivity := make(map[string]int)

	for _, snap := range snapshots {
		totalProposals++
		if snap.Consensus.Decision == "APPROVE" {
			approvedCount++
		}
		totalConsensusTime += snap.Consensus.ExecutionTime

		for _, agent := range snap.Agents {
			agentActivity[agent.Name]++
		}
	}

	approvalRate := 0.0
	if totalProposals > 0 {
		approvalRate = float64(approvedCount) / float64(totalProposals) * 100
	}

	avgConsensusTime := int64(0)
	if totalProposals > 0 {
		avgConsensusTime = totalConsensusTime / int64(totalProposals)
	}

	mostActive := ""
	maxActivity := 0
	for agent, count := range agentActivity {
		if count > maxActivity {
			maxActivity = count
			mostActive = agent
		}
	}

	return HistorySummary{
		TotalProposals:   totalProposals,
		ApprovalRate:     approvalRate,
		AvgConsensusTime: avgConsensusTime,
		MostActiveAgent:  mostActive,
	}
}
```

```go
// watcher.go
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileWatcher monitors metrics file for changes
type FileWatcher struct {
	path         string
	lastModTime  time.Time
	pollInterval time.Duration
	hub          WebSocketHub
	store        MetricsStore
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(path string, hub WebSocketHub, store MetricsStore) *FileWatcher {
	return &FileWatcher{
		path:         path,
		pollInterval: 5 * time.Second,
		hub:          hub,
		store:        store,
	}
}

// Watch starts monitoring the file
func (w *FileWatcher) Watch() {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		info, err := os.Stat(w.path)
		if err != nil {
			log.Printf("Error stat file: %v", err)
			continue
		}

		modTime := info.ModTime()
		if modTime.After(w.lastModTime) {
			w.lastModTime = modTime
			w.onFileChanged()
		}
	}
}

// onFileChanged handles file modification events
func (w *FileWatcher) onFileChanged() {
	snapshot, err := w.store.GetCurrent()
	if err != nil {
		log.Printf("Error reading snapshot: %v", err)
		return
	}

	msg := WSMessage{
		Type:    "snapshot",
		Payload: snapshot,
	}

	w.hub.Broadcast(msg)
	log.Printf("Broadcast new snapshot: %s", snapshot.Timestamp)
}
```

```go
// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Setup monitoring directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}
	monitoringDir := filepath.Join(homeDir, "monitoring")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(monitoringDir, "history"), 0755); err != nil {
		log.Fatal("Error creating monitoring directory:", err)
	}

	// Initialize components
	store := NewFileMetricsStore(monitoringDir)
	alertEngine := NewAlertEngine()
	hub := NewHub()

	// Start WebSocket hub
	go hub.Run()

	// Start file watcher
	metricsFile := filepath.Join(monitoringDir, "metrics.json")
	watcher := NewFileWatcher(metricsFile, hub, store)
	go watcher.Watch()

	// Create API server
	server := NewAPIServer(store, alertEngine, hub)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/api/current", server.HandleGetCurrent).Methods("GET")
	r.HandleFunc("/api/history", server.HandleGetHistory).Methods("GET")
	r.HandleFunc("/api/quota", server.HandleGetQuota).Methods("GET")
	r.HandleFunc("/api/export", server.HandleExport).Methods("GET")
	r.HandleFunc("/ws", server.HandleWebSocket)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Dashboard server starting on :%s", port)
	log.Printf("Monitoring directory: %s", monitoringDir)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
```

```go
// go.mod
module github.com/clawdlinux/dashboard-server

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/rs/cors v1.10.1
)
```