package main

import (
	"context"
	"mi-LModel/client"
)

func main() {
	ctx := context.Background()
	cli := client.Client

	param := client.DidiChatCompletionNewParams()

	client.Demo(ctx, cli, param)
}
