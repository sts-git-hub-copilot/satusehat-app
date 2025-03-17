package ssdoctordao

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputUpdateIdss struct {
	Tx    gldb.Tx
	Audit gldata.AuditData

	DoctorId int64  `json:"doctorId" validate:"required"`
	Idss     string `json:"idss" validate:"required"`
}

func UpdateIdss(input InputUpdateIdss) (*tables.Doctor, error) {
	result := tables.Doctor{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" UPDATE ", tables.SS_DOCTOR, " SET ").
		AddSetNext("idss", input.Idss).
		AddPrepareUpdateAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" WHERE doctor_id = :doctorId ").
		SetParam("doctorId", input.DoctorId).
		Add(` RETURNING `, result.Columns()), &result)
	if err != nil {
		return nil, errors.New("failed update idss doctor : " + err.Error())
	}

	return &result, nil
}
