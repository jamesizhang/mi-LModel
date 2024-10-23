package client

import "testing"

func TestFetchAndDisplayWeather(t *testing.T) {
	apiKey := "9dd8b5a4f0425cab0d43c44212e464a7" // 使用实际的API密钥
	address := "北京昌平区"

	FetchAndDisplayWeather(address, apiKey)
}
