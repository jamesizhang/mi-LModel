package net

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"
	"io"
	"mi-LModel/client"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Param struct {
	ID       string `json:"id"`
	Question string `json:"question"`
}

type ParamVal struct {
	paramToModel *openai.ChatCompletionNewParams
	Created      time.Time
	Updated      time.Time
}

type ServerResponse struct {
	Response string `json:"response"`
}

var (
	params      = make(map[string]*ParamVal)
	mu          sync.Mutex
	checkTime   = 60 * time.Second
	ctx         = context.Background()
	cli         = client.Client
	latestID    string
	xiaoDisplay bool
)

func CheckExpiredParams() {
	ticker := time.NewTicker(checkTime)
	for range ticker.C {
		now := time.Now()
		mu.Lock()
		for id, param := range params {
			if now.Sub(param.Updated) > 30*time.Minute {
				delete(params, id)
			}
		}
		mu.Unlock()
	}
}

func AskHandler(w http.ResponseWriter, r *http.Request) {
	var req Param
	//var bytedata []byte
	//n, err := r.Body.Read(bytedata)
	//if err != nil && err != io.EOF {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//str := string(bytedata[0:n])
	//fmt.Printf("Body %v\n", r.Body)
	fmt.Printf("%v\n", r)
	//fmt.Printf("%v\n", n)
	//fmt.Printf("%v\n", str)
	//fmt.Println(r.PostFormValue("id"))
	//req.ID = r.PostFormValue("id")
	//req.Question = r.PostFormValue("question")
	err := json.NewDecoder(r.Body).Decode(&req)
	if req.ID == "" && req.Question == "" {
		http.Error(w, "参数有错误", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	mu.Lock()
	param, exists := params[req.ID]
	if !exists {
		param = &ParamVal{paramToModel: client.DidiChatCompletionNewParams(), Created: time.Now(), Updated: time.Now()}
		params[req.ID] = param
	} else {
		param.Updated = time.Now() // 更新创建时间
	}
	mu.Unlock()
	_, exists = params[latestID]
	if !exists {
		latestID = req.ID
	} else if params[latestID].Created.Before(params[req.ID].Created) {
		latestID = req.ID
	}
	if strings.Contains(req.Question, "天气") {
		weather := client.FetchAndDisplayWeather("北京昌平区", "9dd8b5a4f0425cab0d43c44212e464a7")
		req.Question = fmt.Sprintf("这是用户关于天气的询问：%s，这是查询国家气象局的结果：%s,跟进天气状况引导用户打车出行！", req.Question, weather)
	}

	//调用cli.Ask的实际逻辑
	answer := cli.Ask(ctx, req.Question, param.paramToModel)

	if checkSuccess(answer) {
		if !xiaoDisplay {
			xiaoDisplay = true
			err := startApp("com.nylo.android", ".MainActivity")
			if err != nil {
				//return
			} else {
				fmt.Printf("%s\n", "调用成功")
			}
		} else {
			err := stopApp("com.nylo.android")
			if err != nil {
				//return
			}
			err = startApp("com.nylo.android", "MainActivity")
			if err != nil {
				//return
			} else {
				fmt.Printf("%s\n", "调用成功")
			}
		}
	}

	resp := ServerResponse{Response: answer}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

// 检查answer字符串是否包含特定的成功标志
func checkSuccess(answer string) bool {
	return strings.Contains(answer, "预约成功") ||
		strings.Contains(answer, "叫车成功") ||
		strings.Contains(answer, "成功")
}

// 执行ADB命令来停止指定的应用
func stopApp(packageName string) error {
	cmd := exec.Command("adb", "shell", "am", "force-stop", packageName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行ADB命令失败: %v", err)
	}
	return nil
}

// 执行ADB命令来启动指定的应用
func startApp(packageName, className string) error {
	cmd := exec.Command("adb", "shell", "am", "start", "-n", packageName+"/"+className)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行ADB命令失败: %v", err)
	}
	return nil
}

func GetLatestOrderHandler(w http.ResponseWriter, r *http.Request) {
	print(r)
	if latestID == "" {
		http.Error(w, "该用户没有打车记录", http.StatusBadRequest)
		return
	}
	pararNeed2Model, exist := params[latestID]
	if !exist {
		http.Error(w, "打车记录已经过期", http.StatusBadRequest)
		return
	}
	//调用cli.GetLastInfo
	answer, err := cli.GetLastInfo(ctx, pararNeed2Model.paramToModel)
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(answer)
	if err != nil {
		return
	}
}
