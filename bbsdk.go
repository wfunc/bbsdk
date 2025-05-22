package bbsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wfunc/util/xhash"
	"github.com/wfunc/util/xmap"
)

const (
	KeyXBBID    = "X-BB-Id"
	KeyXBBTOKEN = "X-BB-Token"

	CreateUserURL        = "/api/v2/createUser"
	LoginURL             = "/api/v2/login"
	LogoutURL            = "/api/v2/logout"
	TransferURL          = "/api/v2/transfer"
	GetGameRecordURL     = "/api/v2/getGameRecord"
	GetBalanceURL        = "/api/v2/getBalance"
	GetUserControlURL    = "/api/v2/getUserControl"
	SetUserControlURL    = "/api/v2/setUserControl"
	CancelUserControlURL = "/api/v2/cancelUserControl"
	SetUserLimitURL      = "/api/v2/setUserLimit"
	password             = "123"
)

var (
	URL       string
	BBID      string
	BBTOKEN   string
	Uppername string
	Client    *http.Client
)

func Bootstrap(url, bbid, bbtoken, uppername string) {
	URL = url
	BBID = bbid
	BBTOKEN = bbtoken
	Uppername = uppername
	Client = &http.Client{
		Timeout: 5 * time.Second,
	}
}

//	type CreateUserParams struct {
//		Username string `form:"Username" binding:"required"`
//		Password string `form:"Password" binding:"required"`
//	}
func CreateUser(username string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s&Password=%s", URL, CreateUserURL, username, password)
	return GET(url)
}

func Transfer(username, remitNo, action string, remit string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s&RemitNo=%s&Action=%s&Remit=%s", URL, TransferURL, username, remitNo, action, remit)
	return GET(url)
}

func GetGameRecord(action, date, starttime, endtime, gametype, page, pagelimit string) (xmap.M, error) {
	params := fmt.Sprintf("action=%s&date=%s&starttime=%s&endtime=%s", action, date, starttime, endtime)
	// &gametype=%s&page=%s&pagelimit=%s
	if len(gametype) > 0 {
		params += fmt.Sprintf("&gametype=%s", gametype)
	}
	if len(page) > 0 {
		params += fmt.Sprintf("&page=%s", page)
	}
	if len(pagelimit) > 0 {
		params += fmt.Sprintf("&pagelimit=%s", pagelimit)
	}
	url := fmt.Sprintf("%s%s?%s", URL, GetGameRecordURL, params)
	return GET(url)
}

func GetBalance(username string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s", URL, GetBalanceURL, username)
	return GET(url)
}

func GetUserControl(username string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s", URL, GetUserControlURL, username)
	return GET(url)
}

func SetUserControl(username, winRate, loseRate, tieRate, balance string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s&WinRate=%s&LoseRate=%s&TieRate=%s&Balance=%s", URL, SetUserControlURL, username, winRate, loseRate, tieRate, balance)
	return GET(url)
}

func CancelUserControl(username string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s", URL, CancelUserControlURL, username)
	return GET(url)
}

func SetUserLimit(username, btn, btx string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s&Btn=%s&Btx=%s", URL, SetUserLimitURL, username, btn, btx)
	return GET(url)
}

func Logout(username string) (xmap.M, error) {
	url := fmt.Sprintf("%s%s?Username=%s", URL, LogoutURL, username)
	return GET(url)
}

func AdjustDateByTimezone() string {
	// 目标时区：GMT+4
	loc, _ := time.LoadLocation("Etc/GMT+4")

	// 当前时间
	currentTime := time.Now()

	// 转换为 GMT+4 时区时间
	adjustedTime := currentTime.In(loc)

	// 返回 GMT+4 时间格式化输出
	return adjustedTime.Format("2006-01-02")
}

func Login(username string) string {
	date := AdjustDateByTimezone()
	before := date + "Uppername:" + Uppername + "Username:" + username + BBTOKEN
	key := xhash.SHA256([]byte(before))
	url := fmt.Sprintf("%s%s?Username=%s&Uppername=%s&Key=%s", URL, LoginURL, username, Uppername, key)
	return url
}

func GET(url string) (xmap.M, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(KeyXBBID, BBID)
	req.Header.Set(KeyXBBTOKEN, BBTOKEN)

	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	var result xmap.M
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}
