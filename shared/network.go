// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package shared

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// RFC3493Dialer dialer
func RFC3493Dialer(network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}
	for _, a := range addrs {
		c, err := net.DialTimeout(network, net.JoinHostPort(a, port), 10*time.Second)
		if err != nil {
			continue
		}
		return c, err
	}
	return nil, fmt.Errorf("Unable to connect to: " + address)
}

// InitTLSConfig returns a tls.Config populated with default encryption
// parameters. This is used as baseline config server certificate
func InitTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		},
		PreferServerCipherSuites: true,
	}
}

func finalizeTLSConfig(tlsConfig *tls.Config, tlsRemoteCert *x509.Certificate) {
	// Trusted certificates
	if tlsRemoteCert != nil {
		caCertPool := tlsConfig.RootCAs
		if caCertPool == nil {
			caCertPool = x509.NewCertPool()
		}

		// Make it a valid RootCA
		tlsRemoteCert.IsCA = true
		tlsRemoteCert.KeyUsage = x509.KeyUsageCertSign

		// Setup the pool
		caCertPool.AddCert(tlsRemoteCert)
		tlsConfig.RootCAs = caCertPool

		// Set the ServerName
		if tlsRemoteCert.DNSNames != nil {
			tlsConfig.ServerName = tlsRemoteCert.DNSNames[0]
		}
	}

	tlsConfig.BuildNameToCertificate()
}

// GetTLSConfig returns a tls config
func GetTLSConfig(tlsClientCertFile, tlsClientKeyFile, tlsClientCAFile string, tlsRemoteCert *x509.Certificate) (*tls.Config, error) {
	tlsConfig := InitTLSConfig()

	// Client authentication
	if tlsClientCertFile != "" && tlsClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(tlsClientCertFile, tlsClientKeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if tlsClientCAFile != "" {
		caCertificates, err := ioutil.ReadFile(tlsClientCAFile)
		if err != nil {
			return nil, err
		}

		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCertificates)

		tlsConfig.RootCAs = caPool
	}

	finalizeTLSConfig(tlsConfig, tlsRemoteCert)
	return tlsConfig, nil
}

// GetTLSConfigMem returns a tls config
func GetTLSConfigMem(tlsClientCert string, tlsClientKey string, tlsClientCA string, tlsRemoteCertPEM string, insecureSkipVerify bool) (*tls.Config, error) {
	tlsConfig := InitTLSConfig()
	tlsConfig.InsecureSkipVerify = insecureSkipVerify
	// Client authentication
	if tlsClientCert != "" && tlsClientKey != "" {
		cert, err := tls.X509KeyPair([]byte(tlsClientCert), []byte(tlsClientKey))
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	var tlsRemoteCert *x509.Certificate
	if tlsRemoteCertPEM != "" {
		// Ignore any content outside of the PEM bytes we care about
		certBlock, _ := pem.Decode([]byte(tlsRemoteCertPEM))
		if certBlock == nil {
			return nil, fmt.Errorf("Invalid remote certificate")
		}

		var err error
		tlsRemoteCert, err = x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, err
		}
	}

	if tlsClientCA != "" {
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM([]byte(tlsClientCA))

		tlsConfig.RootCAs = caPool
	}

	finalizeTLSConfig(tlsConfig, tlsRemoteCert)

	return tlsConfig, nil
}

// WebsocketSendStream manages the send stream of the websocket
func WebsocketSendStream(conn *websocket.Conn, r io.Reader, bufferSize int) chan bool {
	ch := make(chan bool)

	if r == nil {
		close(ch)
		return ch
	}

	go func(conn *websocket.Conn, r io.Reader) {
		in := ReaderToChannel(r, bufferSize)
		for {
			buf, ok := <-in
			if !ok {
				break
			}

			w, err := conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("Got error getting next writer %s", err)
				break
			}

			_, err = w.Write(buf)
			w.Close()
			if err != nil {
				log.Printf("Got err writing %s", err)
				break
			}
		}
		conn.WriteMessage(websocket.TextMessage, []byte{})
		ch <- true
	}(conn, r)

	return ch
}

// WebsocketRecvStream manages the recv stream of the socket
func WebsocketRecvStream(w io.Writer, conn *websocket.Conn) chan bool {
	ch := make(chan bool)

	go func(w io.Writer, conn *websocket.Conn) {
		for {
			mt, r, err := conn.NextReader()
			if mt == websocket.CloseMessage {
				log.Printf("Got close message for reader")
				break
			}

			if mt == websocket.TextMessage {
				log.Printf("got message barrier")
				break
			}

			if err != nil {
				switch err.(type) {
				case *websocket.CloseError:
				default:
					log.Printf("Got error getting next reader %s, %s", err, w)
				}
				break
			}

			buf, err := ioutil.ReadAll(r)
			if err != nil {
				log.Printf("Got error writing to writer %s", err)
				break
			}

			if w == nil {
				continue
			}

			i, err := w.Write(buf)
			if i != len(buf) {
				log.Printf("Didn't write all of buf")
				break
			}
			if err != nil {
				log.Printf("Error writing buf %s", err)
				break
			}
		}
		ch <- true
	}(w, conn)

	return ch
}

