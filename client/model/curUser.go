package model

import (
	"gin_example/chartroom/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
