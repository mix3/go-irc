package irc

import (
	"errors"
	"fmt"
	"strings"
)

type Event struct {
	Code   string
	Raw    string
	Nick   string
	Host   string
	Source string
	User   string
	Args   []string
}

func ParseMessage(msg string) (*Event, error) {
	msg = msg[:len(msg)-2]
	event := &Event{Raw: msg}
	if msg[0] == ':' {
		if i := strings.Index(msg, " "); i > -1 {
			event.Source = msg[1:i]
			msg = msg[i+1 : len(msg)]
		} else {
			return nil, errors.New(fmt.Sprintf("Misformed msg from server: %#s", msg))
		}

		if i, j := strings.Index(event.Source, "!"), strings.Index(event.Source, "@"); i > -1 && j > -1 {
			event.Nick = event.Source[0:i]
			event.User = event.Source[i+1 : j]
			event.Host = event.Source[j+1 : len(event.Source)]
		}
	}

	split := strings.SplitN(msg, " :", 2)
	args := strings.Split(split[0], " ")
	event.Code = strings.ToUpper(args[0])
	event.Args = args[1:]
	if len(split) > 1 {
		event.Args = append(event.Args, split[1])
	}

	event.ctcp()

	return event, nil
}

func (e *Event) Message() string {
	if len(e.Args) == 0 {
		return ""
	}
	return e.Args[len(e.Args)-1]
}

func (e *Event) ctcp() {
	msg := e.Message()
	if e.Code == "PRIVMSG" && len(msg) > 0 && msg[0] == '\x01' {
		e.Code = "CTCP"

		if i := strings.LastIndex(msg, "\x01"); i > -1 {
			msg = msg[1:i]
		}

		if msg == "VERSION" {
			e.Code = "CTCP_VERSION"
		} else if msg == "TIME" {
			e.Code = "CTCP_TIME"
		} else if msg[0:4] == "PING" {
			e.Code = "CTCP_PING"
		} else if msg == "USERINFO" {
			e.Code = "CTCP_USERINFO"
		} else if msg == "CLIENTINFO" {
			e.Code = "CTCP_CLIENTINFO"
		} else if msg[0:6] == "ACTION" {
			e.Code = "CTCP_ACTION"
			msg = msg[7:]
		}

		e.Args[len(e.Args)-1] = msg
	}
}
