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

type InputFinishEncounter struct {
	Tx                   gldb.Tx
	Audit                gldata.AuditData
	FinishEncounterInput FinishEncounterInput
}

type FinishEncounterInput struct {
	IdssEncounter string `json:"idssEncounter"`
	IdssCondition string `json:"idssCondition"`
	OrgId         string `json:"orgId"`
	IdssPatient   string `json:"idssPatient"`
	PatientName   string `json:"patientName"`
	IdssDoctor    string `json:"idssDoctor"`
	DoctorName    string `json:"doctorName"`
	IdssLocation  string `json:"idssLocation"`
	LocationName  string `json:"locationName"`
	StartTime     string `json:"startTime"`
	ProcessTime   string `json:"processTime"`
	EndTime       string `json:"endTime"`
	Duration      int64  `json:"duration"`
}

func CallApiFinishEncounter(input InputFinishEncounter) (goleafcore.Dto, error) {
	patternBody := `{
    "resourceType": "Encounter",
    "id": "{{.IdssEncounter}}",
    "identifier": [{
        "system": "http://sys-ids.kemkes.go.id/encounter/{{.OrgId}}"
    }],
    "status": "finished",
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
        "start": "{{.StartTime}}",
        "end": "{{.EndTime}}"
    },
    "length": {
        "value": {{.Duration}},
        "unit": "min",
        "system": "http://unitsofmeasure.org",
        "code": "min"
    },
    "location": [{
        "location": {
            "reference": "Location/{{.IdssLocation}}",
            "display": "{{.LocationName}}"
        },
        "period": {
            "start": "{{.ProcessTime}}",
            "end": "{{.ProcessTime}}"
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
    "diagnosis": [
        {
            "condition": {
                "reference": "Condition/{{.IdssCondition}}"
            },
            "use": {
                "coding": [{
                    "system": "http://terminology.hl7.org/CodeSystem/diagnosis-role",
                    "code": "CC",
                    "display": "Chief Complaint"
                }]
            }
        }
    ],
    "statusHistory": [
        {
            "status": "arrived",
            "period": {
                "start": "{{.StartTime}}",
                "end": "{{.StartTime}}"
            }
        },
        {
            "status": "in-progress",
            "period": {
                "start": "{{.ProcessTime}}",
                "end": "{{.EndTime}}"
            }
        },
        {
            "status": "finished",
            "period": {
                "start": "{{.EndTime}}",
                "end": "{{.EndTime}}"
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
		Path:   "/Encounter/" + input.FinishEncounterInput.IdssEncounter,
		Method: fiber.MethodPut,
		Body:   goleafcore.NewOrEmpty(glutil.ParseTemplateOrDefault(patternBody, input.FinishEncounterInput)),
	})

	return output, err
}

func FetchFinishEncounterInput(tx gldb.Tx, trxId int64) (*FinishEncounterInput, error) {
	result := FinishEncounterInput{}

	err := gldb.SelectRowQTx(tx, *gldb.NewQBuilder().
		Add(` SELECT C.idssen AS idss_encounter, D.idssen AS idss_condition,
                     fss_get_org_id(A.tenant_id) AS org_id,
                     fss_get_idss_patient(A.patient_id) AS idss_patient,
                     fss_get_patient_name(A.patient_id) AS patient_name,
                     fss_get_idss_doctor(B.doctor_id) AS idss_doctor,
                     fss_get_doctor_name(B.doctor_id) AS doctor_name,
                     fss_get_idss_location(B.location_id) AS idss_location,
                     fss_get_location_name(B.location_id) AS location_name,
                     datetime_in AS start_time, datetime_out AS end_time, datetime_process AS process_time,
                     min_duration AS duration
              FROM `, tables.SS_TRX, ` A
              INNER JOIN `, tables.SS_TRX_ITEM, ` B ON B.trx_id = A.trx_id
              INNER JOIN `, tables.SS_TRX_IDSSEN, ` C ON C.trx_id = A.trx_id AND C.status = :ENCOUNTER
              INNER JOIN `, tables.SS_TRX_IDSSEN, ` D ON D.trx_id = A.trx_id AND D.status = :CONDITION
              WHERE A.trx_id = :trxId `).
		SetParam("ENCOUNTER", constants.COMBO_VALUE_ENCOUNTER).
		SetParam("CONDITION", constants.COMBO_VALUE_CONDITION).
		SetParam("trxId", trxId), &result)
	if err != nil {
		return nil, glerr.New("error failed to fetch finish encounter input", err.Error())
	}

	result.StartTime = glutil.ParseDateOrDefault(result.StartTime).Format("2006-01-02T15:04:05-07:00")
	result.ProcessTime = glutil.ParseDateOrDefault(result.ProcessTime).Format("2006-01-02T15:04:05-07:00")
	result.EndTime = glutil.ParseDateOrDefault(result.EndTime).Format("2006-01-02T15:04:05-07:00")

	logrus.Debug(" -- Fetch Finish Encounter Input -> ", goleafcore.NewOrEmpty(result).ToJsonString())

	return &result, nil
}

// example response
// {
//   "class": {
//     "code": "AMB",
//     "display": "ambulatory",
//     "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode"
//   },
//   "diagnosis": [
//     {
//       "condition": {
//         "reference": "Condition/3938a300-84e8-4583-9b12-ed5c5e83215a"
//       },
//       "use": {
//         "coding": [
//           {
//             "code": "CC",
//             "display": "Chief Complaint",
//             "system": "http://terminology.hl7.org/CodeSystem/diagnosis-role"
//           }
//         ]
//       }
//     }
//   ],
//   "id": "6a25f268-8179-4592-a2ea-e405419fc80f",
//   "identifier": [
//     {
//       "system": "http://sys-ids.kemkes.go.id/encounter/a33ade0f-a26c-45e0-a19f-e206e1ccc255"
//     }
//   ],
//   "length": {
//     "code": "min",
//     "system": "http://unitsofmeasure.org",
//     "unit": "min",
//     "value": 63
//   },
//   "location": [
//     {
//       "extension": [
//         {
//           "extension": [
//             {
//               "url": "value",
//               "valueCodeableConcept": {
//                 "coding": [
//                   {
//                     "code": "reguler",
//                     "display": "Kelas Reguler",
//                     "system": "http://terminology.kemkes.go.id/CodeSystem/locationServiceClass-Outpatient"
//                   }
//                 ]
//               }
//             }
//           ],
//           "url": "https://fhir.kemkes.go.id/r4/StructureDefinition/ServiceClass"
//         }
//       ],
//       "location": {
//         "display": "Klinik Sehat Sejahtera",
//         "reference": "Location/6e7cc6bc-ed60-4eb8-a4f1-b48405c54c04"
//       },
//       "period": {
//         "end": "2025-03-11T04:00:00+00:00",
//         "start": "2025-03-11T04:00:00+00:00"
//       }
//     }
//   ],
//   "meta": {
//     "lastUpdated": "2025-03-13T07:38:53.410532+00:00",
//     "versionId": "MTc0MTg1MTUzMzQxMDUzMjAwMA"
//   },
//   "participant": [
//     {
//       "individual": {
//         "display": "dr. Alexander",
//         "reference": "Practitioner/10009880728"
//       },
//       "type": [
//         {
//           "coding": [
//             {
//               "code": "ATND",
//               "display": "attender",
//               "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType"
//             }
//           ]
//         }
//       ]
//     }
//   ],
//   "period": {
//     "end": "2025-03-11T05:00:00+00:00",
//     "start": "2025-03-11T03:56:49+00:00"
//   },
//   "resourceType": "Encounter",
//   "serviceProvider": {
//     "reference": "Organization/a33ade0f-a26c-45e0-a19f-e206e1ccc255"
//   },
//   "status": "finished",
//   "statusHistory": [
//     {
//       "period": {
//         "end": "2025-03-11T03:56:49+00:00",
//         "start": "2025-03-11T03:56:49+00:00"
//       },
//       "status": "arrived"
//     },
//     {
//       "period": {
//         "end": "2025-03-11T05:00:00+00:00",
//         "start": "2025-03-11T04:00:00+00:00"
//       },
//       "status": "in-progress"
//     },
//     {
//       "period": {
//         "end": "2025-03-11T05:00:00+00:00",
//         "start": "2025-03-11T05:00:00+00:00"
//       },
//       "status": "finished"
//     }
//   ],
//   "subject": {
//     "display": "Ardianto Putra",
//     "reference": "Patient/P02478375538"
//   }
// }
