package translate

import (
	"encoding/json"
	"github.com/gaaon/quote-collector/pkg/constant"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
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

func FindTranslationByKakao(content string) (string, error) {
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

func FindTranslationByKakaoAndSave(content string, entity model.QuoteEntity) (string, interface{}, error) {
	translated, err := FindTranslationByKakao(content)
	if err != nil {
		return "", nil, err
	}

	_id, err := repository.InsertTranslation(model.TranslationEntity{
		Content:   translated,
		Vendor:    constant.KAKAO,
		QuoteId:   entity.Id,
		CreatorId: entity.CreatorId,
	})

	return translated, _id, err
}