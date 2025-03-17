package sstokendao

import (
	"errors"
	"time"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputAddToken struct {
	Tx          gldb.Tx
	Audit       gldata.AuditData
	AccessToken string    `json:"accessToken" validate:"required"`
	IssuedAt    time.Time `json:"issuedAt" validate:"required"`
	ExpiresIn   int64     `json:"expiresIn" validate:"required"` // Expiry time in seconds
	DataJson    string    `json:"dataJson" default:"{}"`
}

func Add(input InputAddToken) (*tables.Token, error) {
	result := tables.Token{}
	if err := glutil.FetchDefAndValidate(&input); err != nil {
		return nil, err
	}

	// Calculate expiredToken based on issuedAt + expiresIn
	expiredToken := input.IssuedAt.Add(time.Duration(input.ExpiresIn) * time.Second)

	err := gldb.SelectOneQTx(input.Tx, *gldb.NewQBuilder().
		Add(" INSERT INTO ", tables.SS_TOKEN).
		AddValue("tenant_id", input.Audit.TenantLoginId).
		AddValue("access_token", input.AccessToken).
		AddValue("issued_at", input.IssuedAt).
		AddValue("expires_in", input.ExpiresIn).
		AddValue("expired_token", expiredToken).
		AddValue("data_json", input.DataJson).
		AddValueInsertAudit(input.Audit.Datetime(), input.Audit.UserLoginId).
		Add(` RETURNING token_id, tenant_id, access_token, issued_at, expires_in, 
			 expired_token, data_json, create_datetime, update_datetime, 
			 create_user_id, update_user_id, version `), &result)

	if err != nil {
		return nil, errors.New("failed insert token: " + err.Error())
	}

	return &result, nil
}
