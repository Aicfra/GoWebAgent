package model

type ReqDTO struct {
	HttpMethod string  `json:"httpMethod"`
	HttpBody   []byte  `json:"httpBody"`
	MsgID      string  `json:"msgId"`
}