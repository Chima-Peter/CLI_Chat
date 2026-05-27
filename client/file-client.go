package client

import (
	"crypto/tls"
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
		return fmt.Errorf("connect to file server: %w", err)
	}
	defer conn.Close()

	zipPath, err := streamZipToServer(conn, filepath, nick)
	if err != nil {
		return fmt.Errorf("send file: %w", err)
	}
	fmt.Printf("File sent to %s (%s)\n", friend, zipPath)
	return nil
}

func streamZipToServer(conn net.Conn, zipPath string, nick string) (string, error) {
	info, err := os.Stat(zipPath)
	if err != nil {
		return "", fmt.Errorf("stat zip %q: %w", zipPath, err)
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
		return "", fmt.Errorf("send file metadata: %w", err)
	}

	f, err := os.Open(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip %q: %w", zipPath, err)
	}
	defer f.Close()

	n, err := io.CopyN(conn, f, info.Size())
	if err != nil {
		return "", fmt.Errorf("stream zip %q: %w", zipPath, err)
	}
	if n != info.Size() {
		return "", fmt.Errorf("stream zip %q: sent %d of %d bytes", zipPath, n, info.Size())
	}
	return zipPath, nil
}
