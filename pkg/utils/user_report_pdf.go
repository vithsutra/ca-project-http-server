package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/signintech/gopdf"
	"github.com/vithsutra/ca_project_http_server/internals/models"
)

const LOGO_IMAGE_PATH = "./assets/vithsutra_logo.png"

func OuterBorderSection(pdf *gopdf.GoPdf) {
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.05)
	pdf.Line(1, 1, 20, 1)
	pdf.Line(1, 1, 1, 28.7)
	pdf.Line(1, 28.7, 20, 28.7)
	pdf.Line(20, 1, 20, 28.7)
}

func HeaderSection(pdf *gopdf.GoPdf) error {

	OuterBorderSection(pdf)

	x := 1.8
	y := 1.6
	pdf.Image(LOGO_IMAGE_PATH, x, y, &gopdf.Rect{
		H: 2,
		W: 3.5,
	})

	if err := pdf.SetFont("bold-font", "", 17); err != nil {
		return err
	}

	heading1 := "VITHSUTRA TECHNOLOGIES PVT. LTD."
	phone_number := "+919113068170"
	email := "contact@vithsutra.com"
	web_address := "www.vithsutra.com"

	x, y = 6.5, 2

	pdf.SetXY(x, y)
	pdf.Text(heading1)

	if err := pdf.SetFont("light-font", "", 14); err != nil {
		return err
	}

	x, y = 6, 2.75
	pdf.SetXY(x, y)
	pdf.Text("Phone: ")

	x = pdf.GetX()

	pdf.SetXY(x, y)
	pdf.Text(phone_number)

	x, y = pdf.GetX()+0.3, 2.75
	pdf.SetXY(x, y)
	pdf.Text("Email: ")

	x = pdf.GetX()

	pdf.SetXY(x, y)
	pdf.Text(email)

	x, y = 8.25, pdf.GetY()+0.75
	pdf.SetXY(x, y)
	pdf.Text("Web: ")

	x = pdf.GetX()

	pdf.SetXY(x, y)
	pdf.Text(web_address)

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.05)
	pdf.Line(1.5, 4, 19.5, 4)

	return nil
}

func EmployeeInfoSection(pdf *gopdf.GoPdf, employeeName, employeeCategory, date string) error {
	x, y := 1.5, 5.5

	if err := pdf.SetFont("bold-font", "", 15); err != nil {
		return err
	}

	pdf.SetXY(x, y)
	pdf.Text("Name:  ")

	if err := pdf.SetFont("light-font", "", 15); err != nil {
		return err
	}

	x = pdf.GetX()

	pdf.SetXY(x, y)
	pdf.Text(employeeName)

	if err := pdf.SetFont("bold-font", "", 15); err != nil {
		return err
	}

	x, y = 1.5, pdf.GetY()+0.7

	pdf.SetXY(x, y)
	pdf.Text("Position:  ")

	if err := pdf.SetFont("light-font", "", 15); err != nil {
		return err
	}

	x = pdf.GetX()
	pdf.SetXY(x, y)
	pdf.Text(employeeCategory)

	if err := pdf.SetFont("bold-font", "", 15); err != nil {
		return err
	}

	x, y = 15.25, 5.5

	pdf.SetXY(x, y)
	pdf.Text("Date:  ")

	if err := pdf.SetFont("light-font", "", 15); err != nil {
		return err
	}

	x = pdf.GetX()

	pdf.SetXY(x, y)
	pdf.Text(date)

	return nil

}

func TableHeaderSection(pdf *gopdf.GoPdf, startY float64) error {

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.05)
	pdf.Line(1, startY, 20, startY)
	pdf.Line(1, startY+1.5, 20, startY+1.5)

	pdf.Line(3, startY, 3, 28.7)
	pdf.Line(7.5, startY, 7.5, 28.7)
	pdf.Line(17, startY, 17, 28.7)

	if err := pdf.SetFont("bold-font", "", 14); err != nil {
		return err
	}

	heading1 := "SN"
	heading2 := "Date"
	heading3 := "Work Summary"
	heading4 := "Hrs"

	textWidth, err := pdf.MeasureTextWidth(heading1)

	if err != nil {
		return err
	}

	x, y := (2 - (textWidth / 2)), startY+0.9

	pdf.SetXY(x, y)
	pdf.Text(heading1)

	textWidth, err = pdf.MeasureTextWidth(heading2)

	if err != nil {
		return err
	}

	x, y = (5.25 - (textWidth / 2)), startY+0.9

	pdf.SetXY(x, y)
	pdf.Text(heading2)

	textWidth, err = pdf.MeasureTextWidth(heading3)

	if err != nil {
		return err
	}

	x, y = (12.25 - (textWidth / 2)), startY+0.9

	pdf.SetXY(x, y)
	pdf.Text(heading3)

	textWidth, err = pdf.MeasureTextWidth(heading4)

	if err != nil {
		return err
	}

	x, y = (18.5 - (textWidth / 2)), startY+0.9

	pdf.SetXY(x, y)
	pdf.Text(heading4)

	return nil
}

