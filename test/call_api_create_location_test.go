package test_test

import (
	"testing"
	"time"

	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gltest"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func TestCallApiCreateLocation(t *testing.T) {
	gltest.TestApi(t, func(app *fiber.App, tx gldb.Tx) error {

		return nil
	}, func(assert *gltest.Assert, app *fiber.App, tx gldb.Tx, i int) interface{} {
		err := gldb.ExecQTx(tx, *gldb.NewQBuilder().Add(`SET SEARCH_PATH TO demo`))
		if err != nil {
			return err
		}

		resp, err := utils.CallApiCreateLoc(utils.InputCallApiCreateLoc{
			Tx: tx,
			Audit: gldata.AuditData{
				UserLoginId:   10,
				TenantLoginId: 10,
				RoleLoginId:   10,
				SessionId:     glutil.RandomUuid(),
				Timestamp:     time.Now(),
			},
		})
		if err != nil {
			return err
		}

		logrus.Debug("OUTPUT : ", resp.PrettyString())

		return nil
	})
}
