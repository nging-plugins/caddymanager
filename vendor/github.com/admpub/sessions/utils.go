package sessions

import "encoding/gob"

func RegisterGob(v interface{}) {
	defer recover()
	gob.Register(v)
}
