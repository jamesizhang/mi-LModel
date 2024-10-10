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
)

func CheckExpiredParams() {
	ticker := time.NewTicker(checkTime)
	for range ticker.C {
		now := time.Now()
		mu.Lock()
		for id, param := range params {
			if now.Sub(param.Created) > 5*time.Minute {
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
		param = &ParamVal{paramToModel: client.DidiChatCompletionNewParams(), Created: time.Now()}
		params[req.ID] = param
	} else {
		param.Created = time.Now() // 更新创建时间
	}
	mu.Unlock()

	//调用cli.Ask的实际逻辑
	answer := cli.Ask(ctx, req.Question, param.paramToModel)

	resp := ServerResponse{Response: answer}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
