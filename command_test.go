package irc

import "testing"

func TestCommand(t *testing.T) {
	out := make(chan string)
	conn := &Conn{out: out}
	// Raw
	func() {
		go func() {
			conn.Raw("PING str")
		}()
		ret := <-out
		expect(t, ret, "PING str")
	}()
	// Pass
	func() {
		go func() {
			conn.Pass("password")
		}()
		ret := <-out
		expect(t, ret, "PASS password")
	}()
	// Nick
	func() {
		go func() {
			conn.Nick("nickname")
		}()
		ret := <-out
		expect(t, ret, "NICK nickname")
	}()
	// Join
	func() {
		go func() {
			conn.Join("#channel")
		}()
		ret := <-out
		expect(t, ret, "JOIN #channel")
	}()
	// Part
	func() {
		go func() {
			conn.Part("#channel")
		}()
		ret := <-out
		expect(t, ret, "PART #channel")
		go func() {
			conn.Part("#channel", "message1", "message2")
		}()
		ret = <-out
		expect(t, ret, "PART #channel :message1 message2")
	}()
	// Kick
	func() {
		go func() {
			conn.Kick("#channel", "mix3")
		}()
		ret := <-out
		expect(t, ret, "KICK #channel mix3")
		go func() {
			conn.Kick("#channel", "mix3", "message1", "message2")
		}()
		ret = <-out
		expect(t, ret, "KICK #channel mix3 :message1 message2")
	}()
	// Quit
	func() {
		go func() {
			conn.Quit()
		}()
		ret := <-out
		expect(t, ret, "QUIT :GoBye!")
		go func() {
			conn.Quit("message1", "message2")
		}()
		ret = <-out
		expect(t, ret, "QUIT :message1 message2")
	}()
	// Whois
	func() {
		go func() {
			conn.Whois("mix3")
		}()
		ret := <-out
		expect(t, ret, "WHOIS mix3")
	}()
	// Who
	func() {
		go func() {
			conn.Who("mix3")
		}()
		ret := <-out
		expect(t, ret, "WHO mix3")
	}()
	// Privmsg
	func() {
		go func() {
			conn.Privmsg("#channel", "message")
		}()
		ret := <-out
		expect(t, ret, "PRIVMSG #channel :message")
	}()
	// Notice
	func() {
		go func() {
			conn.Notice("#channel", "message")
		}()
		ret := <-out
		expect(t, ret, "NOTICE #channel :message")
	}()
	// Ctcp
	func() {
		go func() {
			conn.Ctcp("#channel", "ACTION")
		}()
		ret := <-out
		expect(t, ret, "PRIVMSG #channel :\001ACTION\001")
		go func() {
			conn.Ctcp("#channel", "ACTION", "arg1", "arg2")
		}()
		ret = <-out
		expect(t, ret, "PRIVMSG #channel :\001ACTION arg1 arg2\001")
	}()
	// CtcpReply
	func() {
		go func() {
			conn.CtcpReply("#channel", "ACTION")
		}()
		ret := <-out
		expect(t, ret, "NOTICE #channel :\001ACTION\001")
		go func() {
			conn.CtcpReply("#channel", "ACTION", "arg1", "arg2")
		}()
		ret = <-out
		expect(t, ret, "NOTICE #channel :\001ACTION arg1 arg2\001")
	}()
	// Version
	func() {
		go func() {
			conn.Version("#channel")
		}()
		ret := <-out
		expect(t, ret, "PRIVMSG #channel :\001VERSION\001")
	}()
	// Action
	func() {
		go func() {
			conn.Action("#channel", "msg")
		}()
		ret := <-out
		expect(t, ret, "PRIVMSG #channel :\001ACTION msg\001")
	}()
	// Topic
	func() {
		go func() {
			conn.Topic("#channel")
		}()
		ret := <-out
		expect(t, ret, "TOPIC #channel")
		go func() {
			conn.Topic("#channel", "arg1", "arg2")
		}()
		ret = <-out
		expect(t, ret, "TOPIC #channel :arg1 arg2")
	}()
	// Mode
	func() {
		go func() {
			conn.Mode("#channel")
		}()
		ret := <-out
		expect(t, ret, "MODE #channel")
		go func() {
			conn.Mode("#channel", "arg1", "arg2")
		}()
		ret = <-out
		expect(t, ret, "MODE #channel arg1 arg2")
	}()
	// Away
	func() {
		go func() {
			conn.Away()
		}()
		ret := <-out
		expect(t, ret, "AWAY")
		go func() {
			conn.Away("arg1", "arg2")
		}()
		ret = <-out
		expect(t, ret, "AWAY :arg1 arg2")
	}()
	// Invite
	func() {
		go func() {
			conn.Invite("nick", "channel")
		}()
		ret := <-out
		expect(t, ret, "INVITE nick channel")
	}()
	// Oper
	func() {
		go func() {
			conn.Oper("nick", "channel")
		}()
		ret := <-out
		expect(t, ret, "OPER nick channel")
	}()
	// Ping
	func() {
		go func() {
			conn.Ping("message")
		}()
		ret := <-out
		expect(t, ret, "PING :message")
	}()
	// Pong
	func() {
		go func() {
			conn.Pong("message")
		}()
		ret := <-out
		expect(t, ret, "PONG :message")
	}()
}
