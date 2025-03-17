package queueprocess

import (
	"context"
	"errors"

	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glconstant"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
)

type InputSendWebhookUpdateStatus struct {
	Mt          context.Context
	Audit       gldata.AuditData
	TrxId       int64  `json:"trxId"`
	FailMessage string `json:"failMessage"`
	Status      string `json:"status"`
	DataJson    string `json:"dataJson"`
}

func SendWebhookUpdateStatus(input InputSendWebhookUpdateStatus) error {

	dataJson := input.DataJson
	if input.TrxId != glconstant.NULL_REF_VALUE_FOR_LONG {
		err := gldb.SelectRowQMt(input.Mt, *gldb.NewQBuilder().
			Add(` SELECT data_json FROM `, tables.SS_TRX, ` WHERE trx_id = :trxId `).
			SetParam("trxId", input.TrxId), &dataJson)
		if err != nil {
			return glerr.New("fail get data json trx ", err.Error())
		}

	}
	err := glqueue.Enqueue(glqueue.QueueData{
		Mt:          input.Mt,
		AuditData:   input.Audit,
		QueueName:   constants.QUEUE_NAME_WEBHOOK_UPDATE_STATUS,
		ProcessName: constants.QUEUE_NAME_WEBHOOK_UPDATE_STATUS,
		Data: goleafcore.NewOrEmpty(goleafcore.Dto{
			"eventCode": constants.EVENT_UPDATE_STATUS,
			"changeLog": []goleafcore.Dto{
				{
					"failMessage": input.FailMessage,
					"status":      input.Status,
					"dataJson":    dataJson,
				},
			},
		}),
	})
	if err != nil {
		return errors.New("error queue webhook update status : " + err.Error())
	}

	return nil
}
