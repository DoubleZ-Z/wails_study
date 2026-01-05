package history

import (
	"strconv"
	"time"
	"wails_study/project/frontend"
	"wails_study/project/tcp/packetV2"
	"wails_study/project/util"
)

type MessageHistory struct {
	TraceId         string `json:"traceId"`
	StationNo       string `json:"stationNo"`
	Category        string `json:"category,omitempty"`
	Success         bool   `json:"success,omitempty"`
	ResultCode      int    `json:"resultCode,omitempty"`
	ResultMessage   string `json:"resultMessage,omitempty"`
	Action          string `json:"action,omitempty"`
	RequestMessage  string `json:"requestMessage,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
	RequestTime     string `json:"requestTime,omitempty"`
	ReceiveTime     string `json:"receiveTime,omitempty"`
	ResponseTime    string `json:"responseTime,omitempty"`
	Duration        int64  `json:"duration,omitempty"`
	RequestDelay    int64  `json:"requestDelay,omitempty"`
}

var messageHistoryMap *util.SafeMap[string, *MessageHistory]

func init() {
	messageHistoryMap = util.NewSafeMap[string, *MessageHistory]()
}

func SetMessageHistory(origin string, packet packetV2.ProtonPacket[any], error error) {
	var requestDelay int64
	timestamp, parseError := strconv.ParseInt(packet.Header.Timestamp, 10, 64)
	if parseError == nil {
		requestDelay = time.Now().UnixMilli() - timestamp
	}
	messageHistory := &MessageHistory{
		TraceId:        packet.Header.Trace,
		StationNo:      packet.Ext.Station,
		Category:       packet.Reason,
		Success:        error == nil,
		Action:         packet.Header.Action,
		RequestTime:    packet.Header.Time,
		ReceiveTime:    util.UnixToTimestampString(time.Now().UnixMilli()),
		RequestMessage: origin,
		RequestDelay:   requestDelay,
	}
	messageHistoryMap.Set(packet.Header.Trace, messageHistory)
	touchChange()
}

func GetMessageHistoryList() []*MessageHistory {
	slice := messageHistoryMap.ToSlice()
	return slice
}

func UpdateMessageHistory(origin string, packet packetV2.ProtonPacket[any]) {
	if messageHistory, ok := messageHistoryMap.Get(packet.Header.Trace); ok {
		responseTime := util.UnixToTimestampString(time.Now().UnixMilli())
		milliseconds := util.DiffMilliseconds(messageHistory.ReceiveTime, responseTime)
		messageHistory.ResponseMessage = origin
		messageHistory.ResponseTime = responseTime
		messageHistory.Success = true
		messageHistory.Duration = milliseconds
		touchChange()
	}
}

func ClearMessageHistory() {
	messageHistoryMap = util.NewSafeMap[string, *MessageHistory]()
}

func touchChange() {
	frontend.Emit("messageHistoryReflush", GetMessageHistoryList())
}