func TextWrapper(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine + " " + word
		width, _ := pdf.MeasureTextWidth(testLine)

		if width > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func TableSection(pdf *gopdf.GoPdf, startY float64, history []*models.UserWorkHistoryForPdf) (float64, string, error) {

	var isFirstPage bool = true

	var y float64 = startY + 2.25

	var userPerDayWorkHours []string

	for index1, history := range history {
		if isFirstPage {
			TableHeaderSection(pdf, startY)
			if err := pdf.SetFont("light-font", "", 14); err != nil {
				return 0.0, "", err
			}
			isFirstPage = false
		}

		if y > 28.7 {
			pdf.AddPage()
			TableHeaderSection(pdf, 5)
			if err := pdf.SetFont("light-font", "", 14); err != nil {
				return 0.0, "", err
			}
			y = 7.25
		}

		textWidth, err := pdf.MeasureTextWidth(strconv.Itoa(index1 + 1))

		if err != nil {
			return 0.0, "", err
		}

		x := (2 - (textWidth / 2))

		pdf.SetXY(x, y)
		pdf.Text(strconv.Itoa(index1 + 1))

		textWidth, err = pdf.MeasureTextWidth(history.Date)

		if err != nil {
			return 0.0, "", err
		}

		x = (5.25 - (textWidth / 2))

		pdf.SetXY(x, y)
		pdf.Text(history.Date)

		diffTime, err := CalculateTimeDiff(history.LoginTime, history.LogoutTime)

		if err != nil {
			return 0.0, "", err
		}

		userPerDayWorkHours = append(userPerDayWorkHours, diffTime)

		textWidth, err = pdf.MeasureTextWidth(diffTime)

		if err != nil {
			return 0.0, "", err
		}

		x = (18.5 - (textWidth / 2))

		pdf.SetXY(x, y)
		pdf.Text(diffTime)

		workSummaryLines := TextWrapper(pdf, history.WorkSummary, 8)

		var lineCounter float64 = 0

		for _, line := range workSummaryLines {

			textWidth, err = pdf.MeasureTextWidth(line)

			if err != nil {
				return 0.0, "", err
			}

			localX, localY := (12.25 - (textWidth / 2)), y+(lineCounter/1.8)
			pdf.SetXY(localX, localY)

			if pdf.GetY() > 28.2 {
				pdf.AddPage()
				TableHeaderSection(pdf, 5)
				if err := pdf.SetFont("light-font", "", 14); err != nil {
					return 0.0, "", err
				}
				y = 7.25
				localY = y
				lineCounter = 0
				pdf.SetXY(localX, localY)
			}
			pdf.Text(line)
			lineCounter++
		}

		pdf.SetY(pdf.GetY() + 0.5)

		if pdf.GetY() <= 28 {
			pdf.Line(1, pdf.GetY(), 20, pdf.GetY())
		}

		y = pdf.GetY() + 0.75

	}

	totalWorkHours, err := SumTimes(userPerDayWorkHours)

	if err != nil {
		return 0.0, "", err
	}

	return y, totalWorkHours, nil
}

func TotalWorkHoursSection(pdf *gopdf.GoPdf, startY float64, totalHrs string) error {
	if startY >= 27 {
		pdf.AddPage()
		TableHeaderSection(pdf, 5)
	}

	pdf.Line(1, 27.25, 20, 27.25)

	if err := pdf.SetFont("bold-font", "", 14); err != nil {
		return err
	}

	heading := "Total Work Hours"

	textWidth, err := pdf.MeasureTextWidth(heading)

	if err != nil {
		return err
	}

	x, y := (12.25 - (textWidth / 2)), 28.15

	pdf.SetXY(x, y)
	pdf.Text(heading)

	if err := pdf.SetFont("light-font", "", 14); err != nil {
		return err
	}

	textWidth, err = pdf.MeasureTextWidth(totalHrs)

	if err != nil {
		return err
	}

	x, y = (18.5 - (textWidth / 2)), 28.15
	pdf.SetXY(x, y)
	pdf.Text(totalHrs)

	return nil
}

func GenerateUserReportPdf(data *models.UserReportPdf) (string, error) {
	pdf := gopdf.GoPdf{}

	pdf.Start(
		gopdf.Config{
			PageSize: *gopdf.PageSizeA4,
			Unit:     gopdf.UnitCM,
		},
	)

	if err := pdf.AddTTFFont("bold-font", "./fonts/Roboto/static/Roboto-Bold.ttf"); err != nil {
		return "", err
	}

	if err := pdf.AddTTFFont("light-font", "./fonts/Roboto/static/Roboto-Regular.ttf"); err != nil {
		return "", err
	}

	pdf.AddHeader(
		func() {
			if err := HeaderSection(&pdf); err != nil {
				log.Println(err)
			}
		},
	)

	pdf.AddPage()

	currentDate := time.Now().Format("02-01-2006")

	if err := EmployeeInfoSection(&pdf, data.Name, data.Position, currentDate); err != nil {
		return "", err
	}

	lastY, totalWorkHours, err := TableSection(&pdf, 7.2, data.History)

	if err != nil {
		return "", err
	}

	if err := TotalWorkHoursSection(&pdf, lastY, totalWorkHours); err != nil {
		return "", err
	}

	uid := uuid.New().String()

	if err := pdf.WritePdf(fmt.Sprintf("./users_cache/%s.pdf", uid)); err != nil {
		return "", err
	}

	return uid, nil

}
