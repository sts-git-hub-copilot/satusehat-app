package tables

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Patient struct {
	PatientId   int64  `json:"patientId"`
	TenantId    int64  `json:"tenantId"`
	Nik         string `json:"nik"`
	Idss        string `json:"idss"`
	PatientName string `json:"patientName"`
	Address     string `json:"address"`
	DataJson    string `json:"dataJson"`

	glentity.MasterEntity
}

func (d Patient) TableName() string {
	return SS_PATIENT
}

func (t Patient) Columns() string {
	return `patient_id, tenant_id, nik, idss, patient_name, address, data_json, 
		version, create_datetime, create_user_id, update_datetime, update_user_id, 
		active, active_datetime, non_active_datetime`
}
