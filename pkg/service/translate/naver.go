package translate

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gaaon/quote-collector/pkg/constant"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type translateNaverResponse struct {
	Translated string `json:"translatedText"`
}

func FindTranslationByNaver(content string) (string, error) {

	transId := uuid.New().String()

	jsonRaw := "\"" + transId + "\"," +
		"\"dict\":false,\"honorific\":false,\"instant\":false," +
		"\"source\":\"en\",\"target\":\"ko\",\"text\":\"" + strings.ReplaceAll(content, "\"", "\\\"") + "\"}"

	encoded := "rlWxMKMcL2IWMPV6" + base64.StdEncoding.EncodeToString([]byte(jsonRaw))

	form := url.Values{}
	form.Add("data", encoded)

	req, _ := http.NewRequest(
		"POST",
		"https://papago.naver.com/apis/n2mt/translate",
		strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var nRes translateNaverResponse
	bodyRaw, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(bodyRaw, &nRes); err != nil {
		log.Fatal(err)
	}

	return nRes.Translated, nil
}

func FindTranslationByNaverAndSave(content string, entity model.QuoteEntity) (string, interface{}, error) {
	translated, err := FindTranslationByNaver(content)
	if err != nil {
		return "", nil, err
	}

	_id, err := repository.InsertTranslation(model.TranslationEntity{
		Content:   translated,
		Vendor:    constant.NAVER,
		QuoteId:   entity.Id,
		CreatorId: entity.CreatorId,
	})

	return translated, _id, err
}
