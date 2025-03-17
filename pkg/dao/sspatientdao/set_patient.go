package sspatientdao

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glconstant"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"github.com/sirupsen/logrus"
)

type InputSet struct {
	Tx    gldb.Tx
	Audit gldata.AuditData

	Nik         string `json:"nik" validate:"required"`
	Idss        string `json:"idss"`
	PatientName string `json:"patientName" validate:"required"`
	Address     string `json:"address"`
	DataJson    string `json:"dataJson" default:"{}"`
}

func Set(input InputSet) (*tables.Patient, error) {
	logrus.Debug(" -- START Set Patient -> ")
	result := tables.Patient{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_PATIENT).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("nik", input.Nik).
		AddValue("idss", input.Idss).
		AddValue("patient_name", input.PatientName).
		AddValue("address", input.Address).
		AddValue("data_json", input.DataJson).
		AddValue("active", glconstant.YES).
		AddValue("active_datetime", input.Audit.Datetime()).
		AddValue("non_active_datetime", glconstant.EMPTY_VALUE).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` ON CONFLICT (tenant_id, nik) DO UPDATE 
            SET patient_name = EXCLUDED.patient_name,
                address = EXCLUDED.address,
                update_datetime = EXCLUDED.update_datetime,
                update_user_id = EXCLUDED.update_user_id
            RETURNING `, result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert patient if not exists : " + err.Error())
	}
	logrus.Debug("Patient Data : ", goleafcore.NewOrEmpty(result).PrettyString())
	logrus.Debug(" -- END Set Patient -> ")
	return &result, nil
}
