package models

import (
	"net"
	"time"
)

type Event struct {
	DataID   []byte
	Author   []byte
	IP       net.IPAddr
	DateTime time.Time
}

type History struct {
	Events []Event
}
