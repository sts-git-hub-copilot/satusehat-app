package tables

import (
	"time"

	"git.solusiteknologi.co.id/goleaf/goleafcore/glentity"
)

type Token struct {
	TokenId      int64     `json:"tokenId"`
	TenantId     int64     `json:"tenantId"`
	AccessToken  string    `json:"accessToken"`
	IssuedAt     time.Time `json:"issuedAt"`
	ExpiresIn    int64     `json:"expiresIn"`
	ExpiredToken time.Time `json:"expiredToken"`
	DataJson     string    `json:"dataJson"`

	glentity.BaseEntity
}

func (d Token) TableName() string {
	return SS_TOKEN
}
