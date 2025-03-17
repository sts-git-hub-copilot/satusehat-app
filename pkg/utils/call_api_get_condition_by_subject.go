package utils

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstrxidssendao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// InputCallApiGetConditionBySubject menyusun parameter untuk pemanggilan API satusehat.
type InputCallApiGetConditionBySubject struct {
	Tx             gldb.Tx
	Audit          gldata.AuditData
	TrxId          int64
	ConditionInput ConditionInput
}

// CallApiGetConditionBySubject melakukan pemanggilan API GET ke endpoint /Condition?subject={name}&encounter={idssencounter}.
func CallApiGetConditionBySubject(input InputCallApiGetConditionBySubject) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat (GET ConditionBySubject) -> ")

	// Memanggil API dengan method GET
	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Condition",
		Method: fiber.MethodGet,
		Query: goleafcore.Dto{
			"subject":   input.ConditionInput.IdssPatient,
			"encounter": input.ConditionInput.IdssEncounter,
		},
	})

	if err != nil {
		logrus.Errorf("Gagal memanggil API Satu Sehat: %v", err)
		return nil, err
	}

	idss, err := getConditionBySubject(&output, input)
	if err != nil {
		return nil, err
	}

	if glutil.IsEmpty(idss) {
		res, err := CallApiCreateCondition(InputCreateCondition(input))
		if err != nil {
			return nil, err
		}
		idss = res.GetString("id")
	}

	logrus.Debug("Condition Input : ", goleafcore.NewOrEmpty(input.ConditionInput).ToJsonString())
	logrus.Debug("IDSS Condition : ", idss)

	_, err = sstrxidssendao.AddIfNotExists(sstrxidssendao.InputAdd{
		Tx:     input.Tx,
		Audit:  input.Audit,
		TrxId:  input.TrxId,
		Idssen: idss,
		Status: constants.COMBO_VALUE_CONDITION,
	})
	if err != nil {
		return nil, err
	}

	logrus.Println(" -- END Call API Satu Sehat (GET Location) -> ")
	output.Put("idss", idss)
	return output, nil
}

type getConditionResponse struct {
	Entry []struct {
		Resource struct {
			Id        string `json:"id"`
			Encounter struct {
				Reference string `json:"reference"`
			} `json:"encounter"`
			Subject struct {
				Reference string `json:"reference"`
			} `json:"subject"`
		} `json:"resource"`
	} `json:"entry"`
}

func getConditionBySubject(dto *goleafcore.Dto, input InputCallApiGetConditionBySubject) (string, error) {
	var response getConditionResponse
	err := dto.ToStruct(&response)
	if err != nil {
		return "", glerr.New("failed to parse location response: ", err)
	}

	// Loop through entries to find matching name
	for _, entry := range response.Entry {
		if entry.Resource.Encounter.Reference == "Encounter/"+input.ConditionInput.IdssEncounter &&
			entry.Resource.Subject.Reference == "Patient/"+input.ConditionInput.IdssPatient {
			return entry.Resource.Id, nil
		}
	}

	return "", nil
}
