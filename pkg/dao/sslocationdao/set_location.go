package sslocationdao

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

	LocationName string `json:"locationName" validate:"required"`
	LocationCode string `json:"locationCode" validate:"required"`
	Idss         string `json:"idss"`
	Address      string `json:"address"`
	DataJson     string `json:"dataJson" default:"{}"`
}

func Set(input InputSet) (*tables.Location, error) {
	logrus.Debug(" -- START Set Location -> ")
	result := tables.Location{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_LOCATION).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("location_name", input.LocationName).
		AddValue("location_code", input.LocationCode).
		AddValue("idss", input.Idss).
		AddValue("address", input.Address).
		AddValue("data_json", input.DataJson).
		AddValue("active", glconstant.YES).
		AddValue("active_datetime", input.Audit.Datetime()).
		AddValue("non_active_datetime", glconstant.EMPTY_VALUE).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` ON CONFLICT (tenant_id, location_code) DO UPDATE 
            SET location_name = EXCLUDED.location_name,
                address = EXCLUDED.address,
                update_datetime = EXCLUDED.update_datetime,
                update_user_id = EXCLUDED.update_user_id
            RETURNING `, result.Columns()), &result)

	if err != nil {
		return nil, errors.New("failed insert location if not exists : " + err.Error())
	}
	logrus.Debug("Location Data : ", goleafcore.NewOrEmpty(result).PrettyString())
	logrus.Debug(" -- END Set Location -> ")
	return &result, nil
}
