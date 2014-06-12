package irc

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	golog "github.com/umisama/golog"
)

type Config struct {
	Nick      string
	User      string
	RealName  string
	SSL       bool
	SSLConfig *tls.Config
	Interval  time.Duration
	PingFreq  time.Duration
	Callback  func(*Conn, *Event)
	Logger    golog.Logger
}

func (c *Config) IsValid() bool {
	if c.Nick == "" {
		return false
	}
	return true
}

func (c *Config) Setup() {
	if c.Interval <= 0 {
		c.Interval = 2 * time.Second
	}
	if c.PingFreq <= 0 {
		c.PingFreq = 15 * time.Minute
	}
	if c.User == "" {
		c.User = c.Nick
	}
	if c.RealName == "" {
		c.RealName = c.User
	}
	if c.Logger == nil {
		level := golog.LogLevel_Info
		if os.Getenv("IKUSAN_DEBUG") != "" {
			level = golog.LogLevel_Debug
		}
		var err error
		c.Logger, err = golog.NewLogger(
			os.Stdout,
			golog.TIME_FORMAT_SEC,
			golog.LOG_FORMAT_SIMPLE,
			level,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
