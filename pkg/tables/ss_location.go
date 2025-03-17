package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Location struct {
	LocationId   int64  `json:"locationId"`
	TenantId     int64  `json:"tenantId"`
	Idss         string `json:"idss"`
	LocationCode string `json:"locationCode"`
	LocationName string `json:"locationName"`
	Address      string `json:"address"`
	DataJson     string `json:"dataJson"`

	glentity.MasterEntity
}

func (d Location) TableName() string {
	return SS_LOCATION
}

func (t Location) Columns() string {
	return `location_id, tenant_id, idss, location_name, location_code, address, data_json, 
		version, create_datetime, create_user_id, update_datetime, update_user_id, 
		active, active_datetime, non_active_datetime`
}
