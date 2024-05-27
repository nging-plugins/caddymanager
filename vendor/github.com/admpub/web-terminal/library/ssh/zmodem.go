package ssh

import (
	"bytes"
	"io"

	"github.com/admpub/errors"
	"github.com/admpub/web-terminal/config"
	websocketx "github.com/admpub/web-terminal/library/websocket"
	"github.com/admpub/websocket"
	"golang.org/x/crypto/ssh"
)

var (
	// ZModemSZStart ...
	// sz fmt.Sprintf("%+q", "rz\r**\x18B00000000000000\r\x8a\x11")
	//ZModemSZStart = []byte{13, 42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}
	ZModemSZStart = []byte{42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}

	// ZModemSZEnd ...
	// sz 结束 fmt.Sprintf("%+q", "\r**\x18B0800000000022d\r\x8a")
	//ZModemSZEnd = []byte{13, 42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}
	ZModemSZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}

	// ZModemSZEndOO ...
	// sz 结束后可能还会发送两个 OO，但是经过测试发现不一定每次都会发送 fmt.Sprintf("%+q", "OO")
	ZModemSZEndOO = []byte{79, 79}

	// ZModemRZStart ...
	// rz fmt.Sprintf("%+q", "**\x18B0100000023be50\r\x8a\x11")
	ZModemRZStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 50, 51, 98, 101, 53, 48, 13, 138, 17}

	// ZModemRZEStart ...
	// rz -e fmt.Sprintf("%+q", "**\x18B0100000063f694\r\x8a\x11")
	ZModemRZEStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 54, 51, 102, 54, 57, 52, 13, 138, 17}

	// ZModemRZSStart ...
	// rz -S fmt.Sprintf("%+q", "**\x18B0100000223d832\r\x8a\x11")
	ZModemRZSStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 50, 51, 100, 56, 51, 50, 13, 138, 17}

	// ZModemRZESStart ...
	// rz -e -S fmt.Sprintf("%+q", "**\x18B010000026390f6\r\x8a\x11")
	ZModemRZESStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 54, 51, 57, 48, 102, 54, 13, 138, 17}

	// ZModemRZEnd ...
	// rz 结束 fmt.Sprintf("%+q", "**\x18B0800000000022d\r\x8a")
	ZModemRZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}

	// ZModemRZCtrlStart ...
	// **\x18B0
	ZModemRZCtrlStart = []byte{42, 42, 24, 66, 48}

	// ZModemRZCtrlEnd1 ...
	// \r\x8a\x11
	ZModemRZCtrlEnd1 = []byte{13, 138, 17}

	// ZModemRZCtrlEnd2 ...
	// \r\x8a
	ZModemRZCtrlEnd2 = []byte{13, 138}

	// ZModemCancel ...
	// zmodem 取消 \x18\x18\x18\x18\x18\x08\x08\x08\x08\x08
	ZModemCancel = []byte{24, 24, 24, 24, 24, 8, 8, 8, 8, 8}
)

var (
	msgSZDisabled = []byte("sz download is disabled")
	msgRZDisabled = []byte("rz upload is disabled")
)

// 发送 ssh 会话的 stdout 和 stdin 数据到 websocket 连接
func TransformChannel(session *ssh.Session, conn websocketx.Writer, cfg *config.TransformConfig) (stdout io.Reader, stderr io.Reader, stdin io.WriteCloser, err error) {
	if cfg == nil {
		err = errors.New(`config.TransformConfig can't be nil`)
		return
	}
	cfg.SetDefaults()
	stdout, err = session.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "get stdout channel error")
		return
	}
	stderr, err = session.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "get stderr channel error")
		return
	}
	stdin, err = session.StdinPipe()
	if err != nil {
		err = errors.Wrap(err, "get stdin channel error")
		return
	}

	transferHandler := func(t MessageType, r io.Reader, w io.WriteCloser) {
		buff := make([]byte, cfg.BufferSize)
		for {
			n, err := r.Read(buff)
			if err != nil {
				return
			}

			//res := fmt.Sprintf("%+q", string(buff[:n]))
			//fmt.Println(buff[:n])
			//fmt.Println(t, n, res)
			operateZModemBytes(n, buff, w, t, conn, cfg)
		}
	}
	go transferHandler(MessageTypeStdout, stdout, stdin)
	go transferHandler(MessageTypeStderr, stderr, stdin)
	return
}

