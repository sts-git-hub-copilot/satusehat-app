package ssresponselogdao

import (
	"context"
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputAdd struct {
	Tx    gldb.Tx
	Audit gldata.AuditData

	Url              string `json:"url" validate:"required"`
	Method           string `json:"method" validate:"required"`
	BodyReq          string `json:"bodyReq" default:"{}"`
	HeaderReq        string `json:"headerReq" default:"{}"`
	Status           string `json:"status" validate:"required"`
	RequestResponse  string `json:"requestResponse" default:"{}"`
	ResponseCode     string `json:"responseCode" validate:"required"`
	ResponseBody     string `json:"responseBody" default:"{}"`
	LatencyMs        int64  `json:"latencyMs"`
	ResponseDatetime string `json:"responseDatetime" validate:"required"`
}

func Add(input InputAdd) (*tables.ResponseLog, error) {
	result := tables.ResponseLog{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_RESPONSE_LOG).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("url", input.Url).
		AddValue("method", input.Method).
		AddValue("body_req", input.BodyReq).
		AddValue("header_req", input.HeaderReq).
		AddValue("status", input.Status).
		AddValue("request_response", input.RequestResponse).
		AddValue("response_code", input.ResponseCode).
		AddValue("response_body", input.ResponseBody).
		AddValue("latency_ms", input.LatencyMs).
		AddValue("response_datetime", input.ResponseDatetime).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` RETURNING response_log_id, tenant_id, url, method, body_req, header_req, status, 
				request_response, response_code, response_body, latency_ms,
				response_datetime, create_datetime, update_datetime, create_user_id, 
				update_user_id, version `), &result)

	if err != nil {
		return nil, errors.New("failed insert ss response log : " + err.Error())
	}

	return &result, nil
}

type InputAddByMt struct {
	Mt    context.Context
	Audit gldata.AuditData

	Url              string `json:"url" validate:"required"`
	Method           string `json:"method" validate:"required"`
	BodyReq          string `json:"bodyReq" default:"{}"`
	HeaderReq        string `json:"headerReq" default:"{}"`
	Status           string `json:"status" validate:"required"`
	RequestResponse  string `json:"requestResponse" default:"{}"`
	ResponseCode     string `json:"responseCode" validate:"required"`
	ResponseBody     string `json:"responseBody" default:"{}"`
	LatencyMs        int64  `json:"latencyMs"`
	ResponseDatetime string `json:"responseDatetime" validate:"required"`
}

func AddByMt(input InputAddByMt) (*tables.ResponseLog, error) {
	result := tables.ResponseLog{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQMt(input.Mt, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_RESPONSE_LOG).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("url", input.Url).
		AddValue("method", input.Method).
		AddValue("body_req", input.BodyReq).
		AddValue("header_req", input.HeaderReq).
		AddValue("status", input.Status).
		AddValue("request_response", input.RequestResponse).
		AddValue("response_code", input.ResponseCode).
		AddValue("response_body", input.ResponseBody).
		AddValue("latency_ms", input.LatencyMs).
		AddValue("response_datetime", input.ResponseDatetime).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` RETURNING response_log_id, tenant_id, url, method, body_req, header_req, status, 
				request_response, response_code, response_body, latency_ms,
				response_datetime, create_datetime, update_datetime, create_user_id, 
				update_user_id, version `), &result)

	if err != nil {
		return nil, errors.New("failed insert ss response log : " + err.Error())
	}

	return &result, nil
}
