package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

type client struct {
	id                      string
	conn                    net.Conn
	nick                    string
	room                    *room
	my_rooms                map[string]*room
	friends                 map[string]struct{}
	pending_friend_requests map[string]struct{}
	sent_friend_request     map[string]struct{}
	blocked_users           map[string]struct{}
	room_invites            map[string]*room
	online                  bool
	mu                      sync.RWMutex
}

func (cl *client) err(err error) {
	cl.send_user_message(map[string]any{}, ERR, err.Error())
}

func (cl *client) send_user_message(payload_data map[string]any, next_action ActionType, response_msg string) {
	response := cl.prepare_response(payload_data, next_action, response_msg)
	encoder := json.NewEncoder(cl.conn)
	encoder.Encode(response)
}

func (cl *client) prepare_response(payload_data map[string]any, next_action ActionType, response_msg string) *Message {
	payload, _ := json.Marshal(payload_data)

	return &Message{
		Action:      next_action,
		ResponseMsg: response_msg,
		Payload:     payload,
	}
}

func (cl *client) isBlockedWith(other *client) (selfBlocked, blockedByOther bool) {
	cl.mu.RLock()
	_, selfBlocked = cl.blocked_users[other.id]
	cl.mu.RUnlock()

	other.mu.RLock()
	_, blockedByOther = other.blocked_users[cl.id]
	other.mu.RUnlock()

	return selfBlocked, blockedByOther
}

func (cl *client) SendFriendRequest(friend *client) {
	if cl.id == friend.id {
		cl.err(fmt.Errorf("You cannot send a friend request to yourself."))
		return
	}

	selfBlocked, blockedByOther := cl.isBlockedWith(friend)
	if selfBlocked {
		cl.err(fmt.Errorf("You have blocked %s. Unblock them to send a friend request.", friend.nick))
		return
	}
	if blockedByOther {
		cl.err(fmt.Errorf("You cannot send a friend request to %s.", friend.nick))
		return
	}

	cl.mu.RLock()
	if hasID(cl.friends, friend.id) {
		cl.mu.RUnlock()
		cl.err(fmt.Errorf("You are already friends with %s.", friend.nick))
		return
	}
	if hasID(cl.sent_friend_request, friend.id) {
		cl.mu.RUnlock()
		cl.err(fmt.Errorf("You have already sent a friend request to %s.", friend.nick))
		return
	}
	if hasID(cl.pending_friend_requests, friend.id) {
		cl.mu.RUnlock()
		cl.err(fmt.Errorf("%s has already sent you a friend request.", friend.nick))
		return
	}
	cl.mu.RUnlock()

	cl.mu.Lock()
	cl.sent_friend_request[friend.id] = struct{}{}
	cl.mu.Unlock()

	friend.mu.Lock()
	friend.pending_friend_requests[cl.id] = struct{}{}
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Friend request sent to %s.", friend.nick))
	friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Received a friend request from %s.", cl.nick))
}

func (cl *client) AcceptFriendRequest(friend *client) {
	cl.mu.RLock()
	exists := hasID(cl.pending_friend_requests, friend.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("You have not received a friend request from this user."))
		return
	}

	cl.mu.Lock()
	cl.friends[friend.id] = struct{}{}
	delete(cl.pending_friend_requests, friend.id)
	cl.mu.Unlock()

	friend.mu.Lock()
	friend.friends[cl.id] = struct{}{}
	delete(friend.sent_friend_request, cl.id)
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("You are now friends with %s.", friend.nick))
	friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s accepted your friend request.", cl.nick))
}

func (cl *client) DeleteFriendRequest(friend *client) {
	cl.mu.RLock()
	exists := hasID(cl.sent_friend_request, friend.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("You have not sent a friend request to this user."))
		return
	}

	friend.mu.RLock()
	ok := hasID(friend.pending_friend_requests, cl.id)
	friend.mu.RUnlock()
	if !ok {
		cl.err(fmt.Errorf("Error deleting pending requests. Record not found."))
		return
	}

	cl.mu.Lock()
	delete(cl.sent_friend_request, friend.id)
	cl.mu.Unlock()

	friend.mu.Lock()
	delete(friend.pending_friend_requests, cl.id)
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Friend request to %s cancelled.", friend.nick))
	friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s cancelled their friend request.", cl.nick))
}

func (cl *client) RejectFriendRequest(friend *client) {
	cl.mu.RLock()
	exists := hasID(cl.pending_friend_requests, friend.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("You have not received a friend request from this user."))
		return
	}

	friend.mu.RLock()
	ok := hasID(friend.sent_friend_request, cl.id)
	friend.mu.RUnlock()
	if !ok {
		cl.err(fmt.Errorf("Error rejecting sent requests. Record not found."))
		return
	}

	cl.mu.Lock()
	delete(cl.pending_friend_requests, friend.id)
	cl.mu.Unlock()

	friend.mu.Lock()
	delete(friend.sent_friend_request, cl.id)
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Friend request from %s rejected.", friend.nick))
	friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s rejected your friend request.", cl.nick))
}

