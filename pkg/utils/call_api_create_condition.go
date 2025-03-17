package utils

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InputCreateCondition struct {
	Tx             gldb.Tx
	Audit          gldata.AuditData
	TrxId          int64
	ConditionInput ConditionInput
}

type ConditionInput struct {
	DiagnosaList  []*tables.Diagnosa `json:"diagnosaList"`
	IdssPatient   string             `json:"idssPatient"`
	PatientName   string             `json:"patientName"`
	IdssEncounter string             `json:"idssEncounter"`
	ProcessTime   string             `json:"processTime"` // format: 2023-06-04T08:30:00+00:00
	DiagnosaTemp  string             `json:"diagnosaTemp"`
}

func CallApiCreateCondition(input InputCreateCondition) (goleafcore.Dto, error) {
	logrus.Debug(" -- START Call API Create Condition -> ")

	patternBody := `{
    "resourceType": "Condition",
    "clinicalStatus": {
        "coding": [{
            "system": "http://terminology.hl7.org/CodeSystem/condition-clinical",
            "code": "active",
            "display": "Active"
        }]
    },
    "category": [{
        "coding": [{
            "system": "http://terminology.hl7.org/CodeSystem/condition-category",
            "code": "encounter-diagnosis",
            "display": "Encounter Diagnosis"
        }]
    }],
    "code": {
        "coding": [
            {{.DiagnosaTemp}}
        ]
    },
    "subject": {
        "reference": "Patient/{{.IdssPatient}}",
        "display": "{{.PatientName}}"
    },
    "encounter": {
        "reference": "Encounter/{{.IdssEncounter}}"
    },
    "onsetDateTime": "{{.ProcessTime}}",
    "recordedDate": "{{.ProcessTime}}"
}`

	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Condition",
		Method: fiber.MethodPost,
		Body:   goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(patternBody, input.ConditionInput)),
	})
	if err != nil {
		return nil, err
	}

	logrus.Debug(" -- END Call API Create Condition -> ")
	return output, nil
}

func FetchConditionInput(tx gldb.Tx, trxId int64) (*ConditionInput, error) {
	result := ConditionInput{}

	err := gldb.SelectRowQTx(tx, *gldb.NewQBuilder().
		Add(` SELECT fss_get_idss_patient(A.patient_id) AS idss_patient,
                     fss_get_patient_name(A.patient_id) AS patient_name,
                     B.idssen AS idss_encounter,
                     A.datetime_process AS process_time
              FROM `, tables.SS_TRX, ` A
              INNER JOIN `, tables.SS_TRX_IDSSEN, ` B ON B.trx_id = A.trx_id AND B.status = :ENCOUNTER
              WHERE A.trx_id = :trxId `).
		SetParam("ENCOUNTER", constants.COMBO_VALUE_ENCOUNTER).
		SetParam("trxId", trxId), &result)
	if err != nil {
		return nil, glerr.New("error failed to fetch encounter input", err.Error())
	}

	result.ProcessTime = glutil.ParseDateOrDefault(result.ProcessTime).Format("2006-01-02T15:04:05-07:00")

	diagnosaList, err := FetchDiagnosaList(tx, trxId)
	if err != nil {
		return nil, err
	}
	result.DiagnosaList = diagnosaList

	for i, diag := range diagnosaList {
		temp := `{
				"system": "http://hl7.org/fhir/sid/icd-10",
				"code": "{{.DiagnosaCode}}",
				"display": "{{.DiagnosaName}}"
			}`
		result.DiagnosaTemp += glutil.ParseTemplateOrDefault(temp, diag)
		if i < len(diagnosaList)-1 {
			result.DiagnosaTemp += ","
		}
	}
	logrus.Debug(" Diagnosa Temp -> ", result.DiagnosaTemp)

	logrus.Debug(" -- Fetch Condition Input -> ", goleafcore.NewOrEmpty(result).ToJsonString())

	return &result, nil
}

func FetchDiagnosaList(tx gldb.Tx, trxId int64) ([]*tables.Diagnosa, error) {
	result := []*tables.Diagnosa{}

	err := gldb.SelectQTx(tx, *gldb.NewQBuilder().
		Add(` SELECT fss_get_diagnosa_code(A.diagnosa_id) diagnosa_code,
                    fss_get_diagnosa_name(A.diagnosa_id) diagnosa_name
              FROM `, tables.SS_TRX_DIAGNOSA, ` A
              WHERE A.trx_id = :trxId `).
		SetParam("trxId", trxId), &result)
	if err != nil {
		return nil, glerr.New("error failed to fetch diagnosa list", err.Error())
	}

	return result, nil
}
