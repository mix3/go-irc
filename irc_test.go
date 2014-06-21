package irc

import (
	"log"
	"testing"
)

type MyConnA struct {
	*Conn
}

func (conn *MyConnA) Callback(e *Event) {
	conn.DefaultCallback(e)
	switch e.Code {
	case "001":
		conn.Join("#ikusan_test")
	case "366":
		conn.Privmsg("#ikusan_test", "Test Message")
		conn.Nick("ikusan-newnick")
	case "NICK":
		conn.Disconnect()
	case DISCONNECTED:
		log.Printf("[INFO   ] disconnect")
	}
}

func TestConnection(t *testing.T) {
	conn, err := New(&Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
	})
	myconn := &MyConnA{conn}
	myconn.SetEmbed(myconn)
	if err != nil {
		t.Fatalf("error new %s", err)
	}
	quit, err := myconn.Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Fatalf("error connect %s", err)
	}
	<-quit
}

type MyConnB struct {
	*Conn
}

func (conn *MyConnB) Callback(e *Event) {
	conn.DefaultCallback(e)
	switch e.Code {
	case "001":
		conn.Disconnect()
	}
}

func TestReconnection(t *testing.T) {
	conn, err := New(&Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
	})
	myconn := &MyConnB{conn}
	myconn.SetEmbed(myconn)
	if err != nil {
		t.Fatalf("error new %s", err)
	}
	quit, err := myconn.Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Fatalf("error connect %s", err)
	}
	count := 1
	for {
		select {
		case <-quit:
			if 3 <= count {
				return
			}
			for {
				if _, err = myconn.Reconnect(); err == nil {
					count++
					break
				}
			}
		}
	}
}
