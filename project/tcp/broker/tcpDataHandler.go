package broker

import (
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/history"
	"wails_study/project/tcp/manager"
	"wails_study/project/tcp/packetV2"
)

type Broker interface {
	HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) (string, error)
	SendAction(stationNo string, action string, actionVar string, content any) string
}

type DefaultBroker struct {
}

func (b *DefaultBroker) HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) string {
	var responseStr string
	switch packet.Type {
	case packetV2.TYPE_REQUEST:
		response := manager.OnRequest(packet, connect)
		responseStr, _ = packetV2.Serialize(response)
		history.UpdateMessageHistory(responseStr, response)
	case packetV2.TYPE_RESPONSE:
		//TODO 响应处理
	}
	return responseStr
}

func (b *DefaultBroker) SendAction(stationNo string, action string, actionVar string, content any) string {
	return ""
}
