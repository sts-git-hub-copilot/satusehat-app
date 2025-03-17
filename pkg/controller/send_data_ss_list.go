package controller

import (
	"context"
	"errors"
	"time"

	"git.solusiteknologi.co.id/goleaf/glqueue"
	"git.solusiteknologi.co.id/goleaf/goleafcore"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glapi"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glconstant"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldata"
	"git.solusiteknologi.co.id/goleaf/goleafcore/gldb"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glerr"
	"git.solusiteknologi.co.id/goleaf/goleafcore/glutil"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/constants"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssanamnesadao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssdiagnosadao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/ssdoctordao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sslocationdao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sspatientdao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstrxanamnesadao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstrxdao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstrxdiagnosadao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/dao/sstrxitemdao"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/queueprocess"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/tables"
	"git.solusiteknologi.co.id/sts-satusehat/ssbackend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type BodySendDataSSList struct {
	Timestamp int64    `json:"timestamp" example:"321362135123" `
	DataList  []DataSS `json:"dataList" example:"[]" validate:"dive"`
}

type DataSS struct {
	NikDoctor      string `json:"nikDoctor" example:"1234567890" validate:"required"`
	DoctorName     string `json:"doctorName" example:"Dr. John Doe" validate:"required"`
	DataJsonDoctor string `json:"dataJsonDoctor" example:"{}" `

	NikPatient      string `json:"nikPatient" example:"1234567890" validate:"required"`
	PatientName     string `json:"patientName" example:"John Doe" validate:"required"`
	PatientAddress  string `json:"patientAddress" example:"Jl. Jend. Sudirman No. 1" validate:"required"`
	DataJsonPatient string `json:"dataJsonPatient" example:"{}"`

	LocationName     string `json:"locationName" example:"RSUD Dr. Soetomo" validate:"required"`
	LocationCode     string `json:"locationCode" example:"1000001" validate:"required"`
	LocationAddress  string `json:"locationAddress" example:"Jl. Jend. Sudirman No. 1" validate:"required"`
	DataJsonLocation string `json:"dataJsonLocation" example:"{}"`

	DatetimeIn      string `json:"datetimeIn" example:"20251202120000" validate:"required"`
	DatetimeProcess string `json:"datetimeProcess" example:"20251202120000" validate:"required"`
	DatetimeEnd     string `json:"datetimeEnd" example:"20251202120000" validate:"required"`

	DiagnosaList []tables.Diagnosa `json:"diagnosaList" validate:"dive,required"`
	AnamnesaList []tables.Anamnesa `json:"anamnesaList" validate:"dive"`
	DataJsonTrx  string            `json:"dataJsonTrx" example:"{}"`
}

type OutSendDataSSList struct {
	SkipList []goleafcore.Dto `json:"skipList"`
}

type ErrorVal struct {
	Row              string `json:"row"`
	ErrorDescription error  `json:"errorDescription"`
}

