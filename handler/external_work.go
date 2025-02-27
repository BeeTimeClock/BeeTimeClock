package handler

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"slices"
	"strconv"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/microsoft"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
)

type ExternalWork struct {
	env          *core.Environment
	user         *repository.User
	externalWork *repository.ExternalWork
	holiday      *repository.Holiday
}

func NewExternalWork(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork, holiday *repository.Holiday) *ExternalWork {
	return &ExternalWork{
		env:          env,
		user:         user,
		externalWork: externalWork,
		holiday:      holiday,
	}
}

func (h *ExternalWork) ExternalWorkGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	externalWorkItems, err := h.externalWork.ExternalWorkFindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkItems))
}

func (h *ExternalWork) ExternalWorkGetInvoiced(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	externalWorkItems, err := h.externalWork.ExternalWorkFindByUserIDAndStatus(user.ID, model.EXTERNAL_WORK_STATUS_INVOICED)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	type ExternalWorkInvoicedInfo struct {
		InvoiceDate time.Time
		Identifier  string
	}

	invoiced := make(map[string]ExternalWorkInvoicedInfo)
	invoicedMapped := []ExternalWorkInvoicedInfo{}

	for _, work := range externalWorkItems {
		invoiceIdentifier := work.InvoiceIdentifier.String()

		if _, exists := invoiced[invoiceIdentifier]; !exists {
			info := ExternalWorkInvoicedInfo{
				InvoiceDate: *work.InvoiceDate,
				Identifier:  invoiceIdentifier,
			}
			invoiced[invoiceIdentifier] = info
			invoicedMapped = append(invoicedMapped, info)
		}
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(invoicedMapped))
}

func (h *ExternalWork) ExternalWorkGetById(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	fromDay := helper.GetDayDate(externalWorkItem.From)
	tillDay := helper.GetDayDate(externalWorkItem.Till)

	currentDay := fromDay
	for currentDay.Before(tillDay.AddDate(0, 0, 1)) {
		if !slices.ContainsFunc(externalWorkItem.WorkExpanses, func(n model.ExternalWorkExpense) bool {
			return n.Date.Round(24*time.Hour).UTC() == currentDay
		}) {
			expanseItem := model.ExternalWorkExpense{
				ExternalWork:   externalWorkItem,
				ExternalWorkID: externalWorkItem.ID,
				Date:           currentDay,
			}

			externalWorkItem.WorkExpanses = append(externalWorkItem.WorkExpanses, expanseItem)
		}

		currentDay = currentDay.AddDate(0, 0, 1)
	}

	holidays, err := h.holiday.HolidayFindByDateRange(fromDay, tillDay)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkItem.Calculate(holidays)))
}

func (h *ExternalWork) ExternalWorkCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}
	var externalWorkCreateRequest model.ExternalWorkCreateRequest

	err = c.BindJSON(&externalWorkCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWork := model.ExternalWork{
		User:                       user,
		ExternalWorkCompensationID: externalWorkCreateRequest.ExternalWorkCompensationID,
		From:                       externalWorkCreateRequest.From,
		Till:                       externalWorkCreateRequest.Till,
		Description:                externalWorkCreateRequest.Description,
		Identifier:                 uuid.New(),
	}

	err = h.externalWork.ExternalWorkInsert(&externalWork)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if microsoft.IsMicrosoftConnected() {
		eventId, err := microsoft.CreateCalendarEntryFromExternalWork(user.Username, &externalWork)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}

		externalWorkExternalEvent := model.ExternalWorkExternalEvent{
			ExternalWork:          externalWork,
			ExternalEventID:       eventId,
			ExternalEventProvider: model.EXTERNAL_EVENT_PROVIDER_MICROSOFT,
		}

		err = h.externalWork.ExternalWorkExternalEventInsert(&externalWorkExternalEvent)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(externalWork))
}

