package queueprocess

import (
	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/utils"
	"github.com/sirupsen/logrus"
)

type InputFinishEncounterQueue struct {
	TrxId int64
}

func FinishEncounterQueueProcess(qtx glqueue.QueueContext) error {

	logrus.Println("\n\n --START Queue Finish Encounter -> ")
	input := InputFinishEncounterQueue{}
	if err := qtx.Data.ToStruct(&input); err != nil {
		logrus.Error("error failed to decode input as InputFinishEncounterQueue", err)
		return err
	}

	if err := gldb.BeginTrxMt(qtx.Mt, func(tx gldb.Tx) error {

		finishInput, err := utils.FetchFinishEncounterInput(tx, input.TrxId)
		if err != nil {
			return err
		}

		_, err = utils.CallApiFinishEncounter(utils.InputFinishEncounter{
			Tx:                   tx,
			Audit:                qtx.AuditData,
			FinishEncounterInput: *finishInput,
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

		err = SendWebhookUpdateStatus(InputSendWebhookUpdateStatus{
			Mt:     qtx.Mt,
			Audit:  qtx.AuditData,
			TrxId:  input.TrxId,
			Status: constants.STATUS_SUCCESS,
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
