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

type InputCreateConditionQueue struct {
	TrxId int64
}

func CreateConditionQueueProcess(qtx glqueue.QueueContext) error {

	logrus.Println("\n\n --START Queue Create Condition -> ")
	input := InputCreateConditionQueue{}
	if err := qtx.Data.ToStruct(&input); err != nil {
		logrus.Error("error failed to decode input as InputCreateDiagnosaQueue", err)
		return err
	}

	if err := gldb.BeginTrxMt(qtx.Mt, func(tx gldb.Tx) error {

		conditionInput, err := utils.FetchConditionInput(tx, input.TrxId)
		if err != nil {
			return err
		}

		_, err = utils.CallApiGetConditionBySubject(utils.InputCallApiGetConditionBySubject{
			Tx:             tx,
			Audit:          qtx.AuditData,
			TrxId:          input.TrxId,
			ConditionInput: *conditionInput,
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
			QueueName:   constants.QUEUE_NAME_FINISH_ENCOUNTER,
			ProcessName: constants.QUEUE_NAME_FINISH_ENCOUNTER,
			Data:        goleafcore.NewOrEmpty(InputFinishEncounterQueue(input)),
		})
		if err != nil {
			return errors.New("error queue create diagnosa : " + err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
