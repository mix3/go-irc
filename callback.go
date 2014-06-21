package irc

const (
	REGISTER     = "REGISTER"
	DISCONNECTED = "DISCONNECTED"
)

func (conn *Conn) callback(e *Event) {
	if conn.embedded != nil {
		conn.embedded.Callback(e)
	}
}

func (conn *Conn) DefaultCallback(e *Event) {
	switch e.Code {
	case REGISTER:
		if 0 < len(conn.password) {
			conn.Pass(conn.password)
		}
		conn.Nick(conn.nick)
		conn.User(conn.cfg.User, conn.cfg.RealName)
	case "001":
		conn.cnick = e.Args[0]
	case "433":
		conn.cnick = conn.cnick + "_"
		conn.Nick(conn.cnick)
	case "437":
		conn.cnick = conn.cnick + "_"
		conn.Nick(conn.cnick)
	case "NICK":
		conn.cnick = e.Message()
	case "PING":
		conn.Pong(e.Message())
	case DISCONNECTED:
	}
}
