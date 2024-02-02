package main

import autoportforward "autoportforward/internal"

func main() {
	config := autoportforward.GetConfig()
	service := autoportforward.NewService(config)
	service.StartService()
}
