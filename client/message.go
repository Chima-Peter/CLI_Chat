package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chima/CLI_Chat/protocol"
)

func isTerminalAction(action protocol.ActionType) bool {
	return action == protocol.DONE || action == protocol.ERR
}

func isDisplayOnlyAction(action protocol.ActionType) bool {
	switch action {
	case protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
		return true
	default:
		return false
	}
}

func needsUserReply(action protocol.ActionType) bool {
	return !isTerminalAction(action) && !isDisplayOnlyAction(action)
}

func buildMessage(line string) (*protocol.Message, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty input")
	}

	if strings.HasPrefix(line, "/") {
		return buildCommandMessage(line)
	}

	payload, _ := json.Marshal(map[string]string{"message": line})
	return &protocol.Message{
		Action:  protocol.SEND_MSG,
		Payload: payload,
	}, nil
}

func buildCommandMessage(line string) (*protocol.Message, error) {
	parts := strings.SplitN(line, " ", 2)
	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}

	marshal := func(fields map[string]string) json.RawMessage {
		payload, _ := json.Marshal(fields)
		return payload
	}

	switch cmd {
	case "/nick":
		if arg == "" {
			return nil, fmt.Errorf("usage: /nick <name>")
		}
		return &protocol.Message{
			Action:  protocol.LOGIN,
			Payload: marshal(map[string]string{"username": arg}),
		}, nil
	case "/join":
		if arg == "" {
			return nil, fmt.Errorf("usage: /join <room>")
		}
		return &protocol.Message{
			Action:  protocol.JOIN_ROOM,
			Payload: marshal(map[string]string{"room": arg}),
		}, nil
	case "/rooms":
		return &protocol.Message{Action: protocol.LIST_ROOMS}, nil
	case "/myrooms":
		return &protocol.Message{Action: protocol.LIST_MY_ROOMS}, nil
	case "/msg":
		if arg == "" {
			return nil, fmt.Errorf("usage: /msg <message>")
		}
		return &protocol.Message{
			Action:  protocol.SEND_MSG,
			Payload: marshal(map[string]string{"message": arg}),
		}, nil
	case "/quit":
		return &protocol.Message{Action: protocol.LOGOUT}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

func buildPromptReply(action protocol.ActionType, serverPayload json.RawMessage, input string) (*protocol.Message, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	fields := map[string]string{}
	if len(serverPayload) > 0 {
		_ = json.Unmarshal(serverPayload, &fields)
	}

	switch action {
	case protocol.LOGIN:
		fields["username"] = input
	case protocol.SET_ROOM_PASSWORD, protocol.GET_ROOM_PASSWORD:
		fields["password"] = input
	case protocol.SEND_MSG, protocol.MESSAGE_FRIEND:
		fields["message"] = input
	default:
		if _, ok := fields["message"]; !ok {
			fields["message"] = input
		}
	}

	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return &protocol.Message{
		Action:  action,
		Payload: payload,
	}, nil
}

func printTerminal(format string, args ...any) {
	fmt.Printf(format, args...)
}
