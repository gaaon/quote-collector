package translate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type translateKakaoResult struct {
	Translated string `json:"translated"`
}

type translateKakaoResponse struct {
	Result translateKakaoResult `json:"result"`
}

func TranslateByKakao(content string) (string, error) {
	form := url.Values{}
	form.Add("lang", "enkr")
	form.Add("q", content)

	req, _ := http.NewRequest(
		"POST",
		"https://translate.kakao.com/translator/translate.json",
		strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Referer", " https://translate.kakao.com/")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var tRes translateKakaoResponse
	bodyRaw, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(bodyRaw, &tRes); err != nil {
		log.Fatal(err)
	}

	return tRes.Result.Translated, nil
}