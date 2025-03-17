package sstrxdao

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

	PatientId       int64  `json:"patientId" validate:"required"`
	DatetimeIn      string `json:"datetimeIn" validate:"required"`
	DatetimeProcess string `json:"datetimeProcess" validate:"required"`
	DatetimeOut     string `json:"datetimeOut" validate:"required"`
	MinDuration     int64  `json:"minDuration"`
	StatusDoc       string `json:"statusDoc" validate:"required"`
	DataJson        string `json:"dataJson" default:"{}"`
}

func Add(input InputAdd) (*tables.Trx, error) {
	result := tables.Trx{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TRX).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("patient_id", input.PatientId).
		AddValue("datetime_in", input.DatetimeIn).
		AddValue("datetime_process", input.DatetimeProcess).
		AddValue("datetime_out", input.DatetimeOut).
		AddValue("min_duration", input.MinDuration).
		AddValue("status_doc", input.StatusDoc).
		AddValue("data_json", input.DataJson).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" RETURNING ", result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert trx : " + err.Error())
	}

	return &result, nil
}
