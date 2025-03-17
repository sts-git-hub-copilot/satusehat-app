package ssanamnesadao

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

	AnamnesaCode    string `json:"anamnesaCode" validate:"required"`
	AnamnesaName    string `json:"anamnesaName" validate:"required"`
	AnamnesaVersion string `json:"anamnesaVersion"`
}

func Set(input InputSet) (*tables.Anamnesa, error) {
	logrus.Debug(" -- START Set Anamnesa -> ")
	result := tables.Anamnesa{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_AMNANESA).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("anamnesa_code", input.AnamnesaCode).
		AddValue("anamnesa_name", input.AnamnesaName).
		AddValue("anamnesa_version", input.AnamnesaVersion).
		AddValue("active", glconstant.YES).
		AddValue("active_datetime", input.Audit.Datetime()).
		AddValue("non_active_datetime", glconstant.EMPTY_VALUE).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` ON CONFLICT (tenant_id, anamnesa_code) DO UPDATE 
            SET anamnesa_name = EXCLUDED.anamnesa_name,
                anamnesa_version = EXCLUDED.anamnesa_version,
                update_datetime = EXCLUDED.update_datetime,
                update_user_id = EXCLUDED.update_user_id
            RETURNING `, result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert anamnesa if not exists : " + err.Error())
	}
	logrus.Debug("Anamnesa Data : ", goleafcore.NewOrEmpty(result).PrettyString())
	logrus.Debug(" -- END Set Anamnesa -> ")
	return &result, nil
}
