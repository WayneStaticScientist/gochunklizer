package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/doquangtan/socketio/v4"
	"opechains.shop/chunklizer/v2/types"
)

type WebSocketManager struct {
	Io *socketio.Io
}

func NewWebSocketManager(io *socketio.Io) *WebSocketManager {
	return &WebSocketManager{
		Io: io,
	}
}
func (wsManager *WebSocketManager) AuthenticateUser(payload *socketio.EventPayload) (string, error) {
	if len(payload.Data) <= 0 {
		return "", nil
	}
	if load, ok := payload.Data[0].(string); ok {
		var myMap map[string]string
		err := json.Unmarshal([]byte(load), &myMap)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return "", err
		}
		return myMap["Authorization"], nil
	}
	return "", nil
}
func (wsManager *WebSocketManager) NewUserConnection(s *socketio.Socket) {
	s.On("authenticate", func(payload *socketio.EventPayload) {
		user, _ := wsManager.AuthenticateUser(payload)
		if user == "" {
			s.Disconnect()
			return
		}
		s.Join(user)
		log.Println("user connected:", user)
		s.To(user).Emit("chunk:start", types.SocketChunckMessage{
			IsSuccess: false,
			Message:   "hello world",
			HasError:  false,
			Data:      ""})
	})

}
