package protocol

import (
	"fmt"

	"lightspeed.2dao3.com/nokserver/consts/errcode"
)

type BaseResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func NewBaseResponse(code int, err error) (resp BaseResponse) {
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("%s:%s", errcode.ErrMsg[code], err)
	} else {
		errMsg = errcode.ErrMsg[code]
	}
	return BaseResponse{
		ErrCode: code,
		ErrMsg:  errMsg,
	}
}

type OnTaskCompleteRequest struct{}
