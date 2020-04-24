package translate

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github/wbellmelodyw/gin-wechat/logger"
	"golang.org/x/text/language"
	"net/url"
	"strings"
)

type GoogleTranslator struct {
	form language.Tag //来源
	to   language.Tag //要翻译成
}

func GetGoogle(form, to language.Tag) *GoogleTranslator {
	return &GoogleTranslator{
		form,
		to,
	}
}

func (g *GoogleTranslator) Text(text string) (string, error) {
	token := GetToken(text)

	urll := "https://translate.google.com/translate_a/single"
	data := map[string]string{
		"client": "gtx",
		"sl":     g.form.String(),
		"tl":     g.to.String(),
		"hl":     g.to.String(),
		// "dt":     []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		"ie":   "UTF-8",
		"oe":   "UTF-8",
		"otf":  "1",
		"ssel": "0",
		"tsel": "0",
		"kc":   "7",
		"q":    text,
		"tk":   token,
	}
	client := resty.New()

	r, err := client.R().SetQueryParams(data).
		SetQueryParamsFromValues(url.Values{
			"dt": []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		}).Get(urll)
	if err != nil {
		return "", err
	}
	//提取翻译
	texts := make([]string, 5)
	rspJson := r.String()
	//词意
	//result := gjson.Get(rspJson, "0.0")
	wordMean := "词意:" + gjson.Get(rspJson, "0.0.0").String()
	texts = append(texts, wordMean)
	//词性
	result := gjson.Get(rspJson, "1")
	wordAtr := "词性:" //少的用+就行 多才用 strings.builder
	for _, attrs := range result.Array() {
		wordAtr += attrs.Get("0").String() + ":"
		for _, attr := range attrs.Get("1").Array() {
			wordAtr += attr.String() + ","
		}
		wordAtr += ";"
		logger.Module("test").Sugar().Info("word2", wordAtr)
	}
	texts = append(texts, wordAtr)
	//解释
	result = gjson.Get(rspJson, "12")
	wordExplain := "解释:" //少的用+就行 多才用 strings.builder
	for _, attrs := range result.Array() {
		wordExplain += attrs.Get("0").String() + ":"
		for _, attr := range attrs.Get("1").Array() {
			wordExplain += attr.Get("0").String() + "|"
		}
		wordExplain = strings.TrimRight(wordExplain, "|")
		wordExplain += ";"
	}
	logger.Module("test").Sugar().Info("word3", wordExplain)
	texts = append(texts, wordExplain)
	//造句

	//废弃
	var resp []interface{}
	err = json.Unmarshal(r.Body(), &resp)
	if err != nil {
		return "", err
	}
	responseText := ""
	logger.Module("test").Sugar().Info("resp", resp)

	//翻译
	for _, obj := range resp[0].([]interface{}) {
		if len(obj.([]interface{})) == 0 {
			break
		}

		t, ok := obj.([]interface{})[0].(string)
		if ok {
			responseText += t
		}
	}
	//词性

	return responseText, nil
}

//func Audio(text string, from language.Tag, to language.Tag) (string, error) {
//
//}
