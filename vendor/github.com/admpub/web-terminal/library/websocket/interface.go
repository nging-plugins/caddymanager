package websocket

// Writer websocket writer
type Writer interface {
	WriteMessage(int, []byte) error
	WriteJSON(v interface{}) error
}

// Reader websocket reader
type Reader interface {
	ReadMessage() (int, []byte, error)
}
