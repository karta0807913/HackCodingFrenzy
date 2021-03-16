package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

const FrenzyHost string = "http://coding-frenzy.arping.me"

type FrenzyClient struct {
	session string
	client  http.Client
}

type FrenzyItem struct {
	Name string
	GID  string
}

func FrenzyLogin(account, password string) (*FrenzyClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Jar: jar,
	}
	data := url.Values{}
	data.Set("主令", "用戶登入")
	data.Set("瘋狂郵學", account)
	data.Set("瘋狂密碼", password)
	data.Set("桌類", "練桌")

	requestURL, err := url.Parse(FrenzyHost)

	if err != nil {
		return nil, err
	}

	requestURL.Path = "/v2-session.php"

	response, err := client.Post(requestURL.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	loginBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(loginBuffer, response.Body)
	if err != nil {
		return nil, err
	}
	loginBytes := loginBuffer.Bytes()
	re := regexp.MustCompile(`,錯誤:'(.*)'`)
	match := re.FindSubmatch(loginBytes)
	if len(match) == 0 {
		return nil, errors.New("登入失敗,請檢查帳號密碼")
	}
	if string(match[1]) != "" {
		return nil, errors.New(string(match[1]))
	}
	re = regexp.MustCompile(`{none:'none',用戶_學號:'(.*)'`)
	match = re.FindSubmatch(loginBytes)
	if len(match) == 0 || string(match[1]) == "訪客" {
		return nil, errors.New("登入失敗,請檢查帳號密碼")
	}

	for _, cookie := range client.Jar.Cookies(requestURL) {
		res := strings.Split(cookie.String(), "=")
		if res[0] != "PHPSESSID" {
			continue
		}
		return &FrenzyClient{
			client:  client,
			session: res[1],
		}, nil
	}
	return nil, errors.New("session not found")
}

func (frenzy *FrenzyClient) Load49() ([]FrenzyItem, error) {
	postURL, err := url.Parse(FrenzyHost)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("主令", "進入課程")
	data.Set("course_gid", "9c06c313a3c0ed313f6562c83691a067")
	data.Set("桌類", "練桌")

	postURL.Path = "/v2-session.php"

	_, err = frenzy.client.Post(postURL.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	data = url.Values{}
	data.Set("主令", "切換主頁")
	data.Set("主頁", "週事內容")
	data.Set("week_gid", "4a42359b3955f382ade06991a4df0198")
	data.Set("桌類", "練桌")

	_, err = frenzy.client.Post(postURL.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	postURL.Path = "/v2-index-student-items.php"

	data = url.Values{}
	data.Set("設定", "取頁")
	data.Set("桌類", "練桌")

	resp, err := frenzy.client.Post(postURL.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, resp.Body)

	items := make([]FrenzyItem, 0, 49)
	re := regexp.MustCompile(`ITEM_GID:'([^']*)'[^：]*：([^<]*) <\/nobr>`)
	for _, val := range re.FindAllStringSubmatch(buffer.String(), -1) {
		items = append(items, FrenzyItem{
			GID:  val[1],
			Name: val[2],
		})
	}
	return items, nil
}

func (frenzy *FrenzyClient) UpdateExamSave(gid, context string) error {
	postURL, err := url.Parse(FrenzyHost)
	if err != nil {
		return err
	}

	postURL.Path = "/frenzy-submit-record.php"

	encodedString := `&%e7%98%8b%e3%80%80%e7%8b%82%e3%80%80%e7%a8%8b%e3%80%80%e8%a8%ad%e3%80%80=` + frenzy.session + `&&item_gid=` + gid + `&%e5%ad%98%e6%a8%a1=%e8%a9%95%e9%87%8f&%e5%ad%98%e7%a2%bc=` + url.QueryEscape(context) + `&%e5%ad%98%e7%a7%92=0`
	resp, err := frenzy.client.Post(postURL.String(), "application/x-www-form-urlencoded", strings.NewReader(encodedString))
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(d) != "評量存檔完畢\r\n" {
		return errors.New(string(d))
	}

	return nil
}
