package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/chima/CLI_Chat/protocol"
	"github.com/chzyer/readline"
)

const defaultPrompt = "> "

type session struct {
	conn     net.Conn
	rl       *readline.Instance
	incoming chan protocol.Message

	displayMu sync.Mutex
	promptMu  sync.Mutex
	prompt    *protocol.Message
}

func runSession(conn net.Conn, rl *readline.Instance) {
	s := &session{
		conn:     conn,
		rl:       rl,
		incoming: make(chan protocol.Message, 16),
	}

	done := make(chan struct{})

	go s.readFromServer(done)
	go s.processServerMessages()
	go func() {
		<-done
		s.rl.Close()
	}()
	s.inputLoop()
}

func (s *session) readFromServer(done chan struct{}) {
	defer close(s.incoming)

	decoder := json.NewDecoder(s.conn)
	for {
		var msg protocol.Message
		if err := decoder.Decode(&msg); err != nil {
			if err != io.EOF {
				s.writeDisplay("Connection closed")
				close(done)
				return
			}
			return
		}
		s.incoming <- msg
	}
}

func (s *session) processServerMessages() {
	for msg := range s.incoming {
		s.handleServerMessage(msg)
	}
}

func (s *session) handleServerMessage(msg protocol.Message) {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()

	switch {
	case needsUserReply(msg.Action):
		s.clearCurrentInputLocked()
		s.setPromptLocked(&msg)
	case isTerminalAction(msg.Action):
		if msg.ResponseMsg != "" {
			s.writeDisplayLocked(msg.ResponseMsg)
		}
	case isDisplayOnlyAction(msg.Action):
		if msg.ResponseMsg != "" {
			s.writeDisplayLocked(msg.ResponseMsg)
		}
	default:
		if msg.ResponseMsg == "" {
			s.writeDisplayLocked(fmt.Sprintf("[server action %d]", msg.Action))
		}
	}

	s.rl.Refresh()
}

func (s *session) writeDisplay(msg string) {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()
	s.writeDisplayLocked(msg)
	s.rl.Refresh()
}

func (s *session) writeDisplayLocked(msg string) {
	if msg == "" {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	s.rl.SetPrompt("")
	_, _ = s.rl.Write([]byte("> " + msg))
}

func (s *session) clearSubmittedInput() {
	s.displayMu.Lock()
	defer s.displayMu.Unlock()
	fmt.Fprint(s.rl.Stdout(), "\033[A\033[2K")
	s.rl.SetPrompt(defaultPrompt)
	s.rl.Refresh()
}

func (s *session) setPromptLocked(msg *protocol.Message) {
	s.promptMu.Lock()
	s.prompt = msg
	s.promptMu.Unlock()

	prompt := defaultPrompt
	if msg != nil && msg.ResponseMsg != "" {
		prompt += msg.ResponseMsg
	}
	s.rl.SetPrompt(prompt)
}

func (s *session) clearPromptLocked() {
	s.promptMu.Lock()
	s.prompt = nil
	s.promptMu.Unlock()
	s.rl.SetPrompt(defaultPrompt)
}

func (s *session) pendingPrompt() *protocol.Message {
	s.promptMu.Lock()
	defer s.promptMu.Unlock()
	return s.prompt
}

func (s *session) clearCurrentInputLocked() {
	_, _ = s.rl.WriteStdin([]byte{
		byte(readline.CharLineEnd),
		byte(readline.CharCtrlU),
		byte(readline.CharKill),
	})
	s.rl.Refresh()
}

func (s *session) inputLoop() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			msg, _ := buildCommandMessage("/quit")
			if err := writeMessage(s.conn, msg); err != nil {
				return
			}
			return
		}

		if prompt := s.pendingPrompt(); prompt != nil {
			if err := s.sendPromptReply(prompt, line); err != nil {
				s.writeDisplay(err.Error())
			}
			continue
		}

		msg, err := buildMessage(line)
		if fsc, ok := errors.AsType[*fileSendCommand](err); ok {
			if sendErr := s.sendFile(fsc.filename); sendErr != nil {
				s.writeDisplay(sendErr.Error())
			} else {
				s.writeDisplay("File sent.")
			}
			s.clearSubmittedInput()
			continue
		}
		if errors.Is(err, errClientOnly) {
			continue
		}
		if msg == nil {
			continue
		}
		if err != nil {
			s.writeDisplay(err.Error())
			continue
		}

		if msg.Action == protocol.LOGOUT {
			if err := writeMessage(s.conn, msg); err != nil {
				s.writeDisplay(err.Error())
				return
			}
			return
		}

		if err := writeMessage(s.conn, msg); err != nil {
			s.writeDisplay(err.Error())
			return
		}
		s.clearSubmittedInput()
	}
}

func (s *session) sendPromptReply(prompt *protocol.Message, input string) error {
	msg, err := buildPromptReply(prompt.Action, prompt.Payload, input)
	if err != nil {
		return err
	}

	s.displayMu.Lock()
	s.clearPromptLocked()
	s.rl.Refresh()
	s.displayMu.Unlock()

	return writeMessage(s.conn, msg)
}

func writeMessage(conn net.Conn, msg *protocol.Message) error {
	return json.NewEncoder(conn).Encode(msg)
}

func (s *session) sendFile(filename string) error {
	zipPath, err := ZipFile(filename)
	if err != nil {
		return err
	}
	return s.streamZipToServer(zipPath)
}

func (s *session) streamZipToServer(zipPath string) error {
	info, err := os.Stat(zipPath)
	if err != nil {
		return fmt.Errorf("stat zip %q: %w", zipPath, err)
	}

	msg := &protocol.Message{
		Action: protocol.SEND_FILE,
		Payload: marshalStringPayload(map[string]string{
			"name": info.Name(),
			"size": strconv.FormatInt(info.Size(), 10),
		}),
	}
	if err := writeMessage(s.conn, msg); err != nil {
		return fmt.Errorf("send file metadata: %w", err)
	}

	f, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("open zip %q: %w", zipPath, err)
	}
	defer f.Close()

	n, err := io.CopyN(s.conn, f, info.Size())
	if err != nil {
		return fmt.Errorf("stream zip %q: %w", zipPath, err)
	}
	if n != info.Size() {
		return fmt.Errorf("stream zip %q: sent %d of %d bytes", zipPath, n, info.Size())
	}
	return nil
}
