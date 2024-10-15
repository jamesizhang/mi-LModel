package net

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"
	"io"
	"mi-LModel/client"
	"net/http"
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
	params    = make(map[string]*ParamVal)
	mu        sync.Mutex
	checkTime = 60 * time.Second
	ctx       = context.Background()
	cli       = client.Client
	latestID  string
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

	//调用cli.Ask的实际逻辑
	answer := cli.Ask(ctx, req.Question, param.paramToModel)

	resp := ServerResponse{Response: answer}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
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
