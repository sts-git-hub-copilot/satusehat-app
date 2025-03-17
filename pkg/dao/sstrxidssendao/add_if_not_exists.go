package sstrxidssendao

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

	TrxId  int64  `json:"trxId" validate:"required"`
	Idssen string `json:"idssen" validate:"required"`
	Status string `json:"status" validate:"required"`
}

func AddIfNotExists(input InputAdd) (*tables.TrxIdssen, error) {
	result := tables.TrxIdssen{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TRX_IDSSEN).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("trx_id", input.TrxId).
		AddValue("idssen", input.Idssen).
		AddValue("status", input.Status).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` ON CONFLICT (trx_id, status) DO UPDATE 
            SET update_datetime = EXCLUDED.update_datetime,
                update_user_id = EXCLUDED.update_user_id,
				version = ss_trx_idssen.version + 1
            RETURNING `, result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert trx idssen : " + err.Error())
	}

	return &result, nil
}
