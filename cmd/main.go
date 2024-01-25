package main

import (
	"github.com/RomanTykhyi/students-api/config"
	"github.com/RomanTykhyi/students-api/internal/common"
	"github.com/RomanTykhyi/students-api/internal/server"
)

func init() {
	common.Init()
}

func main() {
	server.StartServer(config.GetAppConfig().HttpPort)
}
