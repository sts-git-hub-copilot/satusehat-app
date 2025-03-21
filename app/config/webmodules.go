package config

import (
	"strings"
	"time"

	"git.solusiteknologi.co.id/goleaf/glauth"
	authconst "git.solusiteknologi.co.id/goleaf/glauth/constants"
	"git.solusiteknologi.co.id/goleaf/glauth/models"
	"git.solusiteknologi.co.id/goleaf/glautonumber"
	"git.solusiteknologi.co.id/goleaf/glcommon"
	"git.solusiteknologi.co.id/goleaf/glmail"
	"git.solusiteknologi.co.id/goleaf/glmail/model"
	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/glqueue/config"
	"git.solusiteknologi.co.id/goleaf/glrecaptcha"
	"git.solusiteknologi.co.id/goleaf/glregistration"
	"git.solusiteknologi.co.id/goleaf/glwebhook"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	satusehat "git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"

	_ "git.solusiteknologi.co.id/sts-satusehat/ssbackend/docs"
)

func ConfigureFiber(app *fiber.App) {
	apiPrefix := "/api"
	enableUserSchemaMapping := glutil.ToBool(glutil.GetEnv(constants.ENV_ENABLE_USER_SCHEMA_MAPPING))

	app.Get("/api/docs/*", swagger.New(swagger.Config{
		DeepLinking: true,
	})) // default
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/", fiber.StatusTemporaryRedirect)
	}) // default
	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/", fiber.StatusTemporaryRedirect)
	}) // default

	// auth module
	glauth.Setup(app, glauth.Config{
		PathPrefix:              apiPrefix,
		GroupAppConfig:          strings.Split(glutil.GetEnv(constants.ENV_APP_CONFIG_GROUPS, "APPS"), ","),
		EnableUserSchemaMapping: enableUserSchemaMapping,
		UserMaxSession:          int64(glutil.GetEnvInt(constants.ENV_MAX_USER_SESSION_LIMIT)),
		SessionExpired:          convertStringToDuration(glutil.GetEnv(constants.ENV_AUTO_LOGOUT_DURATION, "12h")),
		OnSuccessLogin: func(c *fiber.Ctx, _ goleafcore.Dto, _ *models.OutLogin) {
			c.Locals(authconst.LOCAL_PREFIX_KEYSESS, glutil.ToString(c.Get(authconst.HEADER_CLIENT_TYPE)))
		},
	})

	glqueue.Setup(config.Config{
		AutoFetchInterval: glutil.GetEnv(constants.ENV_CHECK_QUEUE_INTERVAL, "10s"),
	})

	glcommon.Setup(app, glcommon.Config{
		ApiPrefix:               apiPrefix,
		Middleware:              glauth.MiddlewareTask,
		MailConfig:              getMailConfig(),
		EnableUserSchemaMapping: enableUserSchemaMapping,
	})

	glautonumber.Setup(app, glautonumber.Config{
		ApiPrefix:  apiPrefix,
		Middleware: glauth.MiddlewareTask,
	})

	glwebhook.Setup(app, glwebhook.Config{
		ApiPrefix:  apiPrefix,
		Middleware: glauth.MiddlewareTask,
	})

	glrecaptcha.Setup(app, glrecaptcha.Config{
		RecaptchaSecretKey: glutil.GetEnv(constants.ENV_RECAPTCHA_SECRET_KEY),
		Middleware:         glauth.NewAnnonymousMiddleware(),
		ApiPrefix:          apiPrefix,
	})

	glregistration.Setup(app, glregistration.Config{
		ApiPrefix:               apiPrefix,
		Middleware:              glauth.NewAnnonymousMiddleware(),
		MailConfig:              getMailConfig(),
		EnableUserSchemaMapping: enableUserSchemaMapping,
	})

	satusehat.Setup(app, satusehat.Config{
		ApiPrefix:  apiPrefix,
		Middleware: glauth.MiddlewareTask,
	})
}

func convertStringToDuration(durationStr string, def ...time.Duration) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}

		return 5 * time.Minute
	}

	return duration
}

func getMailConfig() model.Config {
	return glmail.Config{
		RetryTimes: 5,
		SmtpConfig: glmail.SmtpConfig{
			Host:     glutil.GetEnv(constants.ENV_MAIL_HOST, "smtp.gmail.com"),
			Port:     glutil.GetEnvInt(constants.ENV_MAIL_PORT, 465),
			Username: glutil.GetEnv(constants.ENV_MAIL_USER, "supicadn"),
			Password: glutil.GetEnv(constants.ENV_MAIL_PASSWORD, "khujzppvisbanuph"),
			From:     glutil.GetEnv(constants.ENV_MAIL_FROM, "tmpl.id<supicadn@gmail.com>"),
		},
	}
}
