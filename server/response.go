package server

import "strings"

// **

type Resp struct {
	Status  RespCode    `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RespCode string

type RespError struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (e Resp) Normalize() Resp {
	// send empty object instead of nil
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}

	// change case message to lower
	e.Message = strings.ToLower(e.Message)

	return e
}

// **

type RespKind int

const (
	ServerRespSuccess RespKind = iota
	ServerRespPending
	ServerRespErrUnauthorized
	ServerRespErrPathNotFound
	ServerRespErrUnexpected
	ServerRespErrValidation
	ServerRespErrProcessFailed
)

var ListResp map[RespKind]Resp = map[RespKind]Resp{
	ServerRespSuccess:          {Status: "00", Message: "success"},
	ServerRespPending:          {Status: "01", Message: "process pending"},
	ServerRespErrUnauthorized:  {Status: "02", Message: "not authorized for action"},
	ServerRespErrPathNotFound:  {Status: "03", Message: "path not found"},
	ServerRespErrUnexpected:    {Status: "05", Message: "unexpected error"},
	ServerRespErrValidation:    {Status: "06", Message: "validation error"},
	ServerRespErrProcessFailed: {Status: "07", Message: "process failed"},
}
