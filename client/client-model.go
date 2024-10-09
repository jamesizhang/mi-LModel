package client

import (
	"context"
	"log"
	"mi-LModel/util"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ClientStuct struct {
	*openai.Client
}

var Client *ClientStuct

func init() {
	Client = &ClientStuct{
		Client: openai.NewClient(
			option.WithAPIKey(util.ApiKay),
			option.WithBaseURL(util.ApiBaseUrl)),
	}
}

func NewChatCompletionNewParams() *openai.ChatCompletionNewParams {
	return &openai.ChatCompletionNewParams{}
}

func DidiChatCompletionNewParams() *openai.ChatCompletionNewParams {
	params := NewChatCompletionNewParams()
	params.Messages = openai.F([]openai.ChatCompletionMessageParamUnion{
		openai.AssistantMessage("现在你是滴滴公司的电话接线员，负责通过电话帮客户进行预约打车需求。" +
			"比如客户说今天想打车到工区，你需要主动询问客户所在地点和打车的类型。其他不需要问。以下是其他的其他的场景问答，" +
			"当客户问你车打好了吗？你要明确告知车牌号和车子颜色。" +
			"当客户问司机还有多久才到？你要明确回复说车距离你还有2公里，预计等待2分钟。" +
			"对于其他问题，回答要简单并且利索。"),
	})
	params.Model = openai.F(util.ChatModelAli)
	params.Seed = openai.F(int64(8))
	return params
}

func (cli *ClientStuct) Ask(ctx context.Context, question string, param *openai.ChatCompletionNewParams) (resp string) {
	param.Messages.Value = append(param.Messages.Value, openai.UserMessage(question))
	response, err := cli.Chat.Completions.New(ctx, *param)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	resp = response.Choices[0].Message.Content
	log.Printf("%v", resp)
	param.Messages.Value = append(param.Messages.Value, openai.AssistantMessage(resp))
	return
}
