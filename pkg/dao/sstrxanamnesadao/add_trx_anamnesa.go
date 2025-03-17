package sstrxanamnesadao

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
	AnamnesaId int64 `json:"anamnesaId" validate:"required"`
}

func Add(input InputAdd) (*tables.TrxAnamnesa, error) {
	result := tables.TrxAnamnesa{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TRX_AMNANESA).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("trx_id", input.TrxId).
		AddValue("anamnesa_id", input.AnamnesaId).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" RETURNING ", result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert trx anamnesa : " + err.Error())
	}

	return &result, nil
}
