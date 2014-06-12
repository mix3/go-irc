package irc

import (
	"log"
	"testing"
)

func TestConnection(t *testing.T) {
	conn, err := New(&Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
		Callback: func(conn *Conn, e *Event) {
			switch e.Code {
			case "001":
				conn.Join("#ikusan_test")
			case "366":
				conn.Privmsg("#ikusan_test", "Test Message")
				conn.Nick("ikusan-newnick")
			case "NICK":
				conn.Disconnect()
			case DISCONNECT:
				log.Printf("[INFO   ] disconnect")
			}
		},
	})
	if err != nil {
		t.Fatalf("error new %s", err)
	}
	quit, err := conn.Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Fatalf("error connect %s", err)
	}
	<-quit
}

func TestReconnection(t *testing.T) {
	conn, err := New(&Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
		Callback: func(conn *Conn, e *Event) {
			switch e.Code {
			case "001":
				conn.Disconnect()
			}
		},
	})
	if err != nil {
		t.Fatalf("error new %s", err)
	}
	quit, err := conn.Connect("irc.freenode.net", 6667)
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
				if _, err = conn.Reconnect(); err == nil {
					count++
					break
				}
			}
		}
	}
}
