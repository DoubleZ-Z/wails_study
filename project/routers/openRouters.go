package routers

import (
	"wails_study/project/controllers/open"
	"wails_study/project/interceptor"

	"github.com/gin-gonic/gin"
)

func OpenRouters(engine *gin.Engine) {
	openGroup := engine.Group("/api/open")
	{
		userController := open.UserController{}
		openGroup.GET("/", userController.Home)
		openGroup.GET("/login", interceptor.TokenAuth, userController.Login)
		openGroup.GET("/find-all", interceptor.TokenAuth, userController.FindAll)
	}
}
