package utils

import (
	"fmt"

	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssdoctordao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InputCallApiGetPractitioner struct {
	Tx       gldb.Tx
	Audit    gldata.AuditData
	DoctorId int64
	Nik      string
}

func CallApiGetPractitionerByNik(input InputCallApiGetPractitioner) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat (GET Practitioner) -> ")

	// Format identifier parameter according to FHIR spec
	identifierQuery := fmt.Sprintf("https://fhir.kemkes.go.id/id/nik|%s", input.Nik)

	// Call API with GET method
	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Practitioner",
		Method: fiber.MethodGet,
		Query: goleafcore.Dto{
			"identifier": identifierQuery,
		},
	})

	if err != nil {
		logrus.Errorf("Failed to call Satu Sehat API: %v", err)
		return nil, err
	}

	idss, err := getIdssPractitioner(&output, input.Nik)
	if err != nil {
		return nil, err
	}

	logrus.Debug("NIK : ", input.Nik)
	logrus.Debug("IDSS : ", idss)

	if glutil.IsEmpty(idss) {
		return nil, glerr.New("IDSS practitioner not found for NIK " + input.Nik)
	}

	_, err = ssdoctordao.UpdateIdss(ssdoctordao.InputUpdateIdss{
		Tx:       input.Tx,
		Audit:    input.Audit,
		DoctorId: input.DoctorId,
		Idss:     idss,
	})
	if err != nil {
		return nil, err
	}

	logrus.Println(" -- END Call API Satu Sehat (GET Practitioner) -> ")
	output.Put("idss", idss)

	return output, nil
}

type practitionerResponse struct {
	Entry []struct {
		Resource struct {
			Id string `json:"id"`
		} `json:"resource"`
	} `json:"entry"`
}

func getIdssPractitioner(dto *goleafcore.Dto, nik string) (string, error) {
	var response practitionerResponse
	err := dto.ToStruct(&response)
	if err != nil {
		return "", glerr.New("failed to parse practitioner response: ", err)
	}

	// Loop through entries to find matching name
	for _, entry := range response.Entry {
		return entry.Resource.Id, nil
	}

	return "", nil
}
