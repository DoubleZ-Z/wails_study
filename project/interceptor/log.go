package interceptor

import (
	"fmt"
	"time"
	"wails_study/project/util"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	cp := c.Copy()
	go fmt.Println(util.UnixToTimeString(time.Now().Unix()), cp.Request.Host, cp.Request.RequestURI, "这里异步记录了一个日志")
}
