package connect

// DefaultServer ...
var DefaultServer *Server

// Connect ...
type Connect struct {
}

// New ...
func New() *Connect {
	return new(Connect)
}