func SendDataSSList(fc *fiber.Ctx) error {
	timeOutUsed := 5 * time.Minute
	return glapi.ApiStdTx(fc, func(mt context.Context, audit *gldata.AuditData, tx gldb.Tx) interface{} {
		body := BodySendDataSSList{}
		out := OutSendDataSSList{}
		errValidationList := make([]ErrorVal, 0)
		skipList := make([]goleafcore.Dto, 0)
		err := glapi.FetchValidBody(fc, &body)
		if err != nil {
			if len(body.DataList) > 0 {
				for _, item := range body.DataList {
					validationField(item, &skipList, &errValidationList)
				}
			}
		}

		if len(skipList) > 0 {
			out.SkipList = skipList
			for _, item := range skipList {
				e := queueprocess.SendWebhookUpdateStatus(queueprocess.InputSendWebhookUpdateStatus{
					Mt:          mt,
					Audit:       *audit,
					TrxId:       glconstant.NULL_REF_VALUE_FOR_LONG,
					Status:      constants.STATUS_FAIL,
					FailMessage: err.Error(),
					DataJson:    item.GetString("dataJsonTrx"),
				})
				if e != nil {
					return e
				}
			}
			return out
		}

		if len(body.DataList) > 0 {
			for _, item := range body.DataList {
				errValidationList = make([]ErrorVal, 0)
				err := actionSend(
					InputActionSend{
						Mt:       mt,
						Tx:       tx,
						Item:     item,
						Audit:    *audit,
						SkipList: &skipList,
						ErrList:  &errValidationList,
					},
				)
				if err != nil {
					if len(errValidationList) > 0 {
						skipList = append(skipList, goleafcore.NewOrEmpty(struct {
							DataSS
							ErrList []ErrorVal
						}{
							DataSS:  item,
							ErrList: errValidationList,
						}))
					} else {
						e := queueprocess.SendWebhookUpdateStatus(queueprocess.InputSendWebhookUpdateStatus{
							Mt:          mt,
							Audit:       *audit,
							TrxId:       glconstant.NULL_REF_VALUE_FOR_LONG,
							Status:      constants.STATUS_FAIL,
							FailMessage: err.Error(),
							DataJson:    item.DataJsonTrx,
						})
						if e != nil {
							return e
						}
						return err
					}
				}
				if len(errValidationList) > 0 {
					skipList = append(skipList, goleafcore.NewOrEmpty(struct {
						DataSS
						ErrList []ErrorVal
					}{
						DataSS:  item,
						ErrList: errValidationList,
					}))
				}
			}
		}

		out.SkipList = skipList
		for _, item := range skipList {
			e := queueprocess.SendWebhookUpdateStatus(queueprocess.InputSendWebhookUpdateStatus{
				Mt:          mt,
				Audit:       *audit,
				TrxId:       glconstant.NULL_REF_VALUE_FOR_LONG,
				Status:      constants.STATUS_FAIL,
				FailMessage: err.Error(),
				DataJson:    item.GetString("dataJsonTrx"),
			})
			if e != nil {
				return e
			}
		}

		return out
	}, timeOutUsed)
}

type InputActionSend struct {
	Mt       context.Context
	Tx       gldb.Tx
	Item     DataSS
	Audit    gldata.AuditData
	SkipList *[]goleafcore.Dto
	ErrList  *[]ErrorVal
}

