package sslocationdao

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputUpdateIdss struct {
	Tx         gldb.Tx
	Audit      gldata.AuditData
	LocationId int64  `json:"locationId" validate:"required"`
	Idss       string `json:"idss" validate:"required"`
}

func UpdateIdss(input InputUpdateIdss) (*tables.Location, error) {
	result := tables.Location{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" UPDATE ", tables.SS_LOCATION, " SET ").
		AddSetNext("idss", input.Idss).
		AddPrepareUpdateAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(" WHERE location_id = :locationId ").
		SetParam("locationId", input.LocationId).
		Add(` RETURNING `, result.Columns()), &result)
	if err != nil {
		return nil, errors.New("failed update idss location : " + err.Error())
	}

	return &result, nil
}
