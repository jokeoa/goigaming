package game

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

type HubManager struct {
	mu          sync.RWMutex
	hubs        map[uuid.UUID]*TableHub
	turnTimeout time.Duration
	broadcaster ports.Broadcaster
	walletSvc   ports.WalletService
	rngSvc      ports.RNGService
	handRepo    ports.PokerHandRepository
	playerRepo  ports.PokerPlayerRepository
	logger      *slog.Logger
	baseCtx     context.Context
}

func NewHubManager(
	baseCtx context.Context,
	turnTimeout time.Duration,
	broadcaster ports.Broadcaster,
	walletSvc ports.WalletService,
	rngSvc ports.RNGService,
	handRepo ports.PokerHandRepository,
	playerRepo ports.PokerPlayerRepository,
	logger *slog.Logger,
) *HubManager {
	return &HubManager{
		hubs:        make(map[uuid.UUID]*TableHub),
		turnTimeout: turnTimeout,
		broadcaster: broadcaster,
		walletSvc:   walletSvc,
		rngSvc:      rngSvc,
		handRepo:    handRepo,
		playerRepo:  playerRepo,
		logger:      logger,
		baseCtx:     baseCtx,
	}
}

func (m *HubManager) GetOrCreateHub(ctx context.Context, table domain.PokerTable) *TableHub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.hubs[table.ID]; ok {
		return hub
	}

	hub := NewTableHub(
		table,
		m.turnTimeout,
		m.broadcaster,
		m.walletSvc,
		m.rngSvc,
		m.handRepo,
		m.playerRepo,
		m.logger,
	)

	m.hubs[table.ID] = hub

	runCtx := m.baseCtx
	if runCtx == nil {
		runCtx = ctx
	}
	go func() {
		hub.Run(runCtx)
		m.mu.Lock()
		delete(m.hubs, table.ID)
		m.mu.Unlock()
	}()

	return hub
}

func (m *HubManager) GetHub(tableID uuid.UUID) *TableHub {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hubs[tableID]
}

func (m *HubManager) RemoveHub(tableID uuid.UUID) {
	m.mu.Lock()
	hub, ok := m.hubs[tableID]
	if ok {
		delete(m.hubs, tableID)
	}
	m.mu.Unlock()

	if hub != nil {
		resultCh := make(chan HubResult, 1)
		if err := hub.Send(HubEvent{Type: EventShutdown, ResultCh: resultCh}); err == nil {
			<-resultCh
		}
	}
}

func (m *HubManager) ShutdownAll() {
	m.mu.Lock()
	hubsCopy := make(map[uuid.UUID]*TableHub, len(m.hubs))
	for id, hub := range m.hubs {
		hubsCopy[id] = hub
	}
	m.hubs = make(map[uuid.UUID]*TableHub)
	m.mu.Unlock()

	for _, hub := range hubsCopy {
		resultCh := make(chan HubResult, 1)
		if err := hub.Send(HubEvent{Type: EventShutdown, ResultCh: resultCh}); err == nil {
			<-resultCh
		}
	}
}
