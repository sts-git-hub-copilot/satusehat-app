package queueprocess

import (
	"errors"

	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/utils"
	"github.com/sirupsen/logrus"
)

type InputSendEncounterQueue struct {
	TrxId int64
}

func SendEncounterQueueProcess(qtx glqueue.QueueContext) error {

	logrus.Println("\n\n --START Queue Send Encounter -> ")
	input := InputSendEncounterQueue{}
	if err := qtx.Data.ToStruct(&input); err != nil {
		logrus.Error("error failed to decode input as InputSendEncounterQueue", err)
		return err
	}

	if err := gldb.BeginTrxMt(qtx.Mt, func(tx gldb.Tx) error {

		encounterInput, err := utils.FetchEncounterInput(tx, input.TrxId)
		if err != nil {
			return err
		}

		_, err = utils.CallApiGetEncounterBySubject(utils.InputCallApiGetEncounterBySubject{
			Tx:             tx,
			Audit:          qtx.AuditData,
			TrxId:          input.TrxId,
			EncounterInput: *encounterInput,
		})
		if err != nil {
			e := SendWebhookUpdateStatus(InputSendWebhookUpdateStatus{
				Mt:          qtx.Mt,
				Audit:       qtx.AuditData,
				TrxId:       input.TrxId,
				Status:      constants.STATUS_FAIL,
				FailMessage: err.Error(),
			})
			if e != nil {
				return e
			}
			return err
		}

		err = glqueue.Enqueue(glqueue.QueueData{
			Mt:          qtx.Mt,
			AuditData:   qtx.AuditData,
			QueueName:   constants.QUEUE_NAME_CREATE_CONDITION,
			ProcessName: constants.QUEUE_NAME_CREATE_CONDITION,
			Data:        goleafcore.NewOrEmpty(InputCreateConditionQueue(input)),
		})
		if err != nil {
			return errors.New("error queue create condition : " + err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

//example response
// {
//   "class": {
//     "code": "AMB",
//     "display": "ambulatory",
//     "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode"
//   },
//   "id": "23de416b-5fec-42ab-b8ef-91b08696365d",
//   "identifier": [
//     {
//       "system": "http://sys-ids.kemkes.go.id/encounter/a33ade0f-a26c-45e0-a19f-e206e1ccc255"
//     }
//   ],
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
//       }
//     }
//   ],
//   "meta": {
//     "lastUpdated": "2025-03-12T15:51:15.200265+00:00",
//     "versionId": "MTc0MTc5NDY3NTIwMDI2NTAwMA"
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
//     "start": "2025-03-11T03:56:49+00:00"
//   },
//   "resourceType": "Encounter",
//   "serviceProvider": {
//     "reference": "Organization/a33ade0f-a26c-45e0-a19f-e206e1ccc255"
//   },
//   "status": "arrived",
//   "statusHistory": [
//     {
//       "period": {
//         "end": "2025-03-11T03:56:49+00:00",
//         "start": "2025-03-11T03:56:49+00:00"
//       },
//       "status": "arrived"
//     }
//   ],
//   "subject": {
//     "display": "Ardianto Putra",
//     "reference": "Patient/P02478375538"
//   }
// }
