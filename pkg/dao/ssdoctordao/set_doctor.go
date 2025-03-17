package ssdoctordao

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

	Nik        string `json:"nik" validate:"required"`
	Idss       string `json:"idss"`
	DoctorName string `json:"doctorName" validate:"required"`
	DataJson   string `json:"dataJson" default:"{}"`
}

func Set(input InputSet) (*tables.Doctor, error) {
	logrus.Debug(" -- START Set Doctor -> ")
	result := tables.Doctor{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_DOCTOR).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("nik", input.Nik).
		AddValue("idss", input.Idss).
		AddValue("doctor_name", input.DoctorName).
		AddValue("data_json", input.DataJson).
		AddValue("active", glconstant.YES).
		AddValue("active_datetime", input.Audit.Datetime()).
		AddValue("non_active_datetime", glconstant.EMPTY_VALUE).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` ON CONFLICT (tenant_id, nik) DO UPDATE 
            SET doctor_name = EXCLUDED.doctor_name,
                update_datetime = EXCLUDED.update_datetime,
                update_user_id = EXCLUDED.update_user_id
            RETURNING `, result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert doctor if not exists : " + err.Error())
	}
	logrus.Debug("Doctor Data : ", goleafcore.NewOrEmpty(result).PrettyString())
	logrus.Debug(" -- END Set Doctor -> ")
	return &result, nil
}
