package utils

type SuccessResp struct {
	Code int `json:"code" example:"0"`
	Msg string `json:"msg" example:""`
	Data interface{} `json:"data"`
}

type FailureResp struct {
	Code int `json:"code" example:"1"`
	Msg string `json:"msg" example:"错误信息对象"`
	Data interface{} `json:"data" example:"nil"`
}
