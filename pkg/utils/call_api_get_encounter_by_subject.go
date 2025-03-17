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

// InputCallApiGetEncounterBySubject menyusun parameter untuk pemanggilan API satusehat.
type InputCallApiGetEncounterBySubject struct {
	Tx             gldb.Tx
	Audit          gldata.AuditData
	TrxId          int64
	EncounterInput EncounterInput
}

// CallApiGetEncounterBySubject melakukan pemanggilan API GET ke endpoint /Encounter?subject={name}.
func CallApiGetEncounterBySubject(input InputCallApiGetEncounterBySubject) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat (GET EncounterBySubject) -> ")

	// Memanggil API dengan method GET
	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Encounter",
		Method: fiber.MethodGet,
		Query: goleafcore.Dto{
			"subject": input.EncounterInput.IdssPatient,
		},
	})

	if err != nil {
		logrus.Errorf("Gagal memanggil API Satu Sehat: %v", err)
		return nil, err
	}

	idss, err := getEncounterBySubject(&output, input)
	if err != nil {
		return nil, err
	}

	if glutil.IsEmpty(idss) {
		res, err := CallApiCreateEncounter(InputCreateEncounter(input))
		if err != nil {
			return nil, err
		}
		idss = res.GetString("id")
	}

	logrus.Debug("Encounter Input : ", goleafcore.NewOrEmpty(input.EncounterInput).ToJsonString())
	logrus.Debug("IDSS Encounter : ", idss)

	_, err = sstrxidssendao.AddIfNotExists(sstrxidssendao.InputAdd{
		Tx:     input.Tx,
		Audit:  input.Audit,
		TrxId:  input.TrxId,
		Idssen: idss,
		Status: constants.COMBO_VALUE_ENCOUNTER,
	})
	if err != nil {
		return nil, err
	}

	logrus.Println(" -- END Call API Satu Sehat (GET Location) -> ")
	output.Put("idss", idss)
	return output, nil
}

type getEncounterResponse struct {
	Entry []struct {
		Resource struct {
			Id         string `json:"id"`
			Idendifier []struct {
				Value string `json:"value"`
			} `json:"identifier"`
		} `json:"resource"`
	} `json:"entry"`
}

func getEncounterBySubject(dto *goleafcore.Dto, input InputCallApiGetEncounterBySubject) (string, error) {
	var response getEncounterResponse
	err := dto.ToStruct(&response)
	if err != nil {
		return "", glerr.New("failed to parse location response: ", err)
	}

	// Loop through entries to find matching name
	for _, entry := range response.Entry {
		for _, identifier := range entry.Resource.Idendifier {
			if identifier.Value == input.EncounterInput.IdssPatient+"."+input.EncounterInput.StartTime {
				return entry.Resource.Id, nil
			}
		}
	}

	return "", nil
}
