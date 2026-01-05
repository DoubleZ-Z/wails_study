package packetV2

import (
	"encoding/json"
	"time"
	"wails_study/project/util"
)

const (
	TYPE_REQUEST   = "request"
	TYPE_RESPONSE  = "response"
	REASON_COMMAND = "command"
	PROTOCOL_VER_2 = 2
)

type ProtonPacketHeader struct {
	Action    string `json:"action"`
	ActionVer string `json:"actionVer"`
	Trace     string `json:"trace"`
	Priority  int    `json:"priority"`
	Time      string `json:"time"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

type ProtonPacketExt struct {
	AppType     string `json:"appType"`
	AppVer      string `json:"appVer"`
	Station     string `json:"station"`
	Timestamp   string `json:"timestamp"`
	ResCode     int    `json:"resCode"`
	ResMsg      string `json:"resMsg"`
	WorksheetNo string `json:"worksheetNo"`
	AccountNo   string `json:"accountNo"`
}

type ProtonPacket[T any] struct {
	Type        string             `json:"type"`
	Reason      string             `json:"reason"`
	ProtocolVer int                `json:"protocolVer"`
	Header      ProtonPacketHeader `json:"header"`
	Ext         ProtonPacketExt    `json:"ext"`
	Payload     any                `json:"payload"`
}

func CheckSign(packet ProtonPacket[any]) (bool, string) {
	sign := Sign(packet, "levent8421")
	if sign == packet.Header.Sign {
		return true, sign
	}
	return false, sign
}

func Deserialize(payloadV2 []byte) (ProtonPacket[any], error) {

	var packet ProtonPacket[any]
	var tempPacket struct {
		Type        string             `json:"type"`
		Reason      string             `json:"reason"`
		ProtocolVer int                `json:"protocolVer"`
		Header      ProtonPacketHeader `json:"header"`
		Ext         ProtonPacketExt    `json:"ext"`
		Payload     json.RawMessage    `json:"payload"`
	}

	err := json.Unmarshal(payloadV2, &tempPacket)
	if err != nil {

		return ProtonPacket[any]{
			Type:        "error",
			Reason:      "json unmarshal error",
			ProtocolVer: 0,
			Header:      ProtonPacketHeader{},
			Ext:         ProtonPacketExt{},
			Payload:     nil,
		}, err
	}

	packet.Type = tempPacket.Type
	packet.Reason = tempPacket.Reason
	packet.ProtocolVer = tempPacket.ProtocolVer
	packet.Header = tempPacket.Header
	packet.Ext = tempPacket.Ext
	packet.SetPayload(tempPacket.Payload)

	return packet, nil
}

func (packet *ProtonPacket[T]) SetPayload(data any) {
	packet.Payload = data
}

func (packet *ProtonPacket[T]) GetPayload() any {
	return packet.Payload
}

func (packet *ProtonPacket[T]) toJsonObject() (T, error) {
	var result T
	if packet.Payload == nil {
		return result, nil
	}
	bytes, err := json.Marshal(packet.Payload)
	if err != nil {
		return result, err
	}
	err2 := json.Unmarshal(bytes, &result)
	if err2 != nil {
		return result, err2
	}
	return result, err
}

func Response(request ProtonPacket[any], payload any, statusCode int, message string) (ProtonPacket[any], error) {
	response := ProtonPacket[any]{}
	response.SetPayload(payload)
	response.Type = TYPE_RESPONSE
	response.ProtocolVer = PROTOCOL_VER_2
	response.Reason = request.Reason
	ext := ProtonPacketExt{
		ResCode: statusCode,
		ResMsg:  message,
	}
	response.Ext = ext
	header := ProtonPacketHeader{
		Action:    request.Header.Action,
		ActionVer: request.Header.ActionVer,
		Trace:     request.Header.Trace,
		Priority:  request.Header.Priority,
		Timestamp: util.UnixToTimestampString(time.Now().Unix()),
	}
	response.Header = header
	return response, nil
}

func Serialize(response ProtonPacket[any]) (string, error) {
	marshal, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
