package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Diagnosa struct {
	DiagnosaId      int64  `json:"diagnosaId"`
	TenantId        int64  `json:"tenantId"`
	DiagnosaCode    string `json:"diagnosaCode" validate:"required"`
	DiagnosaName    string `json:"diagnosaName" validate:"required"`
	DiagnosaVersion string `json:"diagnosaVersion"`

	glentity.MasterEntity
}

func (d Diagnosa) TableName() string {
	return SS_DIAGNOSA
}

func (t Diagnosa) Columns() string {
	return `diagnosa_id, tenant_id, diagnosa_code, diagnosa_name, diagnosa_version,
        version, create_datetime, create_user_id, update_datetime, update_user_id,
        active, active_datetime, non_active_datetime`
}
