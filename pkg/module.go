package satusehat

import (
	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/glwebhook"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glapi"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/controller"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/queueprocess"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	ApiPrefix   string // prefix
	Middleware  fiber.Handler
	QueueConfig glqueue.Config
}

func Setup(app *fiber.App, config Config) {
	config = buildWithDefault(config)

	setupQueueProcess(config)

	var v1ss = config.ApiPrefix + "/v1/satusehat"

	groupSendData := app.Group(v1ss+"/send-data", config.Middleware)
	groupSendData.Post("/send-list", controller.SendDataSSList)

}

func buildWithDefault(config Config) Config {
	return Config{
		ApiPrefix:  glutil.EmptyOrDefault(config.ApiPrefix, glapi.API_PREFIX_DEFAULT),
		Middleware: glapi.OrEmptyMiddleware(config.Middleware),
	}
}

func setupQueueProcess(config Config) {
	glqueue.Setup(config.QueueConfig)

	glqueue.Register(&glqueue.QueueProcess{
		QueueName:  constants.QUEUE_NAME_WEBHOOK_UPDATE_STATUS,
		Process:    glwebhook.QueWebhook,
		Concurrent: 1,
	})

	glqueue.Register(&glqueue.QueueProcess{
		QueueName:  constants.QUEUE_NAME_CREATE_ENCOUNTER,
		Process:    queueprocess.SendEncounterQueueProcess,
		Concurrent: 1,
	})

	glqueue.Register(&glqueue.QueueProcess{
		QueueName:  constants.QUEUE_NAME_CREATE_CONDITION,
		Process:    queueprocess.CreateConditionQueueProcess,
		Concurrent: 1,
	})

	glqueue.Register(&glqueue.QueueProcess{
		QueueName:  constants.QUEUE_NAME_FINISH_ENCOUNTER,
		Process:    queueprocess.FinishEncounterQueueProcess,
		Concurrent: 1,
	})

	glqueue.Register(&glqueue.QueueProcess{
		QueueName: constants.QUEUE_NAME_WEBHOOK_UPDATE_STATUS,
		Process:   glwebhook.QueWebhook,
	})
}
