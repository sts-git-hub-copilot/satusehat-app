package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type TrxDiagnosa struct {
	TrxDiagnosaId int64 `json:"trxDiagnosaId"`
	TenantId      int64 `json:"tenantId"`
	TrxId         int64 `json:"trxId"`
	DiagnosaId    int64 `json:"diagnosaId"`

	glentity.BaseEntity
}

func (d TrxDiagnosa) TableName() string {
	return SS_TRX_DIAGNOSA
}

func (t TrxDiagnosa) Columns() string {
	return `trx_diagnosa_id, tenant_id, trx_id, diagnosa_id,
        version, create_datetime, create_user_id, update_datetime, update_user_id`
}