func actionSend(input InputActionSend) error {

	// set data doctor
	doctor, err := ssdoctordao.Set(ssdoctordao.InputSet{
		Tx:         input.Tx,
		Audit:      input.Audit,
		Nik:        input.Item.NikDoctor,
		DoctorName: input.Item.DoctorName,
		DataJson:   input.Item.DataJsonDoctor,
	})
	if err != nil {
		return err
	}

	if glutil.IsEmpty(doctor.Idss) {
		resp, err := utils.CallApiGetPractitionerByNik(utils.InputCallApiGetPractitioner{
			Tx:       input.Tx,
			Audit:    input.Audit,
			DoctorId: doctor.DoctorId,
			Nik:      doctor.Nik,
		})
		if err != nil {
			return err
		}

		doctor.Idss = resp.GetString("idss")
	}

	patient, err := sspatientdao.Set(sspatientdao.InputSet{
		Tx:          input.Tx,
		Audit:       input.Audit,
		Nik:         input.Item.NikPatient,
		PatientName: input.Item.PatientName,
		Address:     input.Item.PatientAddress,
		DataJson:    input.Item.DataJsonPatient,
	})
	if err != nil {
		return err
	}

	if glutil.IsEmpty(patient.Idss) {
		resp, err := utils.CallApiGetPatientByNik(utils.InputCallApiGetPatient{
			Tx:        input.Tx,
			Audit:     input.Audit,
			PatientId: patient.PatientId,
			Nik:       patient.Nik,
		})
		if err != nil {
			return err
		}

		patient.Idss = resp.GetString("idss")
	}

	location, err := sslocationdao.Set(sslocationdao.InputSet{
		Tx:           input.Tx,
		Audit:        input.Audit,
		LocationName: input.Item.LocationName,
		LocationCode: input.Item.LocationCode,
		Address:      input.Item.LocationAddress,
		DataJson:     input.Item.DataJsonLocation,
	})
	if err != nil {
		return err
	}

	if glutil.IsEmpty(location.Idss) {
		resp, err := utils.CallApiGetLocation(utils.InputCallApiGetLocation{
			Tx:              input.Tx,
			Audit:           input.Audit,
			LocationId:      location.LocationId,
			LocationName:    location.LocationName,
			LocationAddress: location.Address,
		})
		if err != nil {
			return err
		}

		location.Idss = resp.GetString("idss")
	}

	inputTrx := sstrxdao.InputAdd{
		Tx:        input.Tx,
		Audit:     input.Audit,
		PatientId: patient.PatientId,
		StatusDoc: constants.STATUS_DOC_IN_PROGRESS,
		DataJson:  input.Item.DataJsonTrx,
	}
	err = gldb.SelectRowQTx(input.Tx, *gldb.NewQBuilder().
		Add(` SELECT TO_CHAR(
						TO_TIMESTAMP(:datetimeIn, 'YYYYMMDDHH24MISS') AT TIME ZONE 'UTC',
						'YYYYMMDDHH24MISS'
					) datetime_in,
					 TO_CHAR(
						TO_TIMESTAMP(:datetimeProcess, 'YYYYMMDDHH24MISS') AT TIME ZONE 'UTC',
						'YYYYMMDDHH24MISS'
					) datetime_process,
					 TO_CHAR(
						TO_TIMESTAMP(:datetimeEnd, 'YYYYMMDDHH24MISS') AT TIME ZONE 'UTC',
						'YYYYMMDDHH24MISS'
					) datetime_out,
					(EXTRACT(EPOCH FROM (TO_TIMESTAMP(:datetimeEnd, 'YYYYMMDDHH24MISS') - 
							TO_TIMESTAMP(:datetimeIn, 'YYYYMMDDHH24MISS'))) / 60)::BIGINT AS min_duration
		`).
		SetParam("datetimeIn", input.Item.DatetimeIn).
		SetParam("datetimeProcess", input.Item.DatetimeProcess).
		SetParam("datetimeEnd", input.Item.DatetimeEnd), &inputTrx)
	if err != nil {
		return glerr.New("error convertion time : ", err)
	}

	existstrx, err := sstrxdao.IsExists(sstrxdao.InputIsExists{
		Tx:         input.Tx,
		Audit:      input.Audit,
		PatientId:  patient.PatientId,
		DatetimeIn: inputTrx.DatetimeIn,
	})
	if err != nil {
		return err
	}

	if existstrx {
		return glerr.New("Data transaksi sudah ada/duplikat untuk pasien " + patient.PatientName + " dan waktu " + inputTrx.DatetimeIn)
	}

	trx, err := sstrxdao.Add(inputTrx)
	if err != nil {
		return err
	}

	_, err = sstrxitemdao.Add(sstrxitemdao.InputAdd{
		Tx:         input.Tx,
		Audit:      input.Audit,
		TrxId:      trx.TrxId,
		LineNo:     1,
		DoctorId:   doctor.DoctorId,
		LocationId: location.LocationId,
	})
	if err != nil {
		return err
	}

	for _, anamnesa := range input.Item.AnamnesaList {
		exists, amn, err := ssanamnesadao.IsExists(ssanamnesadao.InputIsExists{
			Tx:           input.Tx,
			Audit:        input.Audit,
			AnamnesaCode: anamnesa.AnamnesaCode,
		})
		if err != nil {
			return err
		}
		if !exists {
			return glerr.New("Anamnesa is not found ", anamnesa.AnamnesaCode)
		}

		_, err = sstrxanamnesadao.Add(sstrxanamnesadao.InputAdd{
			Tx:         input.Tx,
			Audit:      input.Audit,
			TrxId:      trx.TrxId,
			AnamnesaId: amn.AnamnesaId,
		})
		if err != nil {
			return err
		}
	}

	for _, diagnosa := range input.Item.DiagnosaList {
		exists, diag, err := ssdiagnosadao.IsExists(ssdiagnosadao.InputIsExists{
			Tx:           input.Tx,
			Audit:        input.Audit,
			DiagnosaCode: diagnosa.DiagnosaCode,
		})
		if err != nil {
			return err
		}
		if !exists {
			return glerr.New("Diagnosa is not found ", diagnosa.DiagnosaCode)
		}

		_, err = sstrxdiagnosadao.Add(sstrxdiagnosadao.InputAdd{
			Tx:         input.Tx,
			Audit:      input.Audit,
			TrxId:      trx.TrxId,
			DiagnosaId: diag.DiagnosaId,
		})
		if err != nil {
			return err
		}
	}

	err = glqueue.Enqueue(glqueue.QueueData{
		Mt:          input.Mt,
		AuditData:   input.Audit,
		QueueName:   constants.QUEUE_NAME_CREATE_ENCOUNTER,
		ProcessName: constants.QUEUE_NAME_CREATE_ENCOUNTER,
		Data: goleafcore.NewOrEmpty(queueprocess.InputSendEncounterQueue{
			TrxId: trx.TrxId,
		}),
	})
	if err != nil {
		return errors.New("error queue create encounter : " + err.Error())
	}

	return nil
}

