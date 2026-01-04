package handler

import (
	"fmt"
	"time"
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/packetV2"
	"wails_study/project/util"
)

type HeartbeatHandler struct {
}

func (h *HeartbeatHandler) HandleAction(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) any {
	fmt.Printf("HeartbeatHandler ======================= [%s]", util.UnixToTimestampString(time.Now().Unix()))
	return nil
}
