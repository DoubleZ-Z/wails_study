package manager

import (
	"fmt"
	"net/http"
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/packetV2"
)

func OnRequest(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) packetV2.ProtonPacket[any] {
	ok, signText := packetV2.CheckSign(packet)
	header := packet.Header
	var response packetV2.ProtonPacket[any]
	var err error

	if !ok {
		response, err = packetV2.Response(packet, nil, http.StatusUnauthorized, "sign key check error!!!,sign : "+signText)
		return response
	}
	handler := GetActionHandler(header.Action, header.ActionVer)
	if handler == nil {
		err = fmt.Errorf("cannot find handler:action:[%s],actionVer:[%s]", header.Action, header.ActionVer)
	} else {
		payload := handler.HandleAction(packet, connect)
		response, err = packetV2.Response(packet, payload, http.StatusOK, "OK")
	}
	if err != nil {
		response, err = packetV2.Response(packet, nil, http.StatusInternalServerError, "handle action["+header.Action+"] err")
	}
	return response
}
