package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type TrxItem struct {
	TrxItemId  int64  `json:"trxItemId"`
	TenantId   int64  `json:"tenantId"`
	TrxId      int64  `json:"trxId"`
	LineNo     int64  `json:"lineNo"`
	DoctorId   int64  `json:"doctorId"`
	LocationId int64  `json:"locationId"`
	DataJson   string `json:"dataJson"`

	glentity.BaseEntity
}

func (d TrxItem) TableName() string {
	return SS_TRX_ITEM
}

func (t TrxItem) Columns() string {
	return `trx_item_id, tenant_id, trx_id, line_no, doctor_id, location_id, data_json,
        version, create_datetime, create_user_id, update_datetime, update_user_id`
}
