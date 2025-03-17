package sspatientdao

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

	PatientId int64  `json:"patientId" validate:"required"`
	Idss      string `json:"idss" validate:"required"`
}

func UpdateIdss(input InputUpdateIdss) (*tables.Patient, error) {
	result := tables.Patient{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" UPDATE ", tables.SS_PATIENT, " SET ").
		AddSetNext("idss", input.Idss).
		AddPrepareUpdateAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" WHERE patient_id = :patientId ").
		SetParam("patientId", input.PatientId).
		Add(` RETURNING `, result.Columns()), &result)
	if err != nil {
		return nil, errors.New("failed update idss patient : " + err.Error())
	}

	return &result, nil
}
