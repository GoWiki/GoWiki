package models

import (
	"time"
	"net"
)

type Event struct {
	DataID []byte
	Author []byte
	IP net.IPAddr
	DateTime time.Time
}


type History stuct {
	Events []Event
}