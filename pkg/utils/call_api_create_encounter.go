package utils

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InputCreateEncounter struct {
	Tx             gldb.Tx
	Audit          gldata.AuditData
	TrxId          int64
	EncounterInput EncounterInput
}

type EncounterInput struct {
	IdssPatient  string `json:"idssPatient"`
	PatientName  string `json:"patientName"`
	IdssDoctor   string `json:"idssDoctor"`
	DoctorName   string `json:"doctorName"`
	IdssLocation string `json:"idssLocation"`
	LocationName string `json:"locationName"`
	OrgId        string `json:"orgId"`
	StartTime    string `json:"startTime"` //format 2024-06-27T12:43:53+07:00 standart start time untuk kirim data di konversi ke UTC
}

func CallApiCreateEncounter(input InputCreateEncounter) (goleafcore.Dto, error) {
	logrus.Debug(" -- START Call API Create Encounter -> ")

	patternBody := `{
    "resourceType": "Encounter",
    "identifier": [{
        "system": "http://sys-ids.kemkes.go.id/encounter/{{.OrgId}}",
        "value": "{{.IdssPatient}}.{{.StartTime}}"
    }],
    "status": "arrived",
    "class": {
        "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
        "code": "AMB",
        "display": "ambulatory"
    },
    "subject": {
        "reference": "Patient/{{.IdssPatient}}",
        "display": "{{.PatientName}}"
    },
    "participant": [{
        "type": [{
            "coding": [{
                "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
                "code": "ATND",
                "display": "attender"
            }]
        }],
        "individual": {
            "reference": "Practitioner/{{.IdssDoctor}}",
            "display": "{{.DoctorName}}"
        }
    }],
    "period": {
        "start": "{{.StartTime}}"
    },
    "location": [{
        "location": {
            "reference": "Location/{{.IdssLocation}}",
            "display": "{{.LocationName}}"
        },
        "extension": [{
            "url": "https://fhir.kemkes.go.id/r4/StructureDefinition/ServiceClass",
            "extension": [{
                "url": "value",
                "valueCodeableConcept": {
                    "coding": [{
                        "system": "http://terminology.kemkes.go.id/CodeSystem/locationServiceClass-Outpatient",
                        "code": "reguler",
                        "display": "Kelas Reguler"
                    }]
                }
            }]
        }]
    }],
    "statusHistory": [
        {
            "status": "arrived",
            "period": {
                "start": "{{.StartTime}}",
                "end": "{{.StartTime}}"
            }
        }
    ],
    "serviceProvider": {
        "reference": "Organization/{{.OrgId}}"
    }
}`

	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Encounter",
		Method: fiber.MethodPost,
		Body:   goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(patternBody, input.EncounterInput)),
	})
	if err != nil {
		return nil, err
	}

	logrus.Debug(" -- END Call API Create Encounter -> ")
	return output, nil
}

func FetchEncounterInput(tx gldb.Tx, trxId int64) (*EncounterInput, error) {
	result := EncounterInput{}

	err := gldb.SelectRowQTx(tx, *gldb.NewQBuilder().
		Add(` SELECT fss_get_idss_patient(A.patient_id) AS idss_patient,
                     fss_get_patient_name(A.patient_id) AS patient_name,
                        fss_get_idss_doctor(B.doctor_id) AS idss_doctor,
                        fss_get_doctor_name(B.doctor_id) AS doctor_name,
                        fss_get_idss_location(B.location_id) AS idss_location,
                        fss_get_location_name(B.location_id) AS location_name,
                        fss_get_org_id(A.tenant_id) AS org_id,
                        datetime_in AS start_time
              FROM `, tables.SS_TRX, ` A
              INNER JOIN `, tables.SS_TRX_ITEM, ` B ON B.trx_id = A.trx_id
              WHERE A.trx_id = :trxId `).
		SetParam("trxId", trxId), &result)
	if err != nil {
		return nil, glerr.New("error failed to fetch encounter input", err.Error())
	}

	result.StartTime = glutil.ParseDateOrDefault(result.StartTime).Format("2006-01-02T15:04:05-07:00")

	logrus.Debug(" -- Fetch Encounter Input -> ", goleafcore.NewOrEmpty(result).ToJsonString())

	return &result, nil
}
