package server

import (
	"encoding/json"
	"game-2048/internal/game"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type GameSession struct {
	ID      string     `json:"id"`
	Game    *game.Game `json:"game"`
	WSConn  *websocket.Conn
	mu      sync.Mutex
	logger  *zap.Logger
}

type RPCRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     int            `json:"id"`
}

type RPCResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  *RPCError  `json:"error,omitempty"`
	ID     int        `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MoveParams struct {
	Direction string `json:"direction"`
}

func (s *GameSession) HandleRPC(message []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var req RPCRequest
	if err := json.Unmarshal(message, &req); err != nil {
		return s.sendError(req.ID, -32700, "Parse error")
	}

	var response RPCResponse
	response.ID = req.ID

	switch req.Method {
	case "newGame":
		s.Game = game.NewGame()
		response.Result = s.Game

	case "move":
		var params MoveParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return s.sendError(req.ID, -32602, "Invalid params")
		}

		var dir game.Direction
		switch params.Direction {
		case "up":
			dir = game.Up
		case "right":
			dir = game.Right
		case "down":
			dir = game.Down
		case "left":
			dir = game.Left
		default:
			return s.sendError(req.ID, -32602, "Invalid direction")
		}

		moved := s.Game.Move(dir)
		response.Result = map[string]interface{}{
			"moved":     moved,
			"grid":      s.Game.Grid,
			"score":     s.Game.Score,
			"bestScore": s.Game.BestScore,
			"gameOver":  s.Game.GameOver,
		}

	default:
		return s.sendError(req.ID, -32601, "Method not found")
	}

	return s.sendResponse(&response)
}

func (s *GameSession) sendResponse(response *RPCResponse) error {
	return s.WSConn.WriteJSON(response)
}

func (s *GameSession) sendError(id int, code int, message string) error {
	response := RPCResponse{
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	return s.sendResponse(&response)
}
