package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Trx struct {
	TrxId           int64  `json:"trxId"`
	TenantId        int64  `json:"tenantId"`
	PatientId       int64  `json:"patientId"`
	DatetimeIn      string `json:"datetimeIn"`
	DatetimeProcess string `json:"datetimeProcess"`
	DatetimeOut     string `json:"datetimeOut"`
	MinDuration     int64  `json:"minDuration"`
	StatusDoc       string `json:"statusDoc"`
	DataJson        string `json:"dataJson"`

	glentity.BaseEntity
}

func (d Trx) TableName() string {
	return SS_TRX
}

func (t Trx) Columns() string {
	return `trx_id, tenant_id, patient_id, datetime_in, datetime_process, datetime_out,
		min_duration, status_doc, data_json,
        version, create_datetime, create_user_id, update_datetime, update_user_id`
}
