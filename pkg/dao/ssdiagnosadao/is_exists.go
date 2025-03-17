package ssdiagnosadao

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputIsExists struct {
	Tx           gldb.Tx
	Audit        gldata.AuditData
	DiagnosaCode string `json:"diagnosaCode" validate:"required"`
}

func IsExists(input InputIsExists) (bool, *tables.Diagnosa, error) {
	results := make([]*tables.Diagnosa, 0)
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return false, nil, err
	}

	err := gldb.SelectQTx(input.Tx, *gldb.NewQBuilder().
		Add(` SELECT diagnosa_id, tenant_id, diagnosa_code, diagnosa_name, diagnosa_version,
              version, create_datetime, create_user_id, update_datetime, update_user_id,
              active, active_datetime, non_active_datetime `).
		Add(" FROM ", tables.SS_DIAGNOSA).
		Add(" WHERE diagnosa_code = :diagnosa_code AND tenant_id = :tenant_id ").
		SetParam("tenant_id", input.Audit.TenantLoginId).
		SetParam("diagnosa_code", input.DiagnosaCode), &results)

	if err != nil {
		return false, nil, glerr.Wrap("error query is exists diagnosa: ", err)
	}

	if len(results) == 0 {
		return false, nil, nil
	}

	return true, results[0], nil
}
