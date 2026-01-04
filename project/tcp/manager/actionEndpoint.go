package manager

import (
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/handler"
	"wails_study/project/tcp/packetV2"
)

type ActionEndpointInterface interface {
	HandleAction(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) any
}

type ActionEndpoint struct {
	Action        string
	ActionVersion string
}

const (
	HEARTBEAT = "heartbeat.upload"
)

const (
	ACTION_VERSION_1 = "V0.1"
	ACTION_VERSION_2 = "V0.2"
)

func init() {
	RegisterAction(HEARTBEAT, ACTION_VERSION_1, &handler.HeartbeatHandler{})
}

var handlerMap = make(map[string]map[string]ActionEndpointInterface)

func RegisterAction(action string, actionVersion string, handler ActionEndpointInterface) {
	if _, ok := handlerMap[action]; !ok {
		handlerMap[action] = make(map[string]ActionEndpointInterface)
	}
	handlerMap[action][actionVersion] = handler
}

func GetActionHandler(action string, actionVersion string) ActionEndpointInterface {
	if actionHandlers, ok := handlerMap[action]; ok {
		if activeHandler, ok := actionHandlers[actionVersion]; ok {
			return activeHandler
		}
	}
	return nil
}
