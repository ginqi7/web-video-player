package main

import (
	"encoding/json"
	"fmt"
)

type Dictionary struct {
	Usphone string `json:"usphone"`
	Ukphone string `json:"ukphone"`
	Trs     []string `json:"trs"`
}

func query(word string) *Dictionary {

	paramsMap := createRequestParams(word)
	header := map[string][]string{}
	result := DoGet("https://dict.youdao.com/jsonapi", header, paramsMap, "application/json")

	if result != nil {
		var jsonObj map[string]interface{}

		// 解析 JSON 字符串
		err := json.Unmarshal(result, &jsonObj)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		words := jsonObj["ec"].(map[string]interface{})["word"].([]interface{})
		dictionary := Dictionary{}
		for _, word := range words {
			wordDict := word.(map[string]interface {})
			dictionary.Usphone = wordDict["usphone"].(string)
			dictionary.Ukphone = wordDict["ukphone"].(string)
			trs := wordDict["trs"].([]interface {})
			for _, transtate := range trs {
				translateStr := transtate.(map[string]interface{})["tr"].([]interface{})[0].(map[string]interface{})["l"].(map[string]interface{})["i"].([]interface{})[0].(string)
				fmt.Println(translateStr)
				dictionary.Trs = append(dictionary.Trs, translateStr)
			}
		}
		return &dictionary

	}
	return nil
}

func createRequestParams(word string) map[string][]string {
	q := word
	langType := "auto"
	jsonversion := "2"
	dicts := "%7B%22count%22%3A99%2C%22dicts%22%3A%5B%5B%22ec%22%2C%22ce%22%2C%22newcj%22%2C%22newjc%22%2C%22kc%22%2C%22ck%22%2C%22fc%22%2C%22cf%22%2C%22multle%22%2C%22jtj%22%2C%22pic_dict%22%2C%22tc%22%2C%22ct%22%2C%22typos%22%2C%22special%22%2C%22tcb%22%2C%22baike%22%2C%22lang%22%2C%22simple%22%2C%22wordform%22%2C%22exam_dict%22%2C%22ctc%22%2C%22web_search%22%2C%22auth_sents_part%22%2C%22ec21%22%2C%22phrs%22%2C%22input%22%2C%22wikipedia_digest%22%2C%22ee%22%2C%22collins%22%2C%22ugc%22%2C%22media_sents_part%22%2C%22syno%22%2C%22rel_word%22%2C%22longman%22%2C%22ce_new%22%2C%22le%22%2C%22newcj_sents%22%2C%22blng_sents_part%22%2C%22hh%22%5D%2C%5B%22ugc%22%5D%2C%5B%22longman%22%5D%2C%5B%22newjc%22%5D%2C%5B%22newcj%22%5D%2C%5B%22web_trans%22%5D%2C%5B%22fanyi%22%5D%5D%7D&keyfrom=mdict.7.2.0.android"
	client := "mobile"
	model := "honor"
	mid := "5.6.1"
	imei := "659135764921685"
	vendor := "wandoujia"
	screen := "1080x1800"
	ssid := "superman"
	network := "wifi"
	abtest := "2"
	xmlVersion := "5.1"

	return map[string][]string{
		"jsonversion": {jsonversion},
		"client":      {client},
		"q":           {q},
		"langType":    {langType},
		"dicts":       {dicts},
		"model":       {model},
		"mid":         {mid},
		"imei":        {imei},
		"vendor":      {vendor},
		"screen":      {screen},
		"ssid":        {ssid},
		"network":     {network},
		"abtest":      {abtest},
		"xmlVersion":  {xmlVersion},
	}
}