func (cl *client) FetchSentRequests(users []map[string]string) {
	if len(users) == 0 {
		cl.send_user_message(map[string]any{"requests": users}, DONE, "You have no sent requests.")
		return
	}

	nicks := userNicks(users)
	cl.send_user_message(
		map[string]any{"requests": users},
		DONE,
		fmt.Sprintf("You have sent request to: %s", strings.Join(nicks, ", ")),
	)
}

func (cl *client) FetchPendingRequests(users []map[string]string) {
	if len(users) == 0 {
		cl.send_user_message(map[string]any{"requests": users}, DONE, "You have no pending requests.")
		return
	}

	nicks := userNicks(users)
	cl.send_user_message(
		map[string]any{"requests": users},
		DONE,
		fmt.Sprintf("You have pending requests from: %s", strings.Join(nicks, ", ")),
	)
}

func (cl *client) GetFriends(users []map[string]string) {
	if len(users) == 0 {
		cl.send_user_message(map[string]any{"friends": users}, DONE, "You have no friends.")
		return
	}

	nicks := userNicks(users)
	cl.send_user_message(
		map[string]any{"friends": users},
		DONE,
		fmt.Sprintf("Your friends: %s", strings.Join(nicks, ", ")),
	)
}

func (cl *client) DeleteFriend(friend *client) {
	cl.mu.RLock()
	exists := hasID(cl.friends, friend.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("You are not friends with %s.", friend.nick))
		return
	}

	cl.mu.Lock()
	delete(cl.friends, friend.id)
	cl.mu.Unlock()

	friend.mu.Lock()
	delete(friend.friends, cl.id)
	friend.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Removed %s from your friends.", friend.nick))

	friend.mu.RLock()
	online := friend.online
	friend.mu.RUnlock()
	if online {
		friend.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s removed you from their friends.", cl.nick))
	}
}

func (cl *client) MessageFriend(friend *client, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		cl.err(fmt.Errorf("Message cannot be empty."))
		return
	}

	cl.mu.RLock()
	isFriend := hasID(cl.friends, friend.id)
	cl.mu.RUnlock()
	if !isFriend {
		cl.err(fmt.Errorf("You are not friends with %s.", friend.nick))
		return
	}

	selfBlocked, blockedByOther := cl.isBlockedWith(friend)
	if selfBlocked {
		cl.err(fmt.Errorf("You have blocked %s.", friend.nick))
		return
	}
	if blockedByOther {
		cl.err(fmt.Errorf("You cannot message %s.", friend.nick))
		return
	}

	friend.mu.RLock()
	online := friend.online
	friend.mu.RUnlock()
	if !online {
		cl.err(fmt.Errorf("%s is offline.", friend.nick))
		return
	}

	friend.send_user_message(
		map[string]any{
			"from_user_id": cl.id,
			"from":         cl.nick,
			"message":      message,
		},
		MESSAGE_FRIEND,
		fmt.Sprintf("Message from %s: %s", cl.nick, message),
	)
	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("Message sent to %s.", friend.nick))
}

func (cl *client) BlockUser(target *client) {
	if cl.id == target.id {
		cl.err(fmt.Errorf("You cannot block yourself."))
		return
	}

	cl.mu.RLock()
	if hasID(cl.blocked_users, target.id) {
		cl.mu.RUnlock()
		cl.err(fmt.Errorf("%s is already blocked.", target.nick))
		return
	}
	cl.mu.RUnlock()

	cl.mu.Lock()
	cl.blocked_users[target.id] = struct{}{}
	cl.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s has been blocked.", target.nick))
}

func (cl *client) UnblockUser(target *client) {
	cl.mu.RLock()
	exists := hasID(cl.blocked_users, target.id)
	cl.mu.RUnlock()
	if !exists {
		cl.err(fmt.Errorf("%s is not blocked.", target.nick))
		return
	}

	cl.mu.Lock()
	delete(cl.blocked_users, target.id)
	cl.mu.Unlock()

	cl.send_user_message(map[string]any{}, DONE, fmt.Sprintf("%s has been unblocked.", target.nick))
}

func (cl *client) setOnline(online bool) {
	cl.mu.Lock()
	cl.online = online
	cl.mu.Unlock()
}

func (cl *client) GetUserStatus(target *client) {
	cl.mu.RLock()
	isBlocked := hasID(cl.blocked_users, target.id)
	cl.mu.RUnlock()

	target.mu.RLock()
	online := target.online
	roomID := ""
	roomName := ""
	if target.room != nil {
		roomID = target.room.id
		roomName = target.room.name
	}
	target.mu.RUnlock()

	status := "offline"
	if online {
		status = "online"
	}

	cl.send_user_message(map[string]any{
		"user_id": target.id,
		"nick":    target.nick,
		"online":  online,
		"blocked": isBlocked,
		"room_id": roomID,
		"room":    roomName,
	}, DONE, fmt.Sprintf("%s is %s.", target.nick, status))
}
