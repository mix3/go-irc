package irc_test

import (
	"testing"

	"github.com/mix3/go-irc"
)

func TestConnection(t *testing.T) {
	conn, err := irc.New(&irc.Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
	})
	conn.CallbackerFunc(func(conn *irc.Conn, e *irc.Event) {
		irc.DefaultCallback(conn, e)
		switch e.Code {
		case "001":
			conn.Join("#ikusan_test")
		case "366":
			conn.Privmsg("#ikusan_test", "Test Message")
			conn.Nick("ikusan-newnick")
		case "NICK":
			conn.Disconnect()
		case irc.DISCONNECTED:
			conn.Logger().Debugf("[INFO   ] disconnect")
		}
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
	cb := irc.CallbackFunc(func(conn *irc.Conn, e *irc.Event) {
		irc.DefaultCallback(conn, e)
		switch e.Code {
		case "001":
			conn.Disconnect()
		case irc.DISCONNECTED:
			conn.Logger().Debugf("[INFO   ] disconnect")
		}
	})
	conn, err := irc.New(&irc.Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
	})
	conn.Callbacker(cb)
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
