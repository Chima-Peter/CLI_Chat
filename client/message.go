package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chima/CLI_Chat/protocol"
)

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
