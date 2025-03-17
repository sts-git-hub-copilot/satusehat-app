package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Anamnesa struct {
	AnamnesaId      int64  `json:"anamnesaId"`
	TenantId        int64  `json:"tenantId"`
	AnamnesaCode    string `json:"anamnesaCode" validate:"required"`
	AnamnesaName    string `json:"anamnesaName" validate:"required"`
	AnamnesaVersion string `json:"anamnesaVersion"`

	glentity.MasterEntity
}

func (d Anamnesa) TableName() string {
	return SS_AMNANESA
}

func (t Anamnesa) Columns() string {
	return `anamnesa_id, tenant_id, anamnesa_code, anamnesa_name, anamnesa_version,
        version, create_datetime, create_user_id, update_datetime, update_user_id,
        active, active_datetime, non_active_datetime`
}
