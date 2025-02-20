package main

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore/glinit"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/app/config"
	"github.com/sirupsen/logrus"
)

// @title Solusi Satu Sehat API
// @version 1.0
// @description API documentation ssbackend
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://sts-satusehat
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type "Bearer" followed by a space and JWT token."

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// @tokenUrl 				  https://api.sts-satusehat/api/v1/auth/login
func main() {
	app := glinit.InitAll(glinit.LogConfig{
		LogFile: glutil.GetEnv(glinit.ENV_LOG_FILE, "./log/ssbackend.log"),
	}, glinit.DbConfig{
		Host:            glutil.GetEnv(glinit.ENV_DB_HOST, "172.17.0.1"),
		Port:            glutil.GetEnvInt(glinit.ENV_DB_PORT, 5432),
		Name:            glutil.GetEnv(glinit.ENV_DB_NAME, "erp_cloud"),
		User:            glutil.GetEnv(glinit.ENV_DB_USER, "sts"),
		Password:        glutil.GetEnv(glinit.ENV_DB_PASSWORD, "Awesome123!"),
		ApplicationName: "ssbackend",
	}, glinit.ServerConfig{
		Port:          glutil.GetEnvInt(glinit.ENV_SERVER_PORT, 5005),
		Multitenant:   true,
		EnableFluentd: true,
	})

	config.ConfigureFiber(app)

	logrus.Info("Starting server")

	glinit.StartServer()

}
