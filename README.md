[![Build Status](https://drone.io/github.com/mix3/go-irc/status.png)](https://drone.io/github.com/mix3/go-irc/latest)

# go-irc

## SYNOPSIS

```
import (
	"log"

	"github.com/mix3/go-irc"
)

func main() {
	conn, err := irc.New(&irc.Config{
		Nick:     "ikusan",
		User:     "ikusan",
		RealName: "ikusan",
		Callback: func(conn *irc.Conn, e *irc.Event) {
			switch e.Code {
			case "001":
				conn.Ping("ping")
			case "PONG":
				conn.Disconnect()
			case irc.DISCONNECTED:
				conn.Logger().Infof("[INFO   ] disconnect")
			}
		},
	})
	if err != nil {
		log.Fatalf("error new %s", err)
	}
	quit, err := conn.Connect("irc.freenode.net", 6667)
	if err != nil {
		log.Fatalf("error connect %s", err)
	}
	count := 1
	for {
		select {
		case <-quit:
			if 3 <= count {
				return
			}
			for {
				_, err = conn.Reconnect()
				if err == nil {
					count++
					break
				}
			}
		}
	}
}
```
