package model

type GeneralResponse struct {
	Code    int         `json:"status_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
