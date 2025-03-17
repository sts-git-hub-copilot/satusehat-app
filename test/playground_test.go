package test_test

import (
	"log"
	"testing"
	"time"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gltest"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func TestTemplateUtil(t *testing.T) {

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
	"telecom": [
		{
			"system": "phone",
			"value": "{{.Phone}}",
			"use": "work"
		},
		{
			"system": "email",
			"value": "{{.Email}}",
			"use": "work"
		},
		{
			"system": "url",
			"value": "{{.UrlWeb}}",
			"use": "work"
		}
	],
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

	location := utils.LocationInput{
		OrgId:        "1",
		OuCode:       "2",
		LocationName: "Ruangan Estetika",
		Address:      "Jl Mt haryono",
		Phone:        "0832312321333",
		Email:        "mail@mail.com",
		UrlWeb:       "manikaaestheticclinic.solusi-iclinic.id",
		Longitude:    decimal.Zero,
		Latitude:     decimal.Zero,
	}

	test2 := utils.ConditionInput{
		DiagnosaList: []*tables.Diagnosa{
			{DiagnosaCode: "1", DiagnosaName: "Diagnosa 1"},
			{DiagnosaCode: "2", DiagnosaName: "Diagnosa 2"},
		},
	}

	templ2 := `[{{range $index, $diag := .DiagnosaList}}{
                "system": "http://hl7.org/fhir/sid/icd-10",
                "code": "{{$diag.DiagnosaCode}}",
                "display": "{{$diag.DiagnosaName}}"
            }{{end}}]`

	log.Println(" cek : ", goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(templ2, test2)).PrettyString())

	log.Println(" cek : ", goleafcore.NewOrEmpty(location).ToJsonString())

	log.Println("dto : ", goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(patternBody, location)).PrettyString())
}

func TestPlayGround(t *testing.T) {
	gltest.TestApi(t, func(app *fiber.App, tx gldb.Tx) error {

		return nil
	}, func(assert *gltest.Assert, app *fiber.App, tx gldb.Tx, i int) interface{} {

		logrus.Debug(glutil.ParseDateOrDefault("20250311105649").Format("2006-01-02T15:04:05-07:00"))

		logrus.Debug(" cekk : ", goleafcore.Dto{
			"sample": time.Now().UTC().Format("2006-01-02T15:04:05-07:00"),
		}.PrettyString())

		logrus.Debug(time.UnixMilli(1741665409000).UTC().Format("2006-01-02T15:04:05-07:00"))

		diagnosaList := []tables.Diagnosa{
			{DiagnosaCode: "1", DiagnosaName: "Diagnosa 1"},
			// {DiagnosaCode: "2", DiagnosaName: "Diagnosa 2"},
			// {DiagnosaCode: "3", DiagnosaName: "Diagnosa 3"},
		}

		template := ""
		for i, diag := range diagnosaList {
			temp := `{
				"system": "http://hl7.org/fhir/sid/icd-10",
				"code": "{{.DiagnosaCode}}",
				"display": "{{.DiagnosaName}}"
			}`
			template += glutil.ParseTemplateOrDefault(temp, diag)
			if i < len(diagnosaList)-1 {
				template += ","
			}
		}
		logrus.Debug(" Diagnosa Temp -> ", template)

		// Duitku signature
		// signatureData := fmt.Sprint("DS20818", 50000, "kSybhpvPoYvThHWjRZSGtfEPKLTPRYMs", "42bb9eba2a841488092f5ae3763f4cc4")
		// hash := md5.New()
		// hash.Write([]byte(signatureData))
		// hashBytes := hash.Sum(nil)
		// Encode ke hexadecimal
		// signature := hex.EncodeToString(hashBytes)

		// logrus.Debug("signature : ", signature)
		// logrus.Debug("signature : 326cba062627f92ac91d91ee48b6c229")

		// logrus.Debug("test : ", goleafcore.Dto{
		// 	"sample": "",
		// }.GetInt64("sample"))

		return nil
	})
}
