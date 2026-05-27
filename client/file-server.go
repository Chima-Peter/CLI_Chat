package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chima/CLI_Chat/protocol"
)

const fileServerIdleTimeout = 20 * time.Minute

func devTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("could not create TLS certificate: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"CLI Chat"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("could not create TLS certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("could not load TLS certificate: %w", err)
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func StartFileServer(notify func(host string, port int) error) error {
	tlsConfig, err := devTLSConfig()
	if err != nil {
		return err
	}

	listener, err := tls.Listen("tcp", ":0", tlsConfig)
	if err != nil {
		return fmt.Errorf("could not listen for incoming files: %w", err)
	}
	defer listener.Close()

	host, port := listenerEndpoint(listener)
	if notify != nil {
		if err := notify(host, port); err != nil {
			return fmt.Errorf("could not register file receiver with the chat server: %w", err)
		}
	}

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
				return nil
			}
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

func writeFileConnError(conn net.Conn, userMsg string) {
	_ = writeMessage(conn, &protocol.Message{
		Action:      protocol.ERR,
		ResponseMsg: userMsg,
		Payload:     marshalStringPayload(map[string]string{"message": userMsg}),
	})
}

func writeFileConnOK(conn net.Conn, userMsg string) {
	_ = writeMessage(conn, &protocol.Message{
		Action:      protocol.DONE,
		ResponseMsg: userMsg,
	})
}

func handleFile(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var request protocol.Message
		if err := decoder.Decode(&request); err != nil {
			return
		}

		if request.Action != protocol.UPLOAD_FILE {
			writeFileConnError(conn, "Unsupported file request.")
			continue
		}

		var meta struct {
			Name string `json:"name"`
			Size string `json:"size"`
			Nick string `json:"nick"`
		}
		if err := json.Unmarshal(request.Payload, &meta); err != nil {
			writeFileConnError(conn, "Invalid file metadata.")
			return
		}

		name := strings.TrimSpace(meta.Name)
		if name == "" {
			writeFileConnError(conn, "File name is missing.")
			return
		}

		size, parseErr := strconv.ParseInt(meta.Size, 10, 64)
		if parseErr != nil || size < 0 {
			writeFileConnError(conn, "Invalid file size.")
			return
		}

		body := io.LimitReader(io.MultiReader(decoder.Buffered(), conn), size)
		received, readErr := io.ReadAll(body)
		if readErr != nil {
			writeFileConnError(conn, "Could not receive the file.")
			return
		}
		if int64(len(received)) != size {
			writeFileConnError(conn, fmt.Sprintf("Incomplete file received (%d of %d bytes).", len(received), size))
			return
		}

		home, err := os.UserHomeDir()
		if err != nil {
			writeFileConnError(conn, "Could not save the file to your Downloads folder.")
			return
		}

		safeName := filepath.Base(name)
		destPath := filepath.Join(home, "Downloads", safeName)

		out, err := os.Create(destPath)
		if err != nil {
			writeFileConnError(conn, fmt.Sprintf("Could not save %q to Downloads.", safeName))
			return
		}
		if _, err = out.Write(received); err != nil {
			_ = out.Close()
			writeFileConnError(conn, fmt.Sprintf("Could not write %q to Downloads.", safeName))
			return
		}
		if err := out.Close(); err != nil {
			writeFileConnError(conn, fmt.Sprintf("Could not finish saving %q to Downloads.", safeName))
			return
		}

		writeFileConnOK(conn, fmt.Sprintf("Saved to Downloads as %s.", safeName))
		fmt.Printf("Document received from %s saved to Downloads as %s\n", meta.Nick, safeName)
		return
	}
}
