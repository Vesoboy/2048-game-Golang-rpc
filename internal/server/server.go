package server

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type GameServer struct {
	sessions map[string]*GameSession
	mu       sync.RWMutex
	logger   *zap.Logger
	upgrader websocket.Upgrader
}

func NewGameServer(logger *zap.Logger) *GameServer {
	return &GameServer{
		sessions: make(map[string]*GameSession),
		logger:   logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // В продакшене нужно настроить CORS
			},
		},
	}
}

func (s *GameServer) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/ws", s.handleWebSocket)
	// Удаляем эту строку, так как статические файлы теперь обрабатываются в main.go
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
}

func (s *GameServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	sessionID := uuid.New().String()
	session := &GameSession{
		ID:     sessionID,
		WSConn: conn,
		logger: s.logger,
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.sessions, sessionID)
		s.mu.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		if err := session.HandleRPC(message); err != nil {
			s.logger.Error("Failed to handle RPC", zap.Error(err))
			break
		}
	}
}