func validationField(item DataSS, skipList *[]goleafcore.Dto, errList *[]ErrorVal) {
	*errList = make([]ErrorVal, 0)
	if item.NikDoctor == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("NIK Dokter tidak boleh kosong"),
		})
	}
	if item.DoctorName == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Nama Dokter tidak boleh kosong"),
		})
	}
	if item.NikPatient == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("NIK Pasien tidak boleh kosong"),
		})
	}
	if item.PatientName == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Nama Pasien tidak boleh kosong"),
		})
	}
	if item.PatientAddress == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Alamat Pasien tidak boleh kosong"),
		})
	}
	if item.LocationName == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Nama Lokasi tidak boleh kosong"),
		})
	}
	if item.LocationCode == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Kode Lokasi tidak boleh kosong"),
		})
	}
	if item.LocationAddress == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Alamat Lokasi tidak boleh kosong"),
		})
	}
	if item.DatetimeIn == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Tanggal tidak boleh kosong"),
		})
	}
	if item.DatetimeProcess == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Tanggal Proses tidak boleh kosong"),
		})
	}
	if item.DatetimeEnd == glconstant.EMPTY_VALUE {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Tanggal Selesai tidak boleh kosong"),
		})
	}
	if len(item.DiagnosaList) > 0 {
		for i, d := range item.DiagnosaList {
			if d.DiagnosaCode == glconstant.EMPTY_VALUE || d.DiagnosaName == glconstant.EMPTY_VALUE {
				*errList = append(*errList, ErrorVal{
					Row:              "Diagnosa" + glutil.ToString(i+1),
					ErrorDescription: glerr.New("Kode atau Nama Diagnosa tidak boleh kosong"),
				})
			}
		}
	} else {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Diagnosa tidak boleh kosong"),
		})
	}
	if len(item.AnamnesaList) > 0 {
		for i, d := range item.AnamnesaList {
			if d.AnamnesaCode == glconstant.EMPTY_VALUE || d.AnamnesaName == glconstant.EMPTY_VALUE {
				*errList = append(*errList, ErrorVal{
					Row:              "Anamnesa" + glutil.ToString(i+1),
					ErrorDescription: glerr.New("Kode atau Nama Anamnesa tidak boleh kosong"),
				})
			}
		}
	} else {
		*errList = append(*errList, ErrorVal{
			ErrorDescription: glerr.New("Anamnesa tidak boleh kosong"),
		})
	}
	if len(*errList) > 0 {
		*skipList = append(*skipList, goleafcore.NewOrEmpty(struct {
			DataSS
			ErrList []ErrorVal
		}{
			DataSS:  item,
			ErrList: *errList,
		}))
	}
}
