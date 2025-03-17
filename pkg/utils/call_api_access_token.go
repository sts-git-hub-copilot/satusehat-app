package utils

import (
	"net/url"
	"time"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glconstant"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glmt"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssresponselogdao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstokendao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// InputCallAccessToken menyusun parameter untuk pemanggilan API access token.
type InputCallAccessToken struct {
	Tx       gldb.Tx
	Audit    gldata.AuditData
	ConfigSS goleafcore.Dto
}

// CallAccessToken melakukan pemanggilan ke endpoint access token dengan body form-urlencoded.
func CallAccessToken(input InputCallAccessToken) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Access Token -> ")

	logrus.Debug("Tenant : ", glmt.GetCurrentTenant(gldb.FetchMtFromTrx(input.Tx)))

	token := glconstant.EMPTY_VALUE
	err := gldb.SelectRowQTx(input.Tx, *gldb.NewQBuilder().
		Add(`SELECT fss_get_access_token(:datetime) `).
		SetParam("datetime", input.Audit.Datetime()), &token)
	if err != nil {
		return nil, glerr.New("error on query get token : ", err)
	}
	if !glutil.IsEmpty(token) {
		return goleafcore.Dto{
			"access_token": token,
		}, nil
	}

	if input.ConfigSS == nil {
		configSS, err := fetchConfSS(input.Tx, input.Audit)
		if err != nil {
			return nil, err
		}
		input.ConfigSS = configSS
	}

	// logrus.Debug("CONFIG SS : ", input.ConfigSS.PrettyString())

	mode := input.ConfigSS.GetString(constants.PARAM_CODE_SS_MODE_ACCESS)
	clientId := input.ConfigSS.GetString(constants.PARAM_CODE_SS_CLIENT_ID)
	clientSecret := input.ConfigSS.GetString(constants.PARAM_CODE_SS_CLIENT_SECRET)
	urlAuth := input.ConfigSS.GetString(constants.PARAM_CODE_SS_AUTH_URL_SANDBOX)
	if mode == constants.MODE_PRODUCTION {
		urlAuth = input.ConfigSS.GetString(constants.PARAM_CODE_SS_AUTH_URL_PRODUCTION)
	}
	apiUrl := urlAuth + "/accesstoken?grant_type=client_credentials"

	// Membuat body form-urlencoded
	formData := url.Values{}
	formData.Set("client_id", clientId)
	formData.Set("client_secret", clientSecret)
	bodyStr := formData.Encode() // Hasil: "client_id=xxx&client_secret=yyy"

	// Menyiapkan header request
	header := goleafcore.Dto{}
	// Untuk pemanggilan access token, Content-Type harus "application/x-www-form-urlencoded"
	header.Put("Content-Type", "application/x-www-form-urlencoded")

	logrus.Debug("Request URL: ", apiUrl)
	logrus.Debug("Request Body: ", bodyStr)
	logrus.Debug("Request Header: ", header.PrettyString())

	// Melakukan pemanggilan HTTP menggunakan helper glutil.CallHttp
	// Perhatikan: pada pemanggilan ini, kita mengirim body berupa string (bukan JSON)
	// Jika glutil.CallHttp hanya mendukung JSON, maka Anda perlu memodifikasi fungsi tersebut
	output, err := glutil.CallHttp(glutil.InputCallHttp{
		FullUrl: apiUrl,
		Method:  fiber.MethodPost,
		// Mengirim body raw berupa string form-urlencoded
		BodyRaw: []byte(bodyStr),
		Header:  header,
	})
	if err != nil {
		logrus.Println("Output response call API access token : ", err.Error())
		_, errIn := ssresponselogdao.Add(ssresponselogdao.InputAdd{
			Tx:        input.Tx,
			Audit:     input.Audit,
			Url:       apiUrl,
			Method:    fiber.MethodPost,
			BodyReq:   bodyStr,
			HeaderReq: header.ToJsonString(),
			Status:    "E",
			RequestResponse: goleafcore.NewOrEmpty(goleafcore.Dto{
				"fullUrl": apiUrl,
				"header":  header,
				"body":    err.Error(),
			}).ToJsonString(),
			ResponseCode:     "401",
			ResponseBody:     err.Error(),
			LatencyMs:        output.ProcessTime.Milliseconds(),
			ResponseDatetime: glutil.DateFrom(time.Now().Add(output.ProcessTime)),
		})
		if errIn != nil {
			return nil, errIn
		}
		return nil, err
	}

	// Parsing respons ke dalam DTO
	outDto := goleafcore.Dto{}
	output.ToStruct(&outDto)

	status := "S"
	if !output.IsOk {
		status = "E"
	}

	_, err = ssresponselogdao.Add(ssresponselogdao.InputAdd{
		Tx:        input.Tx,
		Audit:     input.Audit,
		Url:       apiUrl,
		Method:    fiber.MethodPost,
		BodyReq:   bodyStr,
		HeaderReq: header.ToJsonString(),
		Status:    status,
		RequestResponse: goleafcore.Dto{
			"fullUrl": apiUrl,
			"header":  header,
			"body":    outDto.ToJsonString(),
		}.ToJsonString(),
		ResponseCode:     glutil.ToString(output.ResponseCode),
		ResponseBody:     outDto.ToJsonString(),
		LatencyMs:        output.ProcessTime.Milliseconds(),
		ResponseDatetime: glutil.DateFrom(time.Now().Add(output.ProcessTime)),
	})
	if err != nil {
		return nil, err
	}

	_, err = sstokendao.Add(sstokendao.InputAddToken{
		Tx:          input.Tx,
		Audit:       input.Audit,
		AccessToken: outDto.GetString("access_token"),
		IssuedAt:    ParseMillisToTime(outDto.GetString("issued_at")),
		ExpiresIn:   outDto.GetInt64("expires_in"),
		DataJson:    outDto.ToJsonString(),
	})
	if err != nil {
		return nil, err
	}

	logrus.Println("Output response: ", status, " | ", outDto.PrettyString())
	return outDto, nil
}

// time=2025-02-25T14:02:18+07:00 level=debug msg=OUTPUT : {
//   "access_token": "LfGcAhfchRTfeveEXxiYY6w8h5oA",
//   "api_product_list": "[api-satusehat-prod]",
//   "api_product_list_json": [
//     "api-satusehat-prod"
//   ],
//   "application_name": "42dfe7b9-99d1-4879-a0d6-96c3f0ce9782",
//   "client_id": "Ji93jAdSq8EVCO9QvGt9cTr1S5PAs7VEkuhkYDd7SHGFAelq",
//   "developer.email": "manikaaestheticc@gmail.com",
//   "expires_in": "14399",
//   "issued_at": "1740466938331",
//   "organization_name": "ihs-prod-1",
//   "refresh_count": "0",
//   "refresh_token_expires_in": "0",
//   "scope": "",
//   "status": "approved",
//   "token_type": "BearerToken"
// }
