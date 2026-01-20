package server

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// BatchSender handles batched sending of scan results to Server.
// It accumulates items and sends them in batches to reduce HTTP overhead.
type BatchSender struct {
	ctx       context.Context
	client    ServerClient
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
func NewBatchSender(ctx context.Context, client ServerClient, scanID, targetID int, dataType string, batchSize int) *BatchSender {
	if batchSize <= 0 {
		batchSize = 1000 // default batch size
	}
	return &BatchSender{
		ctx:       ctx,
		client:    client,
		scanID:    scanID,
		targetID:  targetID,
		dataType:  dataType,
		batchSize: batchSize,
		batch:     make([]any, 0, batchSize),
	}
}

// Add adds an item to the batch. Automatically sends when batch is full.
// Returns context.Canceled or context.DeadlineExceeded if context is done.
func (s *BatchSender) Add(item any) error {
	// Check context before processing
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}

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
	// Check context before sending
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}

	s.mu.Lock()
	if len(s.batch) == 0 {
		s.mu.Unlock()
		return nil
	}

	// Copy batch and clear so new items can be queued while sending
	toSend := make([]any, len(s.batch))
	copy(toSend, s.batch)
	s.batch = s.batch[:0]
	s.mu.Unlock()

	if err := s.client.PostBatch(s.ctx, s.scanID, s.targetID, s.dataType, toSend); err != nil {
		// Check if it's a non-retryable error (4xx)
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && !httpErr.IsRetryable() {
			pkg.Logger.Error("Non-retryable error sending batch (data validation issue)",
				zap.String("type", s.dataType),
				zap.Int("count", len(toSend)),
				zap.Int("statusCode", httpErr.StatusCode),
				zap.String("response", httpErr.Body))
		} else {
			pkg.Logger.Error("Failed to send batch after retries",
				zap.String("type", s.dataType),
				zap.Int("count", len(toSend)),
				zap.Error(err))
		}
		// Re-queue batch for retry on next Flush/send, preserving new items
		s.mu.Lock()
		s.batch = append(toSend, s.batch...)
		s.mu.Unlock()

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
