package router

import (
    "github.com/gin-gonic/gin"
    "github.com/CamiloLeonP/parking-radar/internal/app/handler"
)

func SetupRouter() *gin.Engine {
    
  router := gin.Default()

  router.POST("/ping", handler.PinHandler)

  router.GET("/init", handler.AuthMiddleware(), handler.InitHandler)
    
  return router
}