// WebsocketProxy proxies a websocket connection
func WebsocketProxy(source *websocket.Conn, target *websocket.Conn) chan bool {
	forward := func(in *websocket.Conn, out *websocket.Conn, ch chan bool) {
		for {
			mt, r, err := in.NextReader()
			if err != nil {
				break
			}

			w, err := out.NextWriter(mt)
			if err != nil {
				break
			}

			_, err = io.Copy(w, r)
			w.Close()
			if err != nil {
				break
			}
		}

		ch <- true
	}

	chSend := make(chan bool)
	go forward(source, target, chSend)

	chRecv := make(chan bool)
	go forward(target, source, chRecv)

	ch := make(chan bool)
	go func() {
		select {
		case <-chSend:
		case <-chRecv:
		}

		source.Close()
		target.Close()

		ch <- true
	}()

	return ch
}

func defaultReader(conn *websocket.Conn, r io.ReadCloser, readDone chan<- bool) {
	/* For now, we don't need to adjust buffer sizes in
	* WebsocketMirror, since it's used for interactive things like
	* exec.
	 */
	in := ReaderToChannel(r, -1)
	for {
		buf, ok := <-in
		if !ok {
			r.Close()
			log.Printf("sending write barrier")
			conn.WriteMessage(websocket.TextMessage, []byte{})
			readDone <- true
			return
		}
		w, err := conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("Got error getting next writer %s", err)
			break
		}

		_, err = w.Write(buf)
		w.Close()
		if err != nil {
			log.Printf("Got err writing %s", err)
			break
		}
	}
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	conn.WriteMessage(websocket.CloseMessage, closeMsg)
	readDone <- true
	r.Close()
}

// ReaderToChannel reads from websocket and sends to a channel
func ReaderToChannel(r io.Reader, bufferSize int) <-chan []byte {
	if bufferSize <= 128*1024 {
		bufferSize = 128 * 1024
	}

	ch := make(chan ([]byte))

	go func() {
		readSize := 128 * 1024
		offset := 0
		buf := make([]byte, bufferSize)

		for {
			read := buf[offset : offset+readSize]
			nr, err := r.Read(read)
			offset += nr
			if offset > 0 && (offset+readSize >= bufferSize || err != nil) {
				ch <- buf[0:offset]
				offset = 0
				buf = make([]byte, bufferSize)
			}

			if err != nil {
				close(ch)
				break
			}
		}
	}()

	return ch
}

func defaultWriter(conn *websocket.Conn, w io.WriteCloser, writeDone chan<- bool) {
	for {
		mt, r, err := conn.NextReader()
		if err != nil {
			log.Printf("Got error getting next reader %s, %s", err, w)
			break
		}

		if mt == websocket.CloseMessage {
			log.Printf("Got close message for reader")
			break
		}

		if mt == websocket.TextMessage {
			log.Printf("Got message barrier, resetting stream")
			break
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("Got error writing to writer %s", err)
			break
		}
		i, err := w.Write(buf)
		if i != len(buf) {
			log.Printf("Didn't write all of buf")
			break
		}
		if err != nil {
			log.Printf("Error writing buf %s", err)
			break
		}
	}
	writeDone <- true
	w.Close()
}

// WebSocketMirrorReader mirror reader
type WebSocketMirrorReader func(conn *websocket.Conn, r io.ReadCloser, readDone chan<- bool)

// WebSocketMirrorWriter mirror writer
type WebSocketMirrorWriter func(conn *websocket.Conn, w io.WriteCloser, writeDone chan<- bool)

// WebsocketMirror allows mirroring a reader to a websocket and taking the
// result and writing it to a writer. This function allows for multiple
// mirrorings and correctly negotiates stream endings. However, it means any
// websocket.Conns passed to it are live when it returns, and must be closed
// explicitly.
func WebsocketMirror(conn *websocket.Conn, w io.WriteCloser, r io.ReadCloser, Reader WebSocketMirrorReader, Writer WebSocketMirrorWriter) (chan bool, chan bool) {
	readDone := make(chan bool, 1)
	writeDone := make(chan bool, 1)

	ReadFunc := Reader
	if ReadFunc == nil {
		ReadFunc = defaultReader
	}

	WriteFunc := Writer
	if WriteFunc == nil {
		WriteFunc = defaultWriter
	}

	go ReadFunc(conn, r, readDone)
	go WriteFunc(conn, w, writeDone)

	return readDone, writeDone
}

// WebsocketConsoleMirror console mirror
func WebsocketConsoleMirror(conn *websocket.Conn, w io.WriteCloser, r io.ReadCloser) (chan bool, chan bool) {
	readDone := make(chan bool, 1)
	writeDone := make(chan bool, 1)

	go defaultWriter(conn, w, writeDone)

	go func(conn *websocket.Conn, r io.ReadCloser) {
		in := ReaderToChannel(r, -1)
		for {
			buf, ok := <-in
			if !ok {
				r.Close()
				log.Printf("sending write barrier")
				conn.WriteMessage(websocket.TextMessage, []byte{})
				readDone <- true
				return
			}
			w, err := conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("Got error getting next writer %s", err)
				break
			}

			_, err = w.Write(buf)
			w.Close()
			if err != nil {
				log.Printf("Got err writing %s", err)
				break
			}
		}
		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		conn.WriteMessage(websocket.CloseMessage, closeMsg)
		readDone <- true
		r.Close()
	}(conn, r)

	return readDone, writeDone
}

// WebsocketUpgrader websocket.Upgrader
var WebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// AllocatePort asks the kernel for a free open port that is ready to use
func AllocatePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
