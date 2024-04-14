package api

import "github.com/gin-gonic/gin"

func (s *Server) healthcheck(ctx *gin.Context) {
	ctx.String(200, "pong")
}
