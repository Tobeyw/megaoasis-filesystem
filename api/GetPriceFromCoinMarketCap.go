package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"metaoasis-filesystem/consts"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetPriceFromCoinMarketMap(covert string, slug string) (interface{}, interface{}, error) {
	idString := GetIdBySlug(slug)
	if idString == "" {
		return nil, nil, fmt.Errorf("slug invaild")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", consts.CoinMartketCapUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	q := url.Values{}
	q.Add("convert", covert)
	q.Add("slug", slug)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", consts.CoinMarketCapKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var price, updateTime interface{}
	var data map[string]interface{}
	if resp.Status == "200 OK" {
		if err := json.Unmarshal(respBody, &data); err == nil {
			data1 := data["data"].(map[string]interface{})
			id := data1[idString].(map[string]interface{})
			quote := id["quote"].(map[string]interface{})
			currency := quote[covert].(map[string]interface{})
			price = currency["price"]
			lastUpdateTime := currency["last_updated"].(string)
			lastUpdateTime = strings.Replace(lastUpdateTime, "T", " ", 1)
			lastUpdateTime = strings.Replace(lastUpdateTime, "Z", "", 1)

			setDate := "2024-02-05 02:59:00.000" //2024-02-05T02:59:00.000Z
			dateFormate := "2006-01-02 15:04:05"
			loc, _ := time.LoadLocation("UTC")                              //重要：获取时区
			timeObj, err := time.ParseInLocation(dateFormate, setDate, loc) //指定日期 转 当地 日期对象 类型为 time.Time

			if err != nil {
				fmt.Println("parse time failed err :", err)
				return nil, nil, fmt.Errorf("parse time failed err:%s", err)
			}
			updateTime = timeObj.UnixMilli()

		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, fmt.Errorf("failed to get price from CoinMarketCap")
	}
	return price, updateTime, nil
}

func GetIdBySlug(slug string) string {
	if slug == "neo" {
		return "1376"
	}
	if slug == "gas" {
		return "1785"
	}
	return ""
}