func (h *ExternalWork) ExternalWorkExpanseCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var externalWorkExpenseCreateRequest model.ExternalWorkExpenseCreateRequest
	err = c.BindJSON(&externalWorkExpenseCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	if !externalWorkItem.IsEditable() {
		c.JSON(http.StatusForbidden, model.NewErrorResponse(errors.New("can't change in this state")))
		return
	}

	externalWorkExpense := model.ExternalWorkExpense{
		ExternalWork:           externalWorkItem,
		Date:                   externalWorkExpenseCreateRequest.Date,
		DepartureTime:          externalWorkExpenseCreateRequest.DepartureTime,
		ArrivalTime:            externalWorkExpenseCreateRequest.ArrivalTime,
		TravelDurationHours:    externalWorkExpenseCreateRequest.TravelDurationHours,
		PauseDurationHours:     externalWorkExpenseCreateRequest.PauseDurationHours,
		OnSiteFrom:             externalWorkExpenseCreateRequest.OnSiteFrom,
		OnSiteTill:             externalWorkExpenseCreateRequest.OnSiteTill,
		Place:                  externalWorkExpenseCreateRequest.Place,
		TravelWithPrivateCarKm: externalWorkExpenseCreateRequest.TravelWithPrivateCarKm,
		AdditionalOptions:      externalWorkExpenseCreateRequest.AdditionalOptions,
	}

	err = h.externalWork.ExternalWorkExpenseInsert(&externalWorkExpense)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, externalWorkExpense)
}

