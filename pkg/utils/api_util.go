package utils

import (
	"strconv"
	"time"

	"git.solusiteknologi.co.id/goleaf/glcommon/dao/sysconfigdao"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glconstant"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
)

func fetchConfSS(trx gldb.Tx, audit gldata.AuditData) (goleafcore.Dto, error) {
	res, err := sysconfigdao.GetValByParamGroup(trx, constants.PARAM_GROUP_SATUSEHAT, audit)
	if err != nil {
		return nil, glerr.New("parameter.group.not.found", constants.PARAM_GROUP_SATUSEHAT)
	}

	return res, nil
}

func fetchToken(trx gldb.Tx, configSS goleafcore.Dto, audit gldata.AuditData) (string, error) {
	var (
		stop       = false
		result     = glconstant.EMPTY_VALUE
		retryCount = 0
		maxCount   = 3
	)

	for !stop {
		output, err := CallAccessToken(InputCallAccessToken{Tx: trx, Audit: audit, ConfigSS: configSS})
		if err != nil || glutil.IsEmpty(output.GetString("access_token")) {
			retryCount = retryCount + 1
			if retryCount > maxCount {
				stop = true
				return result, glerr.New("failed.to.get.access.token.check.your.config", output)
			}
		} else {
			stop = true
		}

		result = output.GetString("access_token")
	}

	return result, nil
}

func ParseMillisToTime(epochMillisStr string, def ...time.Time) time.Time {
	// Konversi string ke int64
	epochMillis, err := strconv.ParseInt(epochMillisStr, 10, 64)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return time.Now()
	}

	// Konversi milidetik ke time.Time
	timestamp := time.Unix(epochMillis/1000, (epochMillis%1000)*int64(time.Millisecond))

	return timestamp
}
