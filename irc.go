package irc

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	golog "github.com/umisama/golog"
)

type Conn struct {
	sync.WaitGroup
	cfg        *Config
	sock       net.Conn
	nick       string
	cnick      string
	server     string
	password   string
	port       uint
	io         *bufio.ReadWriter
	in         chan *Event
	out        chan string
	end        chan struct{}
	shutdown   chan struct{}
	quit       chan struct{}
	err        chan error
	connected  bool
	callbacker Callbacker
}

func New(cfg *Config) (*Conn, error) {
	if !cfg.IsValid() {
		return nil, errors.New("Config is invalid")
	}
	cfg.Setup()
	conn := &Conn{
		cfg:      cfg,
		nick:     cfg.Nick,
		cnick:    cfg.Nick,
		in:       make(chan *Event, 32),
		out:      make(chan string, 32),
		shutdown: make(chan struct{}, 3),
		quit:     make(chan struct{}),
		err:      make(chan error),
	}
	return conn, nil
}

func (conn *Conn) recv() {
	defer conn.Done()
	for {
		msg, err := conn.io.ReadString('\n')
		conn.Logger().Debugf("[<- RECV] %s", chomp(msg))
		if err != nil {
			conn.err <- err
			return
		}
		event, err := ParseMessage(msg)
		if err != nil {
			conn.err <- err
			return
		}
		conn.in <- event
	}
}

func (conn *Conn) call() {
	defer conn.Done()
	for {
		select {
		case <-conn.end:
			return
		case event := <-conn.in:
			conn.callback(event)
		}
	}
}

func (conn *Conn) send() {
	defer conn.Done()
	for {
		select {
		case <-conn.end:
			return
		case line := <-conn.out:
			conn.Logger().Debugf("[SEND ->] %s", chomp(line))
			conn.write(line)
			time.Sleep(conn.cfg.Interval)
		}
	}
}

func (conn *Conn) write(line string) {
	if _, err := conn.io.WriteString(line + "\r\n"); err != nil {
		conn.err <- err
		return
	}
	if err := conn.io.Flush(); err != nil {
		conn.err <- err
		return
	}
}

func (conn *Conn) ping() {
	defer conn.Done()
	ticker := time.NewTicker(conn.cfg.PingFreq)
	for {
		select {
		case <-conn.end:
			ticker.Stop()
			return
		case <-ticker.C:
			conn.Raw(fmt.Sprintf("PING :%d", time.Now().UnixNano()))
		}
	}
}

func (conn *Conn) errs() {
	for {
		select {
		case err := <-conn.err:
			if err != nil {
				conn.Logger().Warnf("[ ERROR ] %v", err)
				conn.shutdown <- struct{}{}
			} else {
				return
			}
		}
	}
}

func (conn *Conn) down() {
	for {
		select {
		default:
			if !conn.connected {
				conn.quit <- struct{}{}
				return
			}
		case <-conn.shutdown:
			if conn.connected {
				conn.callback(&Event{Code: DISCONNECTED})
				conn.connected = false
				conn.sock.Close()
				close(conn.end)
				conn.Wait()
				close(conn.err)
				conn.sock = nil
				conn.io = nil
			}
		}
	}
}

func (conn *Conn) Connect(server string, port uint, password ...string) (chan struct{}, error) {
	host := fmt.Sprintf("%s:%d", server, port)

	if conn.cfg.SSL {
		if sock, err := tls.Dial("tcp", host, conn.cfg.SSLConfig); err == nil {
			conn.sock = sock
		} else {
			return nil, err
		}
	} else {
		if sock, err := net.Dial("tcp", host); err == nil {
			conn.sock = sock
		} else {
			conn.Logger().Warnf("%s", err)
			return nil, err
		}
	}

	conn.Logger().Infof("[INFO   ] connect to %s", host)

	conn.server = server
	conn.port = port
	if 0 < len(password) {
		conn.password = password[0]
	} else {
		conn.password = ""
	}

	conn.postConnect()

	return conn.quit, nil
}

func (conn *Conn) postConnect() {
	conn.io = bufio.NewReadWriter(
		bufio.NewReader(conn.sock),
		bufio.NewWriter(conn.sock),
	)

	conn.end = make(chan struct{})
	conn.err = make(chan error)
	conn.Add(4)
	go conn.call()
	go conn.send()
	go conn.recv()
	go conn.ping()

	conn.connected = true
	conn.callback(&Event{Code: REGISTER})
	go conn.errs()
	go conn.down()
}

func (conn *Conn) Disconnect() {
	conn.shutdown <- struct{}{}
}

func (conn *Conn) Reconnect() (chan struct{}, error) {
	return conn.Connect(conn.server, conn.port, conn.password)
}

func (conn *Conn) Logger() golog.Logger {
	return conn.cfg.Logger
}

func (conn *Conn) IsConnected() bool {
	return conn.connected
}

func chomp(msg string) string {
	return strings.TrimRight(msg, "\n")
}
