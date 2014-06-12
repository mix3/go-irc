package irc

import "testing"

func TestParseMessage(t *testing.T) {
	func() {
		raw := ":test!test@test.test PRIVMSG #Test :This is a test case"
		event, _ := ParseMessage(raw + "\r\n")
		expect(t, event.Nick, "test")
		expect(t, event.User, "test")
		expect(t, event.Host, "test.test")
		expect(t, event.Source, "test!test@test.test")
		expect(t, event.Args, []string{"#Test", "This is a test case"})
		expect(t, event.Code, "PRIVMSG")
		expect(t, event.Raw, raw)
	}()
	func() {
		raw := ":test!test@test.test PRIVMSG #Test :\x01ACTION is a test case\x01"
		event, _ := ParseMessage(raw + "\r\n")
		expect(t, event.Nick, "test")
		expect(t, event.User, "test")
		expect(t, event.Host, "test.test")
		expect(t, event.Source, "test!test@test.test")
		expect(t, event.Args, []string{"#Test", "is a test case"})
		expect(t, event.Code, "CTCP_ACTION")
		expect(t, event.Raw, raw)
	}()
}
