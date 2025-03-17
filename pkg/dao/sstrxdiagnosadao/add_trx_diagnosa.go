package sstrxdiagnosadao

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputAdd struct {
	Tx    gldb.Tx
	Audit gldata.AuditData

	TrxId      int64 `json:"trxId" validate:"required"`
	DiagnosaId int64 `json:"diagnosaId" validate:"required"`
}

func Add(input InputAdd) (*tables.TrxDiagnosa, error) {
	result := tables.TrxDiagnosa{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TRX_DIAGNOSA).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("trx_id", input.TrxId).
		AddValue("diagnosa_id", input.DiagnosaId).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" RETURNING ", result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert trx diagnosa : " + err.Error())
	}

	return &result, nil
}
