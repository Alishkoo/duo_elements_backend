package websocket

import (
	"duo_elements/internal/data"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	clients       map[string]*websocket.Conn  // Список подключенных клиентов
	playersStates map[string]data.PlayerState // Состояние игроков (позиция, скорость и т.д.)
	mu            sync.Mutex                  // Мьютекс для безопасного доступа к clients
	upgrader      websocket.Upgrader          // WebSocket Upgrader
	shutdown      chan struct{}               // Канал для graceful shutdown
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:       make(map[string]*websocket.Conn),
		playersStates: make(map[string]data.PlayerState),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		shutdown: make(chan struct{}),
	}
}

// HandleConnections обрабатывает входящие WebSocket-соединения
func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка подключения WebSocket: %v", err)
		return
	}
	defer conn.Close()

	playerID := r.URL.Query().Get("playerID")
	if playerID == "" {
		log.Println("playerID не указан")
		return
	}

	// Добавляем клиента в список
	s.mu.Lock()
	s.clients[playerID] = conn
	s.mu.Unlock()
	log.Printf("Игрок подключился: %s", playerID)

	// Чтение сообщений от клиента
	for {
		select {
		case <-s.shutdown:
			log.Printf("Закрытие соединения для игрока: %s", playerID)
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Игрок %s отключился: %v", playerID, err)
				// Удаляем игрока из списка
				s.mu.Lock()
				delete(s.clients, playerID)
				delete(s.playersStates, playerID)
				s.mu.Unlock()
				log.Printf("Игрок %s удален из списка клиентов", playerID)
				return
			}

			// Обрабатываем сообщение
			s.HandleMessage(playerID, message)
		}
	}
}

// HandleMessage обрабатывает сообщения от игроков
func (s *WebSocketServer) HandleMessage(playerID string, message []byte) {
	var baseMessage struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(message, &baseMessage); err != nil {
		log.Printf("Ошибка парсинга сообщения от игрока %s: %v", playerID, err)
		return
	}

	switch baseMessage.Type {
	case "spawn":
		var spawnMsg data.SpawnMessage
		if err := json.Unmarshal(message, &spawnMsg); err != nil {
			log.Printf("Ошибка парсинга SpawnMessage от игрока %s: %v", playerID, err)
			return
		}

		// Сохраняем состояние нового игрока
		s.mu.Lock()
		s.playersStates[playerID] = data.PlayerState{
			Position: spawnMsg.Position,
			Velocity: data.Velocity{DX: 0, DY: 0}, // Начальная скорость
		}

		// Отправляем новому игроку список уже подключенных игроков
		existingPlayers := make([]data.SpawnMessage, 0, len(s.playersStates))
		for id, state := range s.playersStates {
			if id != playerID {
				existingPlayers = append(existingPlayers, data.SpawnMessage{
					Type:      "spawn",
					PlayerID:  id,
					Character: spawnMsg.Character,
					Position:  state.Position,
				})
			}
		}

		// Создаём сообщение для нового игрока
		response := data.ExistingPlayersMessage{
			Type:    "existing_players",
			Players: existingPlayers,
		}

		// Отправляем сообщение новому игроку
		conn := s.clients[playerID]
		if conn != nil {
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Ошибка отправки списка игроков новому игроку %s: %v", playerID, err)
			}
		}

		// Уведомляем остальных игроков о новом участнике
		for id, conn := range s.clients {
			if id != playerID {
				err := conn.WriteJSON(spawnMsg)
				if err != nil {
					log.Printf("Ошибка отправки сообщения игроку %s: %v", id, err)
					conn.Close()
					delete(s.clients, id)
				}
			}
		}
		s.mu.Unlock()

		log.Printf("Игрок %s появился: %+v", playerID, spawnMsg)

	case "sync":
		var syncMsg data.SyncMessage
		if err := json.Unmarshal(message, &syncMsg); err != nil {
			log.Printf("Ошибка парсинга SyncMessage от игрока %s: %v", playerID, err)
			return
		}

		// Сохраняем состояние игрока
		s.mu.Lock()
		s.playersStates[playerID] = data.PlayerState{
			Position: syncMsg.Position,
			Velocity: syncMsg.Velocity,
		}
		s.mu.Unlock()

		log.Printf("Получена синхронизация от игрока %s: %+v", playerID, syncMsg)

		// Пересылаем другим игрокам
		s.BroadcastToOthers(playerID, message)

	case "input":
		var inputMsg data.InputMessage
		if err := json.Unmarshal(message, &inputMsg); err != nil {
			log.Printf("Ошибка парсинга InputMessage от игрока %s: %v", playerID, err)
			return
		}
		log.Printf("Получено действие от игрока %s: %+v", playerID, inputMsg)

		// Пересылаем другим игрокам
		s.BroadcastToOthers(playerID, message)
	default:
		log.Printf("Неизвестный тип сообщения от игрока %s: %s", playerID, baseMessage.Type)
	}
}

// BroadcastToOthers отправляет сообщение всем игрокам, кроме указанного
func (s *WebSocketServer) BroadcastToOthers(playerID string, message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, conn := range s.clients {
		if id == playerID {
			continue // Не отправляем сообщение самому игроку
		}

		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Ошибка отправки сообщения игроку %s: %v", id, err)
			conn.Close()
			delete(s.clients, id)
		}
	}
}

// Stop завершает работу WebSocket-сервера
func (s *WebSocketServer) Stop() {
	log.Println("Завершение работы WebSocket-сервера...")
	close(s.shutdown)

	s.mu.Lock()
	defer s.mu.Unlock()

	for playerID, conn := range s.clients {
		log.Printf("Закрытие соединения для игрока: %s", playerID)
		conn.Close()
	}
	log.Println("WebSocket-сервер завершил работу")
}
