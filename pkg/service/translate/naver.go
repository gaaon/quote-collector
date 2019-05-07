package translate

import (
	"encoding/base64"
	"encoding/json"
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

func TranslateByNaver(content string) (string, error) {

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