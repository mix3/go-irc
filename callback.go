package irc

const (
	REGISTER     = "REGISTER"
	DISCONNECTED = "DISCONNECTED"
)

type Callbacker interface {
	Callback(*Conn, *Event)
}

type CallbackFunc func(*Conn, *Event)

func (cf CallbackFunc) Callback(conn *Conn, e *Event) {
	cf(conn, e)
}

func (conn *Conn) CallbackerFunc(f func(*Conn, *Event)) {
	conn.callbacker = CallbackFunc(f)
}

func (conn *Conn) Callbacker(cf CallbackFunc) {
	conn.callbacker = cf
}

func (conn *Conn) callback(e *Event) {
	conn.callbacker.Callback(conn, e)
}

func DefaultCallback(conn *Conn, e *Event) {
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
