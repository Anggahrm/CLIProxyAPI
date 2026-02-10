package usage

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// globalPersistenceManager is set when usage persistence is initialized.
var globalPersistenceManager *PersistenceManager

// SetGlobalPersistenceManager registers the active persistence manager.
func SetGlobalPersistenceManager(pm *PersistenceManager) { globalPersistenceManager = pm }

// GetGlobalPersistenceManager returns the active persistence manager (may be nil).
func GetGlobalPersistenceManager() *PersistenceManager { return globalPersistenceManager }

// UsagePersister defines the interface for persisting usage data.
type UsagePersister interface {
	PersistUsage(ctx context.Context, data []byte) error
	LoadUsage(ctx context.Context) ([]byte, error)
}

// PersistenceManager handles periodic flushing and startup restoration of usage statistics.
type PersistenceManager struct {
	stats     *RequestStatistics
	persister UsagePersister
	interval  time.Duration
	stopCh    chan struct{}
	stopped   chan struct{}
	mu        sync.Mutex
}

// NewPersistenceManager creates a new persistence manager.
// interval controls how often data is flushed to the store (default 5 minutes if zero).
func NewPersistenceManager(stats *RequestStatistics, persister UsagePersister, interval time.Duration) *PersistenceManager {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	return &PersistenceManager{
		stats:     stats,
		persister: persister,
		interval:  interval,
		stopCh:    make(chan struct{}),
		stopped:   make(chan struct{}),
	}
}

// Restore loads the previously persisted usage snapshot and merges it into the in-memory store.
func (pm *PersistenceManager) Restore(ctx context.Context) error {
	if pm == nil || pm.persister == nil || pm.stats == nil {
		return nil
	}
	data, err := pm.persister.LoadUsage(ctx)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	var snapshot StatisticsSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		log.WithError(err).Warn("usage persistence: failed to unmarshal stored snapshot, starting fresh")
		return nil
	}
	result := pm.stats.MergeSnapshot(snapshot)
	log.Infof("usage persistence: restored %d records (%d skipped)", result.Added, result.Skipped)
	return nil
}

// Start begins the periodic flush goroutine.
func (pm *PersistenceManager) Start() {
	if pm == nil {
		return
	}
	go pm.flushLoop()
}

// Stop terminates the periodic flush and performs a final save.
func (pm *PersistenceManager) Stop() {
	if pm == nil {
		return
	}
	pm.mu.Lock()
	select {
	case <-pm.stopCh:
		pm.mu.Unlock()
		return
	default:
	}
	close(pm.stopCh)
	pm.mu.Unlock()
	<-pm.stopped
}

// Flush persists the current snapshot to the store.
func (pm *PersistenceManager) Flush(ctx context.Context) error {
	if pm == nil || pm.persister == nil || pm.stats == nil {
		return nil
	}
	snapshot := pm.stats.Snapshot()
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return pm.persister.PersistUsage(ctx, data)
}

func (pm *PersistenceManager) flushLoop() {
	defer close(pm.stopped)
	ticker := time.NewTicker(pm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-pm.stopCh:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := pm.Flush(ctx); err != nil {
				log.WithError(err).Error("usage persistence: final flush failed")
			} else {
				log.Info("usage persistence: final flush completed")
			}
			cancel()
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := pm.Flush(ctx); err != nil {
				log.WithError(err).Warn("usage persistence: periodic flush failed")
			}
			cancel()
		}
	}
}
