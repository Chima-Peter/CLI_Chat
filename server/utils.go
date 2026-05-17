package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

type requestPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Nick     string `json:"nick"`
	RoomID   string `json:"room_id"`
	Room     string `json:"room"`
	RoomName string `json:"room_name"`
	NewRoom  string `json:"new_room"`
	Password string `json:"password"`
	Message  string `json:"message"`
	MaxSize  *int   `json:"max_size"`
}

func decodePayload(raw json.RawMessage) requestPayload {
	var p requestPayload
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &p)
	}
	return p
}

func (p requestPayload) roomID() string {
	return strings.TrimSpace(p.RoomID)
}

func (p requestPayload) roomName() string {
	if name := strings.TrimSpace(p.Room); name != "" {
		return name
	}
	return strings.TrimSpace(p.RoomName)
}

func (p requestPayload) userID() string {
	return strings.TrimSpace(p.UserID)
}

func (p requestPayload) userName() string {
	if name := strings.TrimSpace(p.Username); name != "" {
		return name
	}
	return strings.TrimSpace(p.Nick)
}

func hasID(ids map[string]struct{}, id string) bool {
	_, ok := ids[id]
	return ok
}

func copyIDSet(src map[string]struct{}) map[string]struct{} {
	dst := make(map[string]struct{}, len(src))
	for id := range src {
		dst[id] = struct{}{}
	}
	return dst
}

func userNicks(users []map[string]string) []string {
	nicks := make([]string, 0, len(users))
	for _, user := range users {
		nicks = append(nicks, user["nick"])
	}
	return nicks
}

func (s *server) GetClientByID(user_id string) (*client, error) {
	user_id = strings.TrimSpace(user_id)
	if user_id == "" {
		return nil, fmt.Errorf("Provide user id")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	client_data, ok := s.clients[user_id]
	if !ok {
		return nil, fmt.Errorf("Client not found with id: %s", user_id)
	}

	return client_data, nil
}

func (s *server) GetClientByNick(nick string) (*client, error) {
	nick = strings.TrimSpace(nick)
	if nick == "" {
		return nil, fmt.Errorf("Provide a username")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client_data := range s.clients {
		if client_data.nick == nick {
			return client_data, nil
		}
	}
	return nil, fmt.Errorf("Client not found with nick: %s", nick)
}

func (s *server) resolveUser(userID, userName string) (*client, error) {
	if id := strings.TrimSpace(userID); id != "" {
		return s.GetClientByID(id)
	}
	return s.GetClientByNick(userName)
}

func (s *server) clientNick(user_id string) string {
	client_data, err := s.GetClientByID(user_id)
	if err != nil {
		return user_id
	}
	return client_data.nick
}

func (s *server) usersFromIDs(ids map[string]struct{}) []map[string]string {
	users := make([]map[string]string, 0, len(ids))
	for user_id := range ids {
		users = append(users, map[string]string{
			"user_id": user_id,
			"nick":    s.clientNick(user_id),
		})
	}
	return users
}

func (s *server) GetRoomByID(room_id string) (*room, error) {
	room_id = strings.TrimSpace(room_id)
	if room_id == "" {
		return nil, fmt.Errorf("Provide room id")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	room_data, ok := s.rooms[room_id]
	if !ok {
		return nil, fmt.Errorf("No room found with id: %s", room_id)
	}

	return room_data, nil
}

func (s *server) findRoomByName(room_name string) (*room, error) {
	room_name = strings.TrimSpace(room_name)
	if room_name == "" {
		return nil, fmt.Errorf("Provide room name")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, room_data := range s.rooms {
		if room_data.name == room_name {
			return room_data, nil
		}
	}

	return nil, fmt.Errorf("No room found with name: %s", room_name)
}

func (s *server) resolveRoom(roomID, roomName string) (*room, error) {
	if id := strings.TrimSpace(roomID); id != "" {
		return s.GetRoomByID(id)
	}
	return s.findRoomByName(roomName)
}

func (s *server) isRoomNameTaken(name string, excludeRoomID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for id, room_data := range s.rooms {
		if id != excludeRoomID && room_data.name == name {
			return true
		}
	}
	return false
}
