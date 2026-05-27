package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/chima/CLI_Chat/protocol"
)

const fileServerIdleTimeout = 20 * time.Minute

func StartFileServer(notify func(host string, port int) error) error {
	listener, err := tls.Listen("tcp", ":0", &tls.Config{
		Certificates: []tls.Certificate{},
	})
	if err != nil {
		return err
	}
	defer listener.Close()

	host, port := listenerEndpoint(listener)
	if notify != nil {
		if err := notify(host, port); err != nil {
			return err
		}
	}

	log.Printf("file server listening on %s:%d", host, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var timerMu sync.Mutex
	var idleTimer *time.Timer

	extendTimeout := func() {
		timerMu.Lock()
		defer timerMu.Unlock()
		if idleTimer != nil {
			idleTimer.Stop()
		}
		idleTimer = time.AfterFunc(fileServerIdleTimeout, cancel)
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	extendTimeout()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				log.Println("file server stopped:", ctx.Err())
				return ctx.Err()
			}
			log.Println("Unable to accept connection:", err.Error())
			continue
		}
		extendTimeout()
		go handleFile(conn)
	}
}

func listenerEndpoint(l net.Listener) (host string, port int) {
	tcpAddr := l.Addr().(*net.TCPAddr)
	if tcpAddr.IP.IsUnspecified() {
		host = "127.0.0.1"
	} else {
		host = tcpAddr.IP.String()
	}
	return host, tcpAddr.Port
}

func handleFile(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var request protocol.Message
		if err := decoder.Decode(&request); err != nil {
			return
		}

		if request.Action == protocol.UPLOAD_FILE {
			var meta struct {
				Name string `json:"name"`
				Size string `json:"size"`
			}
			_ = json.Unmarshal(request.Payload, &meta)
			size, parseErr := strconv.ParseInt(meta.Size, 10, 64)

			if parseErr != nil || size < 0 {
				log.Printf("file upload: invalid size metadata: %v", parseErr)
				continue
			}

			body := io.LimitReader(io.MultiReader(decoder.Buffered(), conn), size)
			received, readErr := io.ReadAll(body)
			if readErr != nil {
				log.Printf("file upload: read body: %v", readErr)
				continue
			}
			log.Printf("file upload: received %d bytes", len(received))

			home, err := os.UserHomeDir()
			if err != nil {
				log.Printf("file upload: resolve home directory: %v", err)
				continue
			}
			destPath := filepath.Join(home, "Downloads", filepath.Base(meta.Name))

			out, err := os.Create(destPath)
			if err != nil {
				log.Printf("file upload: create file: %v", err)
				continue
			}
			_, err = out.Write(received)
			closeErr := out.Close()
			if err != nil {
				log.Printf("file upload: write file: %v", err)
				continue
			}
			if closeErr != nil {
				log.Printf("file upload: close file: %v", closeErr)
				continue
			}

			log.Printf("file upload: saved to %s", destPath)
			return
		}
	}
}
