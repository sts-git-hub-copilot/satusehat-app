package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type ResponseLog struct {
	ResponseLogId    int64  `json:"responseLogId"`
	TenantId         int64  `json:"tenantId"`
	Url              string `json:"url"`
	Method           string `json:"method"`
	BodyReq          string `json:"bodyReq"`
	HeaderReq        string `json:"headerReq"`
	Status           string `json:"status"`
	RequestResponse  string `json:"requestResponse"`
	ResponseCode     string `json:"responseCode"`
	ResponseBody     string `json:"responseBody"`
	LatencyMs        int64  `json:"latencyMs"`
	ResponseDatetime string `json:"responseDatetime"`

	glentity.BaseEntity
}

func (d ResponseLog) TableName() string {
	return SS_RESPONSE_LOG
}
