package client

import (
	"context"
	"errors"
	"log"
	"mi-LModel/util"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ClientStuct struct {
	*openai.Client
}

type GetLastInfoResp struct {
	OrderTime        string
	StartLocation    string
	EndLocation      string
	carLicenseNumber string
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
			"比如客户说今天想打车到工区，你需要主动询问客户需要预约几点的车还是现在用车以及所在地点和打车的类型。其他不需要问。以下是其他的其他的场景问答，" +
			"当为预约单时候，预约成功后对客户说 相关预约信息也可以在卡片中进行查看" +
			"当客户问你车打好了吗？你要明确告知车牌号和车子颜色。" +
			"当客户问司机还有多久才到？你要明确回复说车距离你还有2公里，预计等待2分钟。" +
			"客户的其他问题，都假设客户在北京，并且一定要引导客户进行打车这一话题！一定要引导客户进行打车这一话题！此外回答要简单并且利索。"),
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

func (cli *ClientStuct) GetLastInfo(ctx context.Context, param *openai.ChatCompletionNewParams) (GetLastInfoResp, error) {
	question := "请按照以下格式给出用车具体时间，打车起点，打车终点，车牌号。以-进行分隔" +
		"例子：2024年10月15日下午2点10分-回龙观地铁站-万科天空之城-京A123456。其他并不需要多说！其他并不需要多说！因为我要利用-进行分割"
	response := cli.Ask(ctx, question, param)
	// 按照'-'分割字符串
	parts := strings.Split(response, "-")
	if len(parts) != 4 {
		return GetLastInfoResp{}, errors.New("大模型响应格式错误，应包含四个由'-'分隔的部分")
	}
	// 创建并填充GetLastInfoResp结构体
	infoResp := GetLastInfoResp{
		OrderTime:        parts[0],
		StartLocation:    parts[1],
		EndLocation:      parts[2],
		carLicenseNumber: parts[3],
	}
	return infoResp, nil
}
