package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type GeoResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		Country          string `json:"country"`
		Province         string `json:"province"`
		Citycode         string `json:"citycode"`
		City             string `json:"city"`
		District         string `json:"district"`
		Adcode           string `json:"adcode"`
		Location         string `json:"location"`
		Level            string `json:"level"`
	} `json:"geocodes"`
}

type WeatherResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Lives    []struct {
		Province      string `json:"province"`
		City          string `json:"city"`
		Adcode        string `json:"adcode"`
		Weather       string `json:"weather"`
		Temperature   string `json:"temperature"`       // 温度
		TemperatureF  string `json:"temperature_float"` // 浮点型温度
		Winddirection string `json:"winddirection"`     // 风向
		Windpower     string `json:"windpower"`         // 风力
		Humidity      string `json:"humidity"`          // 湿度
		HumidityF     string `json:"humidity_float"`    // 浮点型湿度
		Reporttime    string `json:"reporttime"`        // 报告时间
	} `json:"lives"`
}

func getAdcode(address string, apiKey string) (string, error) {
	queryParams := url.Values{}
	queryParams.Add("key", apiKey)
	queryParams.Add("address", address)

	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?%s", queryParams.Encode()))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geoResp GeoResponse
	err = json.Unmarshal(body, &geoResp)
	if err != nil {
		return "", err
	}

	if geoResp.Status == "1" && len(geoResp.Geocodes) > 0 {
		return geoResp.Geocodes[0].Adcode, nil
	}

	return "", fmt.Errorf("无法从地址获取ADCode")
}

func getWeather(adcode string, apiKey string) (WeatherResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("key", apiKey)
	queryParams.Add("city", adcode)

	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?%s", queryParams.Encode()))
	if err != nil {
		return WeatherResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return WeatherResponse{}, err
	}

	if weatherResp.Status == "1" && len(weatherResp.Lives) > 0 {
		return weatherResp, nil
	}

	return WeatherResponse{}, fmt.Errorf("无法获取天气信息")
}

func FetchAndDisplayWeather(address string, apiKey string) string {
	adcode, err := getAdcode(address, apiKey)
	if err != nil {
		log.Printf("获取ADCode失败: %v", err)
		return "获取ADCode失败"
	}

	fmt.Printf("ADCode for %s is: %s\n", address, adcode)

	weatherResp, err := getWeather(adcode, apiKey)
	if err != nil {
		log.Fatalf("获取天气信息失败: %v", err)
		return "获取天气信息失败"
	}

	// 提取关键字段
	if len(weatherResp.Lives) > 0 {
		live := weatherResp.Lives[0]
		weatherSummary := fmt.Sprintf("当前%s%s的天气状况为：%s，气温约为%v摄氏度，湿度为%v%%。",
			live.Province,
			live.City,
			live.Weather,
			live.TemperatureF,
			live.HumidityF)
		fmt.Println(weatherSummary)
		return weatherSummary
	}
	return "该城市没有天气信息"
}
