package utils

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// InputCallApiSS menyusun parameter untuk pemanggilan API satusehat.
type InputCallApiCreateLoc struct {
	Tx         gldb.Tx
	Audit      gldata.AuditData
	LocationId int64
}

type LocationInput struct {
	OrgId        string          `json:"orgId"`        // Digunakan untuk pengisian {{.OrgId}}
	OuCode       string          `json:"ouCode"`       // Digunakan untuk pengisian {{.OuCode}}
	LocationName string          `json:"locationName"` // Digunakan untuk pengisian {{.LocationName}}
	Address      string          `json:"address"`      // Digunakan untuk pengisian {{.Address}}
	Phone        string          `json:"phone"`        // Digunakan untuk pengisian {{.Phone}}
	Email        string          `json:"email"`        // Digunakan untuk pengisian {{.Email}}
	UrlWeb       string          `json:"urlWeb"`       // Digunakan untuk pengisian {{.UrlWeb}}
	Longitude    decimal.Decimal `json:"longitude"`    // Digunakan untuk pengisian {{.Longitude}}
	Latitude     decimal.Decimal `json:"latitude"`     // Digunakan untuk pengisian {{.Latitude}}
}

// CallApiSS melakukan pemanggilan ke endpoint base url.
func CallApiCreateLoc(input InputCallApiCreateLoc) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat -> ")

	locationInput, err := getLocationInput(input)
	if err != nil {
		return nil, err
	}

	patternBody := `{
	"resourceType": "Location",
	"identifier": [
		{
			"system": "http://sys-ids.kemkes.go.id/location/{{.OrgId}}",
			"value": "{{.OuCode}}"
		}
	],
	"status": "active",
	"name": "{{.LocationName}}",
	"description": "{{.Address}}",
	"mode": "instance",
	{{if or .Phone .Email .UrlWeb}}
	"telecom": [
		{{if .Phone}}
		{
			"system": "phone",
			"value": "{{.Phone}}",
			"use": "work"
		},
		{{end}}
		{{if .Email}}
		{
			"system": "email",
			"value": "{{.Email}}",
			"use": "work"
		},
		{{end}}
		{{if .UrlWeb}}
		{
			"system": "url",
			"value": "{{.UrlWeb}}",
			"use": "work"
		}
		{{end}}
	],
	{{end}}
	"physicalType": {
		"coding": [
			{
				"system": "http://terminology.hl7.org/CodeSystem/location-physical-type",
				"code": "ro",
				"display": "Room"
			}
		]
	},
	{{if and (not .Longitude.IsZero) (not .Latitude.IsZero)}}
	"position": {
		"longitude": {{.Longitude}},
		"latitude": {{.Latitude}},
		"altitude": 0
	},
	{{end}}
	"managingOrganization": {
		"reference": "Organization/{{.OrgId}}"
	}
}`

	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Location",
		Method: fiber.MethodPost,
		Body:   goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(patternBody, *locationInput)),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func getLocationInput(input InputCallApiCreateLoc) (*LocationInput, error) {
	result := make([]*LocationInput, 0)

	err := gldb.SelectQTx(input.Tx, *gldb.NewQBuilder().
		Add(` SELECT COALESCE(f_get_value_system_config_by_param_code(:tenantId, :PARAM_CODE),'') org_id, fss_get_ou_code_company() ou_code, location_name, address,
			COALESCE(data_json::jsonb->>'phone','') phone, COALESCE(data_json::jsonb->>'email','') email, 
			COALESCE(data_json::jsonb->>'urlWeb','') url_web,
			COALESCE((data_json::jsonb->>'longitude')::numeric,0) longitude, 
			COALESCE((data_json::jsonb->>'latitude')::numeric,0) latitude
			FROM `, tables.SS_LOCATION, `
			WHERE location_id = :locationId `).
		SetParam("tenantId", input.Audit.TenantLoginId).
		SetParam("PARAM_CODE", constants.PARAM_CODE_SS_ORGANIZATION_ID).
		SetParam("locationId", input.LocationId), &result)
	if err != nil {
		return nil, errors.New("failed get location input : " + err.Error())
	}

	if len(result) == 0 {
		return nil, glerr.New("location not found")
	}

	return result[0], nil
}

// {
//   "description": "Jl Mt haryono",
//   "id": "f1ff604e-0900-48db-8ef8-b59832314640",
//   "identifier": [
//     {
//       "system": "http://sys-ids.kemkes.go.id/location/a33ade0f-a26c-45e0-a19f-e206e1ccc255",
//       "value": "manikaaestheticclinic"
//     }
//   ],
//   "managingOrganization": {
//     "reference": "Organization/a33ade0f-a26c-45e0-a19f-e206e1ccc255"
//   },
//   "meta": {
//     "lastUpdated": "2025-02-26T07:36:30.765804+00:00",
//     "versionId": "MTc0MDU1NTM5MDc2NTgwNDAwMA"
//   },
//   "mode": "instance",
//   "name": "Ruangan Estetika",
//   "physicalType": {
//     "coding": [
//       {
//         "code": "ro",
//         "display": "Room",
//         "system": "http://terminology.hl7.org/CodeSystem/location-physical-type"
//       }
//     ]
//   },
//   "resourceType": "Location",
//   "status": "active"
// }
