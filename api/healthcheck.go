package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/user2410/simplebank/util"
)

func (s *Server) healthcheck(ctx *gin.Context) {
	ip, err := util.ExternalIPv4()
	if err != nil {
		log.Println("Error while getting external IP:", err)
		ctx.String(200, "pong")
	} else {
		ctx.String(200, fmt.Sprintf("pong from %s", string(ip)))
	}
}
