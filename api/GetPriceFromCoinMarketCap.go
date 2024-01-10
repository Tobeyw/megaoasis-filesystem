package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"metaoasis-filesystem/consts"
	"net/http"
	"net/url"
)

func GetPriceFromCoinMarketMap(covert string, slug string) (interface{}, error) {
	idString := GetIdBySlug(slug)
	if idString == "" {
		return nil, fmt.Errorf("slug invaild")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", consts.CoinMartketCapUrl, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("convert", covert)
	q.Add("slug", slug)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", consts.CoinMarketCapKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var price interface{}
	var data map[string]interface{}
	if resp.Status == "200 OK" {
		if err := json.Unmarshal(respBody, &data); err == nil {
			data1 := data["data"].(map[string]interface{})
			id := data1[idString].(map[string]interface{})
			quote := id["quote"].(map[string]interface{})
			currency := quote[covert].(map[string]interface{})
			price = currency["price"]

		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("failed to get price from CoinMarketCap")
	}
	return price, nil
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
