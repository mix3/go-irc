package irc

import "strings"

// copy from github.com/fluffle/goirc/command

func (conn *Conn) Raw(rawline string) {
	conn.out <- rawline
}

func (conn *Conn) Pass(password string) {
	conn.out <- "PASS " + password
}

func (conn *Conn) Nick(nick string) {
	conn.nick = nick
	conn.out <- "NICK " + nick
}

func (conn *Conn) User(user, realname string) {
	conn.out <- "USER " + user + " 12 * :" + realname
}

func (conn *Conn) Join(channel string) {
	conn.out <- "JOIN " + channel
}

func (conn *Conn) Part(channel string, message ...string) {
	msg := strings.Join(message, " ")
	if msg != "" {
		msg = " :" + msg
	}
	conn.out <- "PART " + channel + msg
}

func (conn *Conn) Kick(channel, nick string, message ...string) {
	msg := strings.Join(message, " ")
	if msg != "" {
		msg = " :" + msg
	}
	conn.out <- "KICK" + " " + channel + " " + nick + msg
}

func (conn *Conn) Quit(message ...string) {
	msg := strings.Join(message, " ")
	if msg == "" {
		msg = "GoBye!"
	}
	conn.out <- "QUIT :" + msg
}

func (conn *Conn) Whois(nick string) {
	conn.out <- "WHOIS " + nick
}

func (conn *Conn) Who(nick string) {
	conn.out <- "WHO " + nick
}

func (conn *Conn) Privmsg(t, msg string) {
	conn.out <- "PRIVMSG " + t + " :" + msg
}

func (conn *Conn) Notice(t, msg string) {
	conn.out <- "NOTICE " + t + " :" + msg
}

func (conn *Conn) Ctcp(t, ctcp string, arg ...string) {
	msg := strings.Join(arg, " ")
	if msg != "" {
		msg = " " + msg
	}
	conn.Privmsg(t, "\001"+strings.ToUpper(ctcp)+msg+"\001")
}

func (conn *Conn) CtcpReply(t, ctcp string, arg ...string) {
	msg := strings.Join(arg, " ")
	if msg != "" {
		msg = " " + msg
	}
	conn.Notice(t, "\001"+strings.ToUpper(ctcp)+msg+"\001")
}

func (conn *Conn) Version(t string) {
	conn.Ctcp(t, "VERSION")
}

func (conn *Conn) Action(t, msg string) {
	conn.Ctcp(t, "ACTION", msg)
}

func (conn *Conn) Topic(channel string, topic ...string) {
	t := strings.Join(topic, " ")
	if t != "" {
		t = " :" + t
	}
	conn.out <- "TOPIC " + channel + t
}

func (conn *Conn) Mode(t string, modestring ...string) {
	mode := strings.Join(modestring, " ")
	if mode != "" {
		mode = " " + mode
	}
	conn.out <- "MODE " + t + mode
}

func (conn *Conn) Away(message ...string) {
	msg := strings.Join(message, " ")
	if msg != "" {
		msg = " :" + msg
	}
	conn.out <- "AWAY" + msg
}

func (conn *Conn) Invite(nick, channel string) {
	conn.out <- "INVITE " + nick + " " + channel
}

func (conn *Conn) Oper(user, pass string) {
	conn.out <- "OPER " + user + " " + pass
}

func (conn *Conn) Ping(msg string) {
	conn.out <- "PING :" + msg
}

func (conn *Conn) Pong(msg string) {
	conn.out <- "PONG :" + msg
}
