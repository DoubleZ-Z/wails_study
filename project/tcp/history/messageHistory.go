package history

import (
	"wails_study/project/tcp/packetV2"
	"wails_study/project/util"
)

type MessageHistory struct {
	TraceId         string
	StationNo       string
	Category        string
	Success         bool
	ResultCode      int
	ResultMessage   string
	Action          string
	RequestMessage  string
	ResponseMessage string
	RequestTime     string
	ResponseTime    string
	Duration        int
	RequestDelay    int
}

var messageHistoryMap *util.SafeMap[string, *MessageHistory]

func init() {
	messageHistoryMap = util.NewSafeMap[string, *MessageHistory]()
	SetMessageHistory("origin", packetV2.ProtonPacket[any]{
		Header: packetV2.ProtonPacketHeader{
			Action: "origin",
			Trace:  "0001",
			Time:   "1234124123",
		},
		Ext: packetV2.ProtonPacketExt{
			Station: "yanyizhi",
		},
		Reason: "event",
	}, nil)
}

func SetMessageHistory(origin string, packet packetV2.ProtonPacket[any], error error) {
	messageHistory := &MessageHistory{
		TraceId:        packet.Header.Trace,
		StationNo:      packet.Ext.Station,
		Category:       packet.Reason,
		Success:        error == nil,
		Action:         packet.Header.Action,
		RequestTime:    packet.Header.Time,
		RequestMessage: origin,
	}
	messageHistoryMap.Set(packet.Header.Trace, messageHistory)
}

func GetMessageHistoryList() []*MessageHistory {
	slice := messageHistoryMap.ToSlice()
	return slice
}

func ClearMessageHistory() {
	messageHistoryMap = util.NewSafeMap[string, *MessageHistory]()
}
