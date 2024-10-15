package main

import (
	"mi-LModel/net"
	"net/http"
)

func main() {
	//param := client.DidiChatCompletionNewParams()
	//
	//client.Demo(ctx, cli, param)

	go net.CheckExpiredParams() // 启动定时任务
	http.HandleFunc("/ask", net.AskHandler)
	http.HandleFunc("/getLatestOrder", net.GetLatestOrderHandler)
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		return
	}
}
