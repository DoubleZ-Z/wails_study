package routers

import (
	"wails_study/project/controllers/token"
	"wails_study/project/interceptor"

	"github.com/gin-gonic/gin"
)

func TokenRouters(engine *gin.Engine) {
	group := engine.Group("/api/token", interceptor.TokenAuth)
	{
		controller := token.FileController{}
		group.POST("/file/upload", controller.UploadFile)
	}
}
