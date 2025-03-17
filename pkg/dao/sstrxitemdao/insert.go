package sstrxitemdao

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

	TrxId      int64  `json:"trxId" validate:"required"`
	LineNo     int64  `json:"lineNo" validate:"required"`
	DoctorId   int64  `json:"doctorId" validate:"required"`
	LocationId int64  `json:"locationId" validate:"required"`
	DataJson   string `json:"dataJson" default:"{}"`
}

func Add(input InputAdd) (*tables.TrxItem, error) {
	result := tables.TrxItem{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TRX_ITEM).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("trx_id", input.TrxId).
		AddValue("line_no", input.LineNo).
		AddValue("doctor_id", input.DoctorId).
		AddValue("location_id", input.LocationId).
		AddValue("data_json", input.DataJson).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" RETURNING ", result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert trx item : " + err.Error())
	}

	return &result, nil
}
