package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type TrxIdssen struct {
	IdssenId int64  `json:"idssenId"`
	TenantId int64  `json:"tenantId"`
	TrxId    int64  `json:"trxId"`
	Idssen   string `json:"idssen"`
	Status   string `json:"status"`

	glentity.BaseEntity
}

func (d TrxIdssen) TableName() string {
	return SS_TRX_IDSSEN
}

func (t TrxIdssen) Columns() string {
	return `idssen_id, tenant_id, trx_id, idssen, status,
        version, create_datetime, create_user_id, update_datetime, update_user_id`
}
