package sslocationdao

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputIsExists struct {
	Tx    gldb.Tx
	Audit gldata.AuditData

	LocationCode string `json:"locationCode" validate:"required"`
}

func IsExists(input InputIsExists) (bool, *tables.Location, error) {
	results := make([]*tables.Location, 0)
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return false, nil, err
	}

	err := gldb.SelectQTx(input.Tx, *gldb.NewQBuilder().
		Add(` SELECT location_id, tenant_id, location_name, location_code, idss, address, data_json,
              version, create_datetime, create_user_id, update_datetime, update_user_id,
              active, active_datetime, non_active_datetime `).
		Add(" FROM ", tables.SS_LOCATION).
		Add(" WHERE location_name = :locationCode AND tenant_id = :tenantId ").
		SetParam("tenantId", input.Audit.TenantLoginId).
		SetParam("locationCode", input.LocationCode), &results)

	if err != nil {
		return false, nil, glerr.Wrap("error query is exists location: ", err)
	}

	if len(results) == 0 {
		return false, nil, nil
	}

	return true, results[0], nil
}
