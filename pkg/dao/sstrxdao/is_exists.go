package sstrxdao

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputIsExists struct {
	Tx         gldb.Tx
	Audit      gldata.AuditData
	PatientId  int64  `json:"patientId"`
	DatetimeIn string `json:"datetimeIn"`
}

func IsExists(input InputIsExists) (bool, error) {
	var count int64
	err := gldb.SelectRowQTx(input.Tx, *gldb.NewQBuilder().
		Add(" SELECT COUNT(1) ").
		Add(" FROM ", tables.SS_TRX).
		Add(" WHERE patient_id = :patientId ").
		Add(" AND tenant_id = :tenantId ").
		Add(" AND datetime_in = :datetimeIn ").
		SetParam("tenantId", input.Audit.TenantLoginId).
		SetParam("patientId", input.PatientId).
		SetParam("datetimeIn", input.DatetimeIn), &count)

	if err != nil {
		return false, glerr.New("failed to check trx", err.Error())
	}

	return count > 0, nil
}
