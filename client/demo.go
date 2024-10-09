package client

import (
	"context"
	"github.com/openai/openai-go"
)

func Demo(ctx context.Context, cli *ClientStuct, param *openai.ChatCompletionNewParams) {
	cli.Ask(ctx, "帮我打车到天空之城", param)
	cli.Ask(ctx, "我的位置在回龙观体育中心，想打快车类型", param)
	cli.Ask(ctx, "请问车打好了吗", param)
	cli.Ask(ctx, "司机还有多久才到呢", param)
}
