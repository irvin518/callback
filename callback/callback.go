package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/irvin518/callback/dohttp"
)

type CallBackToWowfish struct {
}

type CallbackBaseData struct {
	From string `json:"from"`
	To   string `json:"to"`
	Ret  int64  `json:"ret"`
	Info string `json:"info"`
	Sign string `json:"sign"`
}

type CallBackToWowfishData struct {
	CallbackBaseData
	Amount string `json:"amount"`
}

type CallBackToTonData struct {
	CallbackBaseData
	Amount  string `json:"amount"`
	Commont string `json:"payload"`
}

type CallBackToWowfishNftData struct {
	CallbackBaseData
	Id string `json:"nft_id"`
}

type CallBackUserWithDrawData struct {
	Payload string `json:"payload"`
	Hash    string `json:"hash"`
}

var ins = CallBackToWowfish{}

func Instance() *CallBackToWowfish {
	return &ins
}

func (c *CallBackToWowfish) Callback(url string, data any) error {
	if url != "" {
		sendData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		rsp, err := dohttp.DoJsonHttp(map[string]string{}, "POST", url, sendData)
		if nil != err {
			return err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("callback error %d", rsp.StatusCode)
		}

		data, err := io.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		type InnerCallbackRet struct {
			Ret  int    `json:"ret"`  //回调返回状态，ret=0 成功
			Info string `json:"info"` //如果有错误，返回错误信息
		}
		resp := &InnerCallbackRet{}
		err = json.Unmarshal(data, resp)
		if err != nil {
			return err
		}
		if resp.Ret != 0 {
			return fmt.Errorf(`callback error:%d, info:%s`, resp.Ret, resp.Info)
		}
		return nil
	}

	return errors.New("callback url is null")
}
