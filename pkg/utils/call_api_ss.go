package utils

import (
	"fmt"
	"time"

	"git.solusiteknologi.co.id/goleaf/glhttp"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glmt"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssresponselogdao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// InputCallApiSS menyusun parameter untuk pemanggilan API satusehat.
type InputCallApiSS struct {
	Tx     gldb.Tx
	Audit  gldata.AuditData
	Path   string
	Body   goleafcore.Dto
	Query  goleafcore.Dto
	Param  goleafcore.Dto
	Method string
}

// CallApiSS melakukan pemanggilan ke endpoint base url.
func CallApiSS(input InputCallApiSS) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat -> ")

	logrus.Debug("Tenant : ", glmt.GetCurrentTenant(gldb.FetchMtFromTrx(input.Tx)))

	configSS, err := fetchConfSS(input.Tx, input.Audit)
	if err != nil {
		return nil, err
	}

	logrus.Debug("CONFIG SS : ", configSS.PrettyString())

	token, err := fetchToken(input.Tx, configSS, input.Audit)
	if err != nil {
		return nil, err
	}

	mode := configSS.GetString(constants.PARAM_CODE_SS_MODE_ACCESS)
	baseUrl := configSS.GetString(constants.PARAM_CODE_SS_BASE_URL_SANDBOX)
	if mode == constants.MODE_PRODUCTION {
		baseUrl = configSS.GetString(constants.PARAM_CODE_SS_BASE_URL_PRODUCTION)
	}
	apiUrl := baseUrl + input.Path

	// Menyiapkan header request
	header := goleafcore.Dto{}
	header.Put("Authorization", fmt.Sprintf("Bearer %s", token))

	logrus.Debug("Request URL: ", apiUrl)
	logrus.Debug("Request Body: ", input.Body.PrettyString())
	logrus.Debug("Request Header: ", header.PrettyString())

	// Melakukan pemanggilan HTTP menggunakan helper glutil.CallHttp
	// Perhatikan: pada pemanggilan ini, kita mengirim body berupa string (bukan JSON)
	// Jika glutil.CallHttp hanya mendukung JSON, maka Anda perlu memodifikasi fungsi tersebut
	output, err := callApi(inputCallApi{
		FullUrl:  apiUrl,
		Query:    input.Query,
		Param:    input.Param,
		Header:   header,
		BodyJson: input.Body,
		Method:   input.Method,
	})
	if err != nil {
		logrus.Println("Output response call API ss : ", err.Error())
		_, errIn := ssresponselogdao.Add(ssresponselogdao.InputAdd{
			Tx:        input.Tx,
			Audit:     input.Audit,
			Url:       apiUrl,
			Method:    input.Method,
			BodyReq:   input.Body.ToJsonString(),
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
		Method:    input.Method,
		BodyReq:   input.Body.ToJsonString(),
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

	logrus.Println("Output response: ", status, " | ", outDto.PrettyString())
	return outDto, nil
}

type inputCallApi struct {
	FullUrl  string `validate:"required"`
	Query    goleafcore.Dto
	Param    goleafcore.Dto
	Header   goleafcore.Dto
	BodyJson goleafcore.Dto
	BodyRaw  []byte
	FilePath string
	Method   string
}

func callApi(input inputCallApi) (glutil.OutputCallHttp, error) {
	output := glutil.OutputCallHttp{}
	if input.Method == fiber.MethodPost || input.Method == fiber.MethodPut {
		resp, err := glutil.CallHttp(glutil.InputCallHttp{
			FullUrl:  input.FullUrl,
			Method:   input.Method,
			BodyJson: input.BodyJson,
			BodyRaw:  input.BodyRaw,
			Header:   input.Header,
		})
		if err != nil {
			return output, glerr.Wrap("failed call api", err)
		}
		output = resp
	}

	if input.Method == fiber.MethodGet {
		resp, err := glhttp.Get(glhttp.InputCallHttpGet{
			FullUrl: input.FullUrl,
			Query:   input.Query,
			Param:   input.Param,
			Header:  input.Header,
		})
		if err != nil {
			return output, glerr.Wrap("failed call api", err)
		}

		output = glutil.OutputCallHttp{
			IsOk:          resp.IsOk,
			ResponseCode:  resp.ResponseCode,
			ResponseBytes: resp.ResponseBytes,
			ProcessTime:   resp.ProcessTime,
		}
	}

	return output, nil
}
