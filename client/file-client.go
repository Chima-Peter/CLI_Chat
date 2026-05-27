package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/chima/CLI_Chat/protocol"
)

func CreateFileClient(host string, port int, nick string, friend string, filepath string) error {
	conn, err := tls.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("could not connect to %s to send the file: %w", friend, err)
	}
	defer conn.Close()

	if err := streamZipToServer(conn, filepath, nick, friend); err != nil {
		return err
	}
	return nil
}

func streamZipToServer(conn net.Conn, zipPath string, nick string, friend string) error {
	info, err := os.Stat(zipPath)
	if err != nil {
		return fmt.Errorf("could not read the prepared file: %w", err)
	}

	msg := &protocol.Message{
		Action: protocol.UPLOAD_FILE,
		Payload: marshalStringPayload(map[string]string{
			"name": info.Name(),
			"size": strconv.FormatInt(info.Size(), 10),
			"nick": nick,
		}),
	}
	if err := writeMessage(conn, msg); err != nil {
		return fmt.Errorf("could not start sending the file to %s: %w", friend, err)
	}

	f, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("could not open the prepared file: %w", err)
	}
	defer f.Close()

	n, err := io.CopyN(conn, f, info.Size())
	if err != nil {
		return fmt.Errorf("connection lost while sending the file to %s: %w", friend, err)
	}
	if n != info.Size() {
		return fmt.Errorf("only sent %d of %d bytes to %s", n, info.Size(), friend)
	}

	return readFileTransferReply(conn, friend)
}

func readFileTransferReply(conn net.Conn, friend string) error {
	var reply protocol.Message
	if err := json.NewDecoder(conn).Decode(&reply); err != nil {
		return fmt.Errorf("no confirmation received from %s", friend)
	}
	if reply.Action == protocol.ERR {
		msg := reply.ResponseMsg
		if msg == "" {
			msg = fmt.Sprintf("%s could not accept the file", friend)
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}
