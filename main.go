package main

import (
	"context"
	"log"

	"mi-LModel/client"
	"mi-LModel/util"

	"github.com/openai/openai-go"
)

func main() {
	ctx := context.Background()
	cli := client.NewClient(util.ApiKay, util.ApiBaseUrl)

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage("现在你是滴滴公司的电话接线员，负责通过电话帮客户进行预约打车需求。" +
				"比如客户说今天想打车到工区，你需要主动询问客户所在地点和打车的类型。其他不需要问。以下是其他的其他的场景问答，" +
				"当客户问你车打好了吗？你要明确告知车牌号和车子颜色。" +
				"当客户问司机还有多久才到？你要明确回复说车距离你还有2公里，预计等待2分钟。" +
				"对于其他问题，回答要简单并且利索。"),
			openai.UserMessage("帮我打车到天空之城"),
		}),
		Model: openai.F(util.ChatModelAli),
		Seed:  openai.F(int64(8)),
	}

	response, err := cli.Chat.Completions.New(ctx, param)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	log.Printf("%v", response.Choices[0].Message.Content)
	param.Messages.Value = append(param.Messages.Value, openai.AssistantMessage(response.Choices[0].Message.Content))
	param.Messages.Value = append(param.Messages.Value, openai.UserMessage("我的位置在回龙观体育中心，想打快车类型"))
	response, err = cli.Chat.Completions.New(ctx, param)

	log.Printf("%v", response.Choices[0].Message.Content)
	param.Messages.Value = append(param.Messages.Value, openai.AssistantMessage(response.Choices[0].Message.Content))
	param.Messages.Value = append(param.Messages.Value, openai.UserMessage("请问车打好了吗"))
	response, err = cli.Chat.Completions.New(ctx, param)

	log.Printf("%v", response.Choices[0].Message.Content)
	param.Messages.Value = append(param.Messages.Value, openai.AssistantMessage(response.Choices[0].Message.Content))
	param.Messages.Value = append(param.Messages.Value, openai.UserMessage("司机还有多久才到呢"))
	response, err = cli.Chat.Completions.New(ctx, param)
	log.Printf("%v", response.Choices[0].Message.Content)
}
