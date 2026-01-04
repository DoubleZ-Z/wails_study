package tcpServer

import (
	"fmt"
	"net"
	"time"
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
			fmt.Println("设置连接读取超时失败:", err)
			break
		}

		n, err := conn.Read(readBuffer)
		if err != nil {
			fmt.Println("读取数据错误或连接已断开:", err)
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
				fmt.Printf("receive packet from %d : %s\n", connect.Conn.RemoteAddr().Network(), string(dataPacket))
				defaultBroker := broker.DefaultBroker{}
				var result string
				packet, err := packetV2.Deserialize(dataPacket)
				history.SetMessageHistory(string(dataPacket), packet, err)
				if err != nil {
					result = defaultBroker.HandlePacket(packet, connect)
				} else {

				}
				processed = endIndex + 1
				startIndex = -1
				endIndex = -1
				if util.IsStringEmpty(result) {
					continue
				}
				fmt.Printf("try send response payload : %s \n", result)
				packetBytes := []byte(result)
				if _, err := conn.Write(supplement(packetBytes)); err != nil {
					fmt.Printf("发送响应数据到客户端失败: %v\n", err)
					break
				}
				fmt.Printf("send response to %d success\n", connect.Conn.RemoteAddr())
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

	// 连接断开时从连接管理器中移除
	connectionAbstract.RemoveChannel(connect.Id)
}

func supplement(packetBytes []byte) []byte {
	responsePacket := make([]byte, 0)
	responsePacket = append(responsePacket, HEADER)
	responsePacket = append(responsePacket, packetBytes...)
	responsePacket = append(responsePacket, FCS)
	return responsePacket
}