func (h *ExternalWork) ExternalWorkExpanseUpdate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	expanseIdParam := c.Param("externalWorkExpanseId")
	expanseId, err := strconv.ParseInt(expanseIdParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var externalWorkExpenseUpdateRequest model.ExternalWorkExpenseUpdateRequest
	err = c.BindJSON(&externalWorkExpenseUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	if !externalWorkItem.IsEditable() {
		c.JSON(http.StatusForbidden, model.NewErrorResponse(errors.New("can't change in this state")))
		return
	}

	externalWorkExpenseItem, err := h.externalWork.ExternalWorkExpenseFindById(uint(expanseId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	externalWorkExpenseItem.DepartureTime = externalWorkExpenseUpdateRequest.DepartureTime
	externalWorkExpenseItem.ArrivalTime = externalWorkExpenseUpdateRequest.ArrivalTime
	externalWorkExpenseItem.TravelDurationHours = externalWorkExpenseUpdateRequest.TravelDurationHours
	externalWorkExpenseItem.PauseDurationHours = externalWorkExpenseUpdateRequest.PauseDurationHours
	externalWorkExpenseItem.OnSiteFrom = externalWorkExpenseUpdateRequest.OnSiteFrom
	externalWorkExpenseItem.OnSiteTill = externalWorkExpenseUpdateRequest.OnSiteTill
	externalWorkExpenseItem.Place = externalWorkExpenseUpdateRequest.Place
	externalWorkExpenseItem.TravelWithPrivateCarKm = externalWorkExpenseUpdateRequest.TravelWithPrivateCarKm
	externalWorkExpenseItem.AdditionalOptions = externalWorkExpenseUpdateRequest.AdditionalOptions

	err = h.externalWork.ExternalWorkExpenseUpdate(&externalWorkExpenseItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, externalWorkExpenseItem)
}

func (h *ExternalWork) ExternalWorkSubmit(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	externalWorkItem.Status = model.EXTERNAL_WORK_STATUS_ACCEPTED

	err = h.externalWork.ExternalWorkUpdate(&externalWorkItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, externalWorkItem)
}

func (h *ExternalWork) AdministrationExternalWorkCompensationGetAll(c *gin.Context) {
	compensations, err := h.externalWork.ExternalWorkCompensationFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(compensations))
}

func (h *ExternalWork) AdministrationExternalWorkCompensationCreate(c *gin.Context) {
	var externalWorkCompensation model.ExternalWorkCompensation

	err := c.BindJSON(&externalWorkCompensation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	err = h.externalWork.ExternalWorkCompensationInsert(&externalWorkCompensation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkCompensation))
}

func (h *ExternalWork) AdministrationExternalWorkCompensationUpdate(c *gin.Context) {
	externalWorkCompensationIdParam := c.Param("externalWorkCompensationId")
	externalWorkCompensationId, err := strconv.ParseInt(externalWorkCompensationIdParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var externalWorkCompensationUpdate model.ExternalWorkCompensation

	err = c.BindJSON(&externalWorkCompensationUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkCompensation, err := h.externalWork.ExternalWorkCompensationFindById(uint(externalWorkCompensationId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	externalWorkCompensation.AdditionalOptions = externalWorkCompensationUpdate.AdditionalOptions
	externalWorkCompensation.PrivateCarKmCompensation = externalWorkCompensationUpdate.PrivateCarKmCompensation
	externalWorkCompensation.WithSocialInsuranceSlots = externalWorkCompensationUpdate.WithSocialInsuranceSlots
	externalWorkCompensation.WithoutSocialInsuranceSlots = externalWorkCompensationUpdate.WithoutSocialInsuranceSlots
	externalWorkCompensation.ValidFrom = externalWorkCompensationUpdate.ValidFrom
	externalWorkCompensation.ValidTill = externalWorkCompensationUpdate.ValidTill

	err = h.externalWork.ExternalWorkCompensationUpdate(&externalWorkCompensation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkCompensation))
}

func (h *ExternalWork) ExternalWorkDownloadPdf(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	invoiceIdentifierParam := c.Param("invoiceIdentifier")
	invoiceIdentifier, err := uuid.Parse(invoiceIdentifierParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	fileName := fmt.Sprintf("external_work_%d_%s.pdf", user.ID, invoiceIdentifier.String())

	fmt.Println(fileName)
	exists, err := helper.ExistsFile(h.env, fileName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		return
	}

	pdfFile, stat, err := helper.GetFile(h.env, fileName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`inline; filename="%s"`, fileName),
	}

	c.DataFromReader(http.StatusOK, stat.Size, "application/pdf", pdfFile, extraHeaders)
}

func (h *ExternalWork) ExternalWorkExportPdf(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	externalWorkItems, err := h.externalWork.ExternalWorkFindByUserIDAndStatus(uint(user.ID), model.EXTERNAL_WORK_STATUS_ACCEPTED)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if len(externalWorkItems) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	minStartDate := externalWorkItems[0].From
	maxEndDate := externalWorkItems[0].Till
	for _, item := range externalWorkItems {
		if item.Till.After(maxEndDate) {
			maxEndDate = item.Till
		}

		if item.From.Before(minStartDate) {
			minStartDate = item.From
		}
	}

	fromDay := helper.GetDayDate(minStartDate)
	tillDay := helper.GetDayDate(maxEndDate)

	holidays, err := h.holiday.HolidayFindByDateRange(fromDay, tillDay)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	pdf := fpdf.New(fpdf.OrientationLandscape, "mm", "A4", "")
	header := []string{
		"Datum",
		"Ort",
		"Nutzung priv. PKW",
		"VMA St./soz. frei",
		"VMA pauschal",
		"Vergütung Überstunden",
	}
	w := []float64{40.0, 50.0, 30.0, 30.0, 30.0, 30.0}

	for _, key := range externalWorkItems[0].ExternalWorkCompensation.AdditionalOptions.Keys() {
		header = append(header, key)
		w = append(w, 30.0)
	}

	lineHeight := 5.0

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	decoder := charmap.Windows1252.NewDecoder()
	euroCP1252 := []byte{0x80}
	euroUTF8, err := decoder.Bytes(euroCP1252)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	euro := func() string {
		return string(euroUTF8)
	}

	pdf.SetTopMargin(30)
	pdf.SetFont("Arial", "", 12)
	improvedTable := func() {
		left := 10.0
		pdf.SetX(left)
		getMaxLines := func(cells []string) int {
			max := 1
			for i, txt := range cells {
				lines := pdf.SplitLines([]byte(txt), w[i])
				if len(lines) > max {
					max = len(lines)
				}
			}
			return max
		}

		drawRow := func(cells []string, align string, style string) {
			pdf.SetFont("Arial", style, 12)
			maxLines := getMaxLines(cells)
			rowHeight := float64(maxLines) * lineHeight
			y := pdf.GetY()
			x := pdf.GetX()

			for i, txt := range cells {
				pdf.Rect(x, y, w[i], rowHeight, "")
				pdf.MultiCell(w[i], lineHeight, tr(txt), "", align, false)
				x += w[i]
				pdf.SetXY(x, y)
			}
			pdf.Ln(rowHeight)
		}

		drawRow(header, "C", "B")

		currency := func(value float64) string {
			return fmt.Sprintf("%.2f %s", value, euro())
		}

		totalExpensesWithSocialInsurance := 0.0
		totalExpensesWithoutSocialInsurance := 0.0
		optionsSums := make(map[string]float64)
		for _, externalWorkItem := range externalWorkItems {
			calculated := externalWorkItem.Calculate(holidays)
			for _, e := range calculated.WorkExpansesCalculated {
				rowData := []string{
					e.Date.Format("02.01.2006"),
					e.Place,
					currency(e.TravelPrivateKmCosts),
					currency(e.ExpensesWithoutSocialInsurance),
					currency(e.ExpensesWithSocialInsurance),
					"",
				}
				for _, key := range e.ExternalWork.ExternalWorkCompensation.AdditionalOptions.Keys() {
					if value, exists := e.AdditionalOptionsUsed[key]; exists {
						rowData = append(rowData, currency(value))
					} else {
						rowData = append(rowData, "")
					}
				}

				drawRow(rowData, "R", "")
			}

			totalExpensesWithSocialInsurance += calculated.TotalExpensesWithSocialInsurance
			totalExpensesWithoutSocialInsurance += calculated.TotalExpensesWithoutSocialInsurance
			for _, key := range calculated.ExternalWorkCompensation.AdditionalOptions.Keys() {
				optionsSums[key] += calculated.TotalOptions[key]
			}
		}
		sums := []string{
			"",
			"Summe",
			"",
			currency(totalExpensesWithoutSocialInsurance),
			currency(totalExpensesWithSocialInsurance),
			"",
		}

		for _, value := range optionsSums {
			sums = append(sums, currency(value))
		}

		drawRow(sums, "R", "B")

	}

	pdf.SetHeaderFuncMode(func() {
		logoName := "logo.png"

		exists, _ := helper.ExistsFile(h.env, logoName)

		if exists {
			file, _, _ := helper.GetFile(h.env, logoName)

			tmpFile, _ := os.CreateTemp("", "*.png")
			defer tmpFile.Close()

			io.Copy(tmpFile, file)

			pdf.ImageOptions(tmpFile.Name(), 5, 5, 70, 20, false, fpdf.ImageOptions{
				ImageType:             "",
				ReadDpi:               false,
				AllowNegativePosition: false,
			}, 0, "")
		}
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 5, "")
		pdf.CellFormat(30, 10, "Aufstellung zur Spesenabrechnung", "", 0, "L", false, 0, "")
		pdf.Ln(7)
		pdf.SetFont("Arial", "", 15)
		pdf.Cell(80, 5, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("Name: %s", user.FullName()), "", 0, "L", false, 0, "")
		pdf.Ln(6)
		pdf.Cell(80, 5, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("Lohnpersonal Nr: %d", user.StaffNumber), "", 0, "L", false, 0, "")
		pdf.Cell(70, 5, "")
		pdf.CellFormat(90, 10, fmt.Sprintf("Zeitraum: %s - %s", fromDay.Format("02.01.2006"), tillDay.Format("02.01.2006")), "1", 0, "R", false, 0, "")
		pdf.Ln(80)
	}, true)

	pdf.AliasNbPages("")
	pdf.AddPage()

	improvedTable()

	pdf.CellFormat(90, 10, fmt.Sprintf("Generiert von: %s", user.FullName()), "", 0, "L", false, 0, "")
	pdf.Ln(8)
	pdf.CellFormat(90, 10, fmt.Sprintf("Generiert am: %s", time.Now().Format("02.01.2006 15:04:05")), "", 0, "L", false, 0, "")

	invoiceUuid := uuid.New()
	fileName := fmt.Sprintf("external_work_%d_%s.pdf", user.ID, invoiceUuid)

	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}
	defer tmpFile.Close()

	err = pdf.Output(tmpFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}
	fileMeta := multipart.FileHeader{
		Filename: fileName,
		Header:   textproto.MIMEHeader{},
		Size:     stat.Size(),
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = helper.SaveFile(h.env, fileName, tmpFile, &fileMeta)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	for _, item := range externalWorkItems {
		now := time.Now()
		item.Status = model.EXTERNAL_WORK_STATUS_INVOICED
		item.InvoiceDate = &now
		item.InvoiceIdentifier = &invoiceUuid

		err = h.externalWork.ExternalWorkUpdate(&item)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	}

	tmpFile.Seek(0, 0)
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`inline; filename="%s"`, fileName),
	}

	c.DataFromReader(http.StatusOK, stat.Size(), "application/pdf", tmpFile, extraHeaders)
}
