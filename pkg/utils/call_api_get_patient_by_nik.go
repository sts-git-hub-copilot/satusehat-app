package utils

import (
	"fmt"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sspatientdao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InputCallApiGetPatient struct {
	Tx        gldb.Tx
	Audit     gldata.AuditData
	PatientId int64
	Nik       string
}

func CallApiGetPatientByNik(input InputCallApiGetPatient) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat (GET Patient) -> ")

	identifierQuery := fmt.Sprintf("https://fhir.kemkes.go.id/id/nik|%s", input.Nik)

	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Patient",
		Method: fiber.MethodGet,
		Query: goleafcore.Dto{
			"identifier": identifierQuery,
		},
	})

	if err != nil {
		logrus.Errorf("Failed to call Satu Sehat API: %v", err)
		return nil, err
	}

	idss, err := getIdssPatient(&output, input.Nik)
	if err != nil {
		return nil, err
	}

	logrus.Debug("NIK : ", input.Nik)
	logrus.Debug("IDSS : ", idss)

	if glutil.IsEmpty(idss) {
		return nil, glerr.New("IDSS patient not found for NIK " + input.Nik)
	}

	_, err = sspatientdao.UpdateIdss(sspatientdao.InputUpdateIdss{
		Tx:        input.Tx,
		Audit:     input.Audit,
		PatientId: input.PatientId,
		Idss:      idss,
	})
	if err != nil {
		return nil, err
	}

	logrus.Println(" -- END Call API Satu Sehat (GET Patient) -> ")
	output.Put("idss", idss)

	return output, nil
}

type patientResponse struct {
	Entry []struct {
		Resource struct {
			Id string `json:"id"`
		} `json:"resource"`
	} `json:"entry"`
}

func getIdssPatient(dto *goleafcore.Dto, nik string) (string, error) {
	var response patientResponse
	err := dto.ToStruct(&response)
	if err != nil {
		return "", glerr.New("failed to parse patient response: ", err)
	}

	// Loop through entries to find matching name
	for _, entry := range response.Entry {
		return entry.Resource.Id, nil
	}

	return "", nil
}
