package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type TrxAnamnesa struct {
	TrxAnamnesaId int64 `json:"trxAnamnesaId"`
	TenantId      int64 `json:"tenantId"`
	TrxId         int64 `json:"trxId"`
	AnamnesaId    int64 `json:"anamnesaId"`

	glentity.BaseEntity
}

func (d TrxAnamnesa) TableName() string {
	return SS_TRX_AMNANESA
}

func (t TrxAnamnesa) Columns() string {
	return `trx_anamnesa_id, tenant_id, trx_id, anamnesa_id,
        version, create_datetime, create_user_id, update_datetime, update_user_id`
}
