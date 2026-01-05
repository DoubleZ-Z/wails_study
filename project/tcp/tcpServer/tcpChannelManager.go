package tcpServer

import (
	"fmt"
	"net"
	"time"
	"wails_study/project/logger"
	"wails_study/project/tcp/broker"
	"wails_study/project/tcp/dto"
	"wails_study/project/tcp/history"
	"wails_study/project/tcp/packetV2"
	"wails_study/project/util"
)

const (
	HEADER = 0x02
	FCS    = 0x03
)

type Manager struct {
	tcpChannelMap map[string]net.TCPAddr
}

func HandleConnect(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("连接关闭错误:", err)
		}
		fmt.Println("连接已关闭并释放资源")
	}()

	connect := dto.TcpConnect{
		Conn: conn,
		Id:   util.GenerateUniqueID(),
	}

	connectionAbstract := GetTcpConnectionAbstract()
	connectionAbstract.AddOpenChannel(connect)

	readBuffer := GetBufferPool().GetBuffer()

	var dataBuffer []byte

	for {
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		if err != nil {
			logger.Warn("设置连接读取超时失败:", err.Error())
			break
		}

		n, err := conn.Read(readBuffer)
		if err != nil {
			logger.Warn("读取数据错误或连接已断开:", err.Error())
			break
		}

		dataBuffer = append(dataBuffer, readBuffer[:n]...)

		processed := 0
		for {
			startIndex := -1
			endIndex := -1

			for i := processed; i < len(dataBuffer); i++ {
				if dataBuffer[i] == HEADER && startIndex == -1 {
					startIndex = i
				} else if dataBuffer[i] == FCS && startIndex != -1 {
					endIndex = i
					break
				}
			}

			if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
				dataPacket := dataBuffer[startIndex+1 : endIndex]
				logger.Infof("receive packet from %s : %s", connect.Conn.RemoteAddr().String(), string(dataPacket))
				defaultBroker := broker.DefaultBroker{}
				var result string
				packet, deserializeErr := packetV2.Deserialize(dataPacket)
				history.SetMessageHistory(string(dataPacket), packet, deserializeErr)
				if deserializeErr == nil {
					result = defaultBroker.HandlePacket(packet, connect)
				} else {
					logger.Warnf("packet deserialize error: %s", deserializeErr.Error())
				}
				processed = endIndex + 1
				startIndex = -1
				endIndex = -1
				if util.IsStringEmpty(result) {
					continue
				}
				logger.Infof("try send to %s response payload : %s", connect.Conn.RemoteAddr().String(), result)
				packetBytes := []byte(result)
				if _, writeErr := conn.Write(supplement(packetBytes)); writeErr != nil {
					logger.Errorf("发送响应数据到客户端失败: %s", writeErr.Error())
					break
				} else {
					logger.Infof("send response to %s success", connect.Conn.RemoteAddr().String())
				}
			} else {
				// 没有找到完整的数据包，将已处理的数据部分移除
				if processed > 0 {
					dataBuffer = dataBuffer[processed:]
					processed = 0
				}
				break
			}
		}
	}

	connectionAbstract.RemoveChannel(connect.Id)
}

func supplement(packetBytes []byte) []byte {
	responsePacket := make([]byte, 0)
	responsePacket = append(responsePacket, HEADER)
	responsePacket = append(responsePacket, packetBytes...)
	responsePacket = append(responsePacket, FCS)
	return responsePacket
}
