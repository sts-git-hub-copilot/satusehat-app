package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Doctor struct {
	DoctorId   int64  `json:"doctorId"`
	TenantId   int64  `json:"tenantId"`
	Nik        string `json:"nik"`
	Idss       string `json:"idss"`
	DoctorName string `json:"doctorName"`
	DataJson   string `json:"dataJson"`

	glentity.MasterEntity
}

func (d Doctor) TableName() string {
	return SS_DOCTOR
}

func (t Doctor) Columns() string {
	return `doctor_id, tenant_id, nik, idss, doctor_name, data_json, 
		version, create_datetime, create_user_id, update_datetime, update_user_id, 
		active, active_datetime, non_active_datetime`
}
