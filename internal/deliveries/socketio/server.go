package socketio

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type SocketIOServer struct {
	server *socketio.Server
}

func NewSocketIOServer() (*SocketIOServer, error) {
	server := socketio.NewServer(nil)

	// Обработка подключения клиента
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Printf("Клиент подключился: %s", s.ID())
		s.SetContext("")
		return nil
	})

	// Обработка отключения клиента
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("Клиент отключился: %s, причина: %s", s.ID(), reason)
	})

	// Обработка пользовательских событий
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Printf("Получено сообщение от клиента %s: %s", s.ID(), msg)
		s.Emit("reply", "Сообщение получено: "+msg)
	})

	return &SocketIOServer{server: server}, nil
}

func (s *SocketIOServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}

func (s *SocketIOServer) Close() {
	s.server.Close()
}
