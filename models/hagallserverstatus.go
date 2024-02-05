package models

import "errors"

// A Hagall server status.
type HagallServerStatus int

// Defined Hagall server status.
const (
	OfflineHagallServer HagallServerStatus = iota
	OnlineHagallServer
	UnhealthyHagallServer
	MaxServerStatus
)

func HagallServerStatusFromString(v string) (HagallServerStatus, error) {
	switch v {
	case "offline":
		return OfflineHagallServer, nil
	case "online":
		return OnlineHagallServer, nil
	case "unhealthy":
		return UnhealthyHagallServer, nil
	}
	return OfflineHagallServer, errors.New("invalid server status")
}

func (h HagallServerStatus) String() string {
	switch h {
	case OfflineHagallServer:
		return "offline"
	case OnlineHagallServer:
		return "online"
	case UnhealthyHagallServer:
		return "unhealthy"
	}
	return "offline"
}
