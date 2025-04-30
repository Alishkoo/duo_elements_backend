// package usecases

// import (
// 	"duo_elements/internal/data"
// 	"encoding/json"
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// type WebSocketUsecase interface {
// 	HandleMessage(playerID string, message []byte, clients map[string]*websocket.Conn)
// }

// type websocketUsecase struct{}

// func NewWebSocketUsecase() WebSocketUsecase {
// 	return &websocketUsecase{}
// }

// func (u *websocketUsecase) HandleMessage(playerID string, message []byte, clients map[string]*websocket.Conn) {
// 	log.Printf("Обработка сообщения от игрока %s: %s", playerID, string(message))
// 	var baseMessage struct {
// 		Type string `json:"type"`
// 	}

// 	if err := json.Unmarshal(message, &baseMessage); err != nil {
// 		log.Printf("Ошибка парсинга сообщения от игрока %s: %v", playerID, err)
// 		return
// 	}

// 	switch baseMessage.Type {
// 	case "input":
// 		var inputMsg data.InputMessage
// 		if err := json.Unmarshal(message, &inputMsg); err != nil {
// 			log.Printf("Ошибка парсинга InputMessage от игрока %s: %v", playerID, err)
// 			return
// 		}
// 		log.Printf("Получено действие от игрока %s: %+v", playerID, inputMsg)

// 		// Пересылаем другим игрокам
// 		u.BroadcastToOthers(playerID, message, clients)

// 	case "sync":
// 		var syncMsg data.SyncMessage
// 		if err := json.Unmarshal(message, &syncMsg); err != nil {
// 			log.Printf("Ошибка парсинга SyncMessage от игрока %s: %v", playerID, err)
// 			return
// 		}

// 		// Сохраняем состояние игрока
// 		log.Printf("Получена синхронизация от игрока %s: %+v", playerID, syncMsg)

// 		// Пересылаем другим игрокам
// 		u.BroadcastToOthers(playerID, message, clients)

// 	case "spawn":
// 		var spawnMsg data.SpawnMessage
// 		if err := json.Unmarshal(message, &spawnMsg); err != nil {
// 			log.Printf("Ошибка парсинга SpawnMessage от игрока %s: %v", playerID, err)
// 			return
// 		}
// 		log.Printf("Игрок %s появился: %+v", playerID, spawnMsg)

// 		// Пересылаем другим игрокам
// 		u.BroadcastToOthers(playerID, message, clients)

// 	case "disconnect":
// 		var disconnectMsg data.DisconnectMessage
// 		if err := json.Unmarshal(message, &disconnectMsg); err != nil {
// 			log.Printf("Ошибка парсинга DisconnectMessage от игрока %s: %v", playerID, err)
// 			return
// 		}
// 		log.Printf("Игрок %s отключился: %+v", playerID, disconnectMsg)

// 	default:
// 		log.Printf("Неизвестный тип сообщения от игрока %s: %s", playerID, baseMessage.Type)
// 	}
// }

// // BroadcastToOthers отправляет сообщение всем игрокам, кроме указанного
// // playerID. Используется для синхронизации состояния игры между игроками.
// func (u *websocketUsecase) BroadcastToOthers(playerID string, message []byte, clients map[string]*websocket.Conn) {
// 	for id, conn := range clients {
// 		if id == playerID {
// 			continue // Не отправляем сообщение самому игроку
// 		}

// 		err := conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			log.Printf("Ошибка отправки сообщения игроку %s: %v", id, err)
// 			conn.Close()
// 			delete(clients, id)
// 		}
// 	}
// }
