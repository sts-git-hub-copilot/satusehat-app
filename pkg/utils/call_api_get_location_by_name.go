package utils

import (
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sslocationdao"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// InputCallApiGetLocation menyusun parameter untuk pemanggilan API satusehat.
type InputCallApiGetLocation struct {
	Tx              gldb.Tx
	Audit           gldata.AuditData
	LocationId      int64
	LocationName    string
	LocationAddress string
}

// CallApiGetLocation melakukan pemanggilan API GET ke endpoint /Location?name={name}.
func CallApiGetLocation(input InputCallApiGetLocation) (goleafcore.Dto, error) {
	logrus.Println(" -- START Call API Satu Sehat (GET Location) -> ")

	// Memanggil API dengan method GET
	output, err := CallApiSS(InputCallApiSS{
		Tx:     input.Tx,
		Audit:  input.Audit,
		Path:   "/Location",
		Method: fiber.MethodGet,
		Query: goleafcore.Dto{
			"name": input.LocationName,
		},
	})

	if err != nil {
		logrus.Errorf("Gagal memanggil API Satu Sehat: %v", err)
		return nil, err
	}

	idss, err := getLocationIdByName(&output, input.LocationName, input.LocationAddress)
	if err != nil {
		return nil, err
	}

	if glutil.IsEmpty(idss) {
		res, err := CallApiCreateLoc(InputCallApiCreateLoc{
			Tx:         input.Tx,
			Audit:      input.Audit,
			LocationId: input.LocationId,
		})
		if err != nil {
			return nil, err
		}
		idss = res.GetString("id")
	}

	logrus.Debug("Location Name : ", input.LocationName)
	logrus.Debug("IDSS : ", idss)

	if glutil.IsEmpty(idss) {
		return nil, glerr.New("IDSS location not found for Location Name " + input.LocationName + " and Location Address " + input.LocationAddress)
	}

	_, err = sslocationdao.UpdateIdss(sslocationdao.InputUpdateIdss{
		Tx:         input.Tx,
		Audit:      input.Audit,
		LocationId: input.LocationId,
		Idss:       idss,
	})
	if err != nil {
		return nil, err
	}

	logrus.Println(" -- END Call API Satu Sehat (GET Location) -> ")
	output.Put("idss", idss)
	return output, nil
}

type locationResponse struct {
	Entry []struct {
		Resource struct {
			Id          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"resource"`
	} `json:"entry"`
}

func getLocationIdByName(dto *goleafcore.Dto, locationName string, locationAddress string) (string, error) {
	var response locationResponse
	err := dto.ToStruct(&response)
	if err != nil {
		return "", glerr.New("failed to parse location response: ", err)
	}

	// Loop through entries to find matching name
	for _, entry := range response.Entry {
		if entry.Resource.Name == locationName && entry.Resource.Description == locationAddress {
			return entry.Resource.Id, nil
		}
	}

	return "", nil
}
