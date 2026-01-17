package server

import (
	"fmt"
	"sync"

	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// BatchSender handles batched sending of scan results to Server.
// It accumulates items and sends them in batches to reduce HTTP overhead.
type BatchSender struct {
	client    *Client
	scanID    int
	targetID  int
	dataType  string // "subdomain", "website", "endpoint", "port"
	batchSize int

	mu      sync.Mutex
	batch   []any
	sent    int // total items sent
	batches int // total batches sent
}

// NewBatchSender creates a new batch sender
func NewBatchSender(client *Client, scanID, targetID int, dataType string, batchSize int) *BatchSender {
	if batchSize <= 0 {
		batchSize = 1000 // default batch size
	}
	return &BatchSender{
		client:    client,
		scanID:    scanID,
		targetID:  targetID,
		dataType:  dataType,
		batchSize: batchSize,
		batch:     make([]any, 0, batchSize),
	}
}

// Add adds an item to the batch. Automatically sends when batch is full.
func (s *BatchSender) Add(item any) error {
	s.mu.Lock()
	s.batch = append(s.batch, item)
	shouldSend := len(s.batch) >= s.batchSize
	s.mu.Unlock()

	if shouldSend {
		return s.sendBatch()
	}
	return nil
}

// Flush sends any remaining items in the batch
func (s *BatchSender) Flush() error {
	s.mu.Lock()
	if len(s.batch) == 0 {
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()

	return s.sendBatch()
}

// Stats returns the total items and batches sent
func (s *BatchSender) Stats() (items, batches int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sent, s.batches
}

// sendBatch sends the current batch to the server
func (s *BatchSender) sendBatch() error {
	s.mu.Lock()
	if len(s.batch) == 0 {
		s.mu.Unlock()
		return nil
	}

	// Copy batch and clear
	toSend := make([]any, len(s.batch))
	copy(toSend, s.batch)
	s.batch = s.batch[:0] // reset slice but keep capacity
	s.mu.Unlock()

	// Build URL and body based on data type (RESTful style)
	var url string
	var body map[string]any

	switch s.dataType {
	case "subdomain":
		url = fmt.Sprintf("%s/api/worker/scans/%d/subdomains/bulk-upsert", s.client.baseURL, s.scanID)
		body = map[string]any{
			"targetId":   s.targetID,
			"subdomains": toSend,
		}
	case "website":
		url = fmt.Sprintf("%s/api/worker/scans/%d/websites/bulk-upsert", s.client.baseURL, s.scanID)
		body = map[string]any{
			"targetId": s.targetID,
			"websites": toSend,
		}
	case "endpoint":
		url = fmt.Sprintf("%s/api/worker/scans/%d/endpoints/bulk-upsert", s.client.baseURL, s.scanID)
		body = map[string]any{
			"targetId":  s.targetID,
			"endpoints": toSend,
		}
	default:
		// Generic fallback
		url = fmt.Sprintf("%s/api/worker/scans/%d/%ss/bulk-upsert", s.client.baseURL, s.scanID, s.dataType)
		body = map[string]any{
			"targetId": s.targetID,
			"items":    toSend,
		}
	}

	if err := s.client.postWithRetry(url, body); err != nil {
		pkg.Logger.Error("Failed to send batch",
			zap.String("type", s.dataType),
			zap.Int("count", len(toSend)),
			zap.Error(err))
		return fmt.Errorf("failed to send %s batch: %w", s.dataType, err)
	}

	s.mu.Lock()
	s.sent += len(toSend)
	s.batches++
	s.mu.Unlock()

	pkg.Logger.Debug("Batch sent",
		zap.String("type", s.dataType),
		zap.Int("count", len(toSend)),
		zap.Int("totalSent", s.sent))

	return nil
}
