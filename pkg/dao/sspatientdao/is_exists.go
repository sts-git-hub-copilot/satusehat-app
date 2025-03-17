package sspatientdao

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputIsExists struct {
	Tx    gldb.Tx
	Audit gldata.AuditData
	Nik   string `json:"nik" validate:"required"`
}

func IsExists(input InputIsExists) (bool, *tables.Patient, error) {
	results := make([]*tables.Patient, 0)
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return false, nil, err
	}

	err := gldb.SelectQTx(input.Tx, *gldb.NewQBuilder().
		Add(` SELECT patient_id, tenant_id, nik, idss, patient_name, address, data_json,
              version, create_datetime, create_user_id, update_datetime, update_user_id,
              active, active_datetime, non_active_datetime `).
		Add(" FROM ", tables.SS_PATIENT).
		Add(" WHERE nik = :nik AND tenant_id = :tenantId ").
		SetParam("tenantId", input.Audit.TenantLoginId).
		SetParam("nik", input.Nik), &results)

	if err != nil {
		return false, nil, glerr.Wrap("error query is exists patient: ", err)
	}

	if len(results) == 0 {
		return false, nil, nil
	}

	return true, results[0], nil
}
