package packetV2

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Sign(packet ProtonPacket[any], signKey string) string {
	var items []string
	items = append(items, "type="+packet.Type)
	items = append(items, "protocolVer="+strconv.Itoa(packet.ProtocolVer))
	items = append(items, "reason="+packet.Reason)
	header := packet.Header
	items = append(items, "action="+header.Action)
	items = append(items, "actionVer="+header.ActionVer)
	items = append(items, "trace="+header.Trace)
	items = append(items, "timestamp="+header.Timestamp)
	items = append(items, "priority="+strconv.Itoa(header.Priority))
	sort.Strings(items)
	builder := strings.Builder{}
	var first = true
	for _, item := range items {
		if first {
			first = false
		} else {
			builder.WriteString("&")
		}
		builder.WriteString(item)
	}
	builder.WriteString(signKey)
	signText := builder.String()
	fmt.Println(signText)
	hasher := md5.New()
	hasher.Write([]byte(signText))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}
