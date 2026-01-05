package handler

import (
	"time"
	"wails_study/project/logger"
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/packetV2"
	"wails_study/project/util"
)

type HeartbeatHandler struct {
}

func (h *HeartbeatHandler) HandleAction(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) any {
	logger.Debugf("HeartbeatHandler ======================= [%s]", util.UnixToTimestampString(time.Now().Unix()))
	return nil
}