func operateZModemBytes(n int, buff []byte, w io.WriteCloser, t MessageType, conn websocketx.Writer, tcfg *config.TransformConfig) {
	cfg := tcfg.ZModemConfig()
	if cfg.GetZModemSZOO() {
		cfg.SetZModemSZOO(false)
		if n < 2 {
			conn.WriteJSON(&Message{Type: t, Data: buff[:n]})
			return
		}
		if n == 2 {
			if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
				conn.WriteMessage(websocket.BinaryMessage, buff[:n])
			} else {
				conn.WriteJSON(&Message{Type: t, Data: buff[:n]})
			}
			return
		}
		if buff[0] == ZModemSZEndOO[0] && buff[1] == ZModemSZEndOO[1] {
			conn.WriteMessage(websocket.BinaryMessage, buff[:2])
			conn.WriteJSON(&Message{Type: t, Data: buff[2:n]})
		} else {
			conn.WriteJSON(&Message{Type: t, Data: buff[:n]})
		}
		return
	}
	if cfg.GetZModemSZ() {
		if n == tcfg.BufferSize {
			// 如果读取的长度为 buffsize，则认为是在传输数据，
			// 这样可以提高 sz 下载速率，很低概率会误判 zmodem 取消操作
			conn.WriteMessage(websocket.BinaryMessage, buff[:n])
			return
		}
		if x, ok := checkByteCommand(buff[:n], ZModemSZEnd); ok {
			cfg.SetZModemSZ(false)
			cfg.SetZModemSZOO(true)
			conn.WriteMessage(websocket.BinaryMessage, ZModemSZEnd)
			if len(x) != 0 {
				conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
			}
			return
		}
		if _, ok := checkByteCommand(buff[:n], ZModemCancel); ok {
			cfg.SetZModemSZ(false)
			conn.WriteMessage(websocket.BinaryMessage, buff[:n])
		} else {
			conn.WriteMessage(websocket.BinaryMessage, buff[:n])
		}
		return
	}
	if cfg.GetZModemRZ() {
		if x, ok := checkByteCommand(buff[:n], ZModemRZEnd); ok {
			cfg.SetZModemRZ(false)
			conn.WriteMessage(websocket.BinaryMessage, ZModemRZEnd)
			if len(x) != 0 {
				conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
			}
			return
		}
		if _, ok := checkByteCommand(buff[:n], ZModemCancel); ok {
			cfg.SetZModemRZ(false)
			conn.WriteMessage(websocket.BinaryMessage, buff[:n])
			return
		}

		// rz 上传过程中服务器端还是会给客户端发送一些信息，比如心跳
		//conn.WriteJSON(&message{Type: messageTypeConsole, Data: buff[:n]})
		//conn.WriteMessage(websocket.BinaryMessage, buff[:n])

		startIndex := bytes.Index(buff[:n], ZModemRZCtrlStart)
		if startIndex != -1 {
			endIndex := bytes.Index(buff[:n], ZModemRZCtrlEnd1)
			if endIndex != -1 {
				ctrl := append(ZModemRZCtrlStart, buff[startIndex+len(ZModemRZCtrlStart):endIndex]...)
				ctrl = append(ctrl, ZModemRZCtrlEnd1...)
				conn.WriteMessage(websocket.BinaryMessage, ctrl)
				info := append(buff[:startIndex], buff[endIndex+len(ZModemRZCtrlEnd1):n]...)
				if len(info) != 0 {
					conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: info})
				}
				return
			}
			endIndex = bytes.Index(buff[:n], ZModemRZCtrlEnd2)
			if endIndex != -1 {
				ctrl := append(ZModemRZCtrlStart, buff[startIndex+len(ZModemRZCtrlStart):endIndex]...)
				ctrl = append(ctrl, ZModemRZCtrlEnd2...)
				conn.WriteMessage(websocket.BinaryMessage, ctrl)
				info := append(buff[:startIndex], buff[endIndex+len(ZModemRZCtrlEnd2):n]...)
				if len(info) != 0 {
					conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: info})
				}
				return
			}
			conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: buff[:n]})
			return
		}
		conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: buff[:n]})
		return
	}
	if x, ok := checkByteCommand(buff[:n], ZModemSZStart); ok {
		if cfg.GetDisableZModemSZ() {
			conn.WriteJSON(&Message{Type: MessageTypeAlert, Data: msgSZDisabled})
			w.Write(ZModemCancel)
			return
		}
		if y, ok := checkByteCommand(x, ZModemCancel); ok {
			// 下载不存在的文件以及文件夹(zmodem 不支持下载文件夹)时
			conn.WriteJSON(&Message{Type: t, Data: y})
		} else {
			cfg.SetZModemSZ(true)
			if len(x) != 0 {
				conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
			}
			conn.WriteMessage(websocket.BinaryMessage, ZModemSZStart)
		}
		return
	}
	if x, ok := checkByteCommand(buff[:n], ZModemRZStart); ok {
		if cfg.GetDisableZModemRZ() {
			conn.WriteJSON(&Message{Type: MessageTypeAlert, Data: msgRZDisabled})
			w.Write(ZModemCancel)
			return
		}
		cfg.SetZModemRZ(true)
		if len(x) != 0 {
			conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
		}
		conn.WriteMessage(websocket.BinaryMessage, ZModemRZStart)
		return
	}
	if x, ok := checkByteCommand(buff[:n], ZModemRZEStart); ok {
		if cfg.GetDisableZModemRZ() {
			conn.WriteJSON(&Message{Type: MessageTypeAlert, Data: msgRZDisabled})
			w.Write(ZModemCancel)
			return
		}
		cfg.SetZModemRZ(true)
		if len(x) != 0 {
			conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
		}
		conn.WriteMessage(websocket.BinaryMessage, ZModemRZEStart)
		return
	}
	if x, ok := checkByteCommand(buff[:n], ZModemRZSStart); ok {
		if cfg.GetDisableZModemRZ() {
			conn.WriteJSON(&Message{Type: MessageTypeAlert, Data: msgRZDisabled})
			w.Write(ZModemCancel)
			return
		}
		cfg.SetZModemRZ(true)
		if len(x) != 0 {
			conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
		}
		conn.WriteMessage(websocket.BinaryMessage, ZModemRZSStart)
		return
	}
	if x, ok := checkByteCommand(buff[:n], ZModemRZESStart); ok {
		if cfg.GetDisableZModemRZ() {
			conn.WriteJSON(&Message{Type: MessageTypeAlert, Data: msgRZDisabled})
			w.Write(ZModemCancel)
			return
		}
		cfg.SetZModemRZ(true)
		if len(x) != 0 {
			conn.WriteJSON(&Message{Type: MessageTypeConsole, Data: x})
		}
		conn.WriteMessage(websocket.BinaryMessage, ZModemRZESStart)
		return
	}
	conn.WriteJSON(&Message{Type: t, Data: buff[:n]})
}

// checkByteCommand ...
func checkByteCommand(x, y []byte) (n []byte, contain bool) {
	index := bytes.Index(x, y)
	if index == -1 {
		return
	}
	lastIndex := index + len(y)
	n = append(x[:index], x[lastIndex:]...)
	return n, true
}

type MessageType string

type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
	Cols int         `json:"cols,omitempty"`
	Rows int         `json:"rows,omitempty"`
}

const (
	MessageTypeStdin   MessageType = "stdin"
	MessageTypeStdout  MessageType = "stdout"
	MessageTypeStderr  MessageType = "stderr"
	MessageTypeResize  MessageType = "resize"
	MessageTypeConsole MessageType = "console"
	MessageTypeAlert   MessageType = "alert"
)
