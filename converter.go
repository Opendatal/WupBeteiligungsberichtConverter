package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

type company struct {
	Name          string         `json:"company_name"`
	AnnualReports []annualReport `json:"annual_report"`
}

type annualReport struct {
	Year    int64   `json:"year"`
	Assets  assets  `json:"assets"`
	Pasiva  pasiva  `json:"pasiva"`
	Profits profits `json:"profits"`
	Loss    loss    `json:"losses"`
}

type assets struct {
	// Activa | Assets
	FixedAssets         int64 `json:"fixed_assets"`          // Anlagevermögen
	CurrentAssets       int64 `json:"current_assets"`        // Umlaufvermögen
	PrepaidExpenses     int64 `json:"prepaid_expenses"`      // Rechnungsabgrenzungsposten
	NotCoveredEquity    int64 `json:"not_covered_equity"`    // Nicht durch EK gedeckter Fehlbetrag
	OwnAccountTransport int64 `json:"own_account_transport"` // Ausgleichsposten Eigenmittelbefö.
}

type pasiva struct {
	// Pasiva
	Equity                  int64 `json:"equity"`                    // Eigenkapital
	SpecialItems            int64 `json:"special_items"`             // Sonderposten + Ertragszuschüsse
	Provisions              int64 `json:"provisions"`                // Rückstellungen
	Liabilities             int64 `json:"liabilities"`               // Verbindlichkeiten
	DeferredIncome          int64 `json:"deferred_income"`           // Rechnungsabgrenzungsposten
	ShareholderLoan         int64 `json:"shareholder_loan"`          // Gesellschafterdarlehen
	LoanSubsidiesAdjustment int64 `json:"loan_subsidies_adjustment"` // Ausgleichsposten aus Darlehensfö.
}

type profits struct {
	Revenues              int64 `json:"Revenues"`                   // Umsatzerlöse
	OtherOperating        int64 `json:"other_operating"`            // Sonstige betriebliche/sonst. Erträge
	OtherInterest         int64 `json:"other_interest"`             // Sonstige Zinsen und ähnliche Erträge
	LossAbsorption        int64 `json:"loss_absorption"`            // Erträge aus Verlustübernahme
	CitySubsidies         int64 `json:"city_subsidies"`             //Zuschüsse Stadt
	WorkProgressInventory int64 `json:"work_in_progress_inventory"` // Bestand in Arbeit befindliche Aufträge
}

type loss struct {
	MaterialExpenses    int64 `json:"material_expenses"`    // Materialaufwand
	PersonnelExpenses   int64 `json:"personnel_expenses"`   // Personalaufwand
	Depreciation        int64 `json:"depreciation"`         // Abschreibungen
	OperatingExpenses   int64 `json:"operating_expenses"`   // Sonstige betriebliche Aufwendungen
	InterestExpenses    int64 `json:"interest_expenses"`    // Zinsen und ähnliche Aufwendungen
	Taxes               int64 `json:"taxes"`                // Steuern
	LossAbsorption      int64 `json:"loss_absorption"`      //Aufwendungen aus Verlustübernahme
	RelatedServices     int64 `json:"related_services"`     // Aufwendungen für bez. Leistungen
	ExtraordinaryResult int64 `json:"extraordinary_result"` // a.o. Ergebnis  --- gibt es seit 2016 nicht me
}

func rowVisitor(r *xlsx.Row) error {

	rowBuffer = []string{}

	cv := func(c *xlsx.Cell) error {
		value, err := c.FormattedValue()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if strings.TrimSpace(value) != "" {
				rowBuffer = append(rowBuffer, value)
			}
		}
		return err
	}

	err := r.ForEachCell(cv)

	sheetbuffer = append(sheetbuffer, rowBuffer)
	cellCount = 0
	rowCount++

	return err
}

var sheetbuffer [][]string
var rowBuffer []string
var rowCount int = 0
var cellCount int = 0

func main() {
	filename := "./res/Beteiligungsdaten-2018-mod.xlsx"
	wb, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}

	var companies []company

	//todo  for each sheet
	for _, sh := range wb.Sheets {
		comp := company{}
		comp.Name = sh.Name

		sh.ForEachRow(rowVisitor)

		fmt.Printf("----------------------------------->> %s\n", sh.Name)

		// there are dublicated names we switch to pasiva if reacht the first occurence
		pasiva := false
		// use a switch for each column because properties are not fixed
		for r, sl := range sheetbuffer {
			// fmt.Printf("%v\n", r)
			for _, value := range sl {
				// fmt.Printf(">%v> %v\n", c, value)

				switch {
				case value == "Bilanz":
					fallthrough
				case value == "Konzern-Bilanz":
					fallthrough
				case value == "Bilanz:":
					// YEARS: find "Bilanz:" -> row -1

					for _, yearValue := range sheetbuffer[r-1] {
						// if year is longer 4 use the last 4 digits
						if len(yearValue) >= 4 {

							year, err := strconv.ParseInt(yearValue[len(yearValue)-4:], 10, 64)
							if err != nil {
								//do nothing if it's not a year
							} else {
								report := annualReport{}
								report.Year = year
								comp.AnnualReports = append(comp.AnnualReports, report)
							}
						}
					}

				case strings.HasPrefix(value, "Anlagevermögen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Assets.FixedAssets = v * 1000
					}
				case strings.HasPrefix(value, "Umlaufvermögen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Assets.CurrentAssets = v * 1000
					}
				case strings.HasPrefix(value, "Rechnungsabgrenzungsposten"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						if pasiva {
							comp.AnnualReports[i].Pasiva.DeferredIncome = v * 1000
						} else {
							comp.AnnualReports[i].Assets.PrepaidExpenses = v * 1000
							pasiva = true
						}
					}
				case strings.HasPrefix(value, "Nicht durch Eigenkapital gedeckter Fehlbetrag"):
					fallthrough
				case strings.HasPrefix(value, "Nicht durch Vermögenseinlagen gedeckter"):
					fallthrough
				case strings.HasPrefix(value, "Nicht durch EK gedeckter Fehlbetrag"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Assets.NotCoveredEquity = v * 1000
					}

				case strings.HasPrefix(value, "Ausgleichsposten Eigenmittelbef"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Assets.OwnAccountTransport = v * 1000
					}

				case strings.HasPrefix(value, "Eigenkapital"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.Equity = v * 1000
					}
				case strings.HasPrefix(value, "Sonderposten"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.SpecialItems = v * 1000
					}
				case strings.HasPrefix(value, "Rückstellungen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.Provisions = v * 1000
					}
				case strings.HasPrefix(value, "Verbindlichkeiten"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.Liabilities = v * 1000
					}
				case strings.HasPrefix(value, "Gesellschafterdarlehen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.ShareholderLoan = v * 1000
					}
				case strings.HasPrefix(value, "Ausgleichsposten aus Darlehensf"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Pasiva.LoanSubsidiesAdjustment = v * 1000
					}

				case strings.HasPrefix(value, "Zuschüsse Stadt"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.CitySubsidies = v * 1000
					}
				case strings.HasPrefix(value, "Umsatzerlöse"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.Revenues = v * 1000
					}
				case strings.HasPrefix(value, "Sonstige betriebliche/sonst. Erträge"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.OtherOperating = v * 1000
					}
				case strings.HasPrefix(value, "Sonstige Zinsen und ähnliche Erträge"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.OtherInterest = v * 1000
					}
				case strings.HasPrefix(value, "Erträge aus Verlustübernahme"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.LossAbsorption = v * 1000
					}
				case strings.HasPrefix(value, "Bestand in Arbeit befindliche Aufträge"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Profits.WorkProgressInventory = v * 1000
					}

				case strings.HasPrefix(value, "Materialaufwand"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.MaterialExpenses = v * 1000
					}
				case strings.HasPrefix(value, "Personalaufwand"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.PersonnelExpenses = v * 1000
					}
				case strings.HasPrefix(value, "Abschreibungen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.Depreciation = v * 1000
					}
				case strings.HasPrefix(value, "Aufwendungen für bez. Leistungen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.RelatedServices = v * 1000
					}
				case strings.HasPrefix(value, "Sonstige/ betriebliche Aufwendungen"):
					fallthrough
				case strings.HasPrefix(value, "Sonstige betriebl. Aufwendungen"):
					fallthrough
				case strings.HasPrefix(value, "Sonstige/betriebliche Aufwendungen"):
					fallthrough
				case strings.HasPrefix(value, "Sonstige betriebliche Aufwendungen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.OperatingExpenses = v * 1000
					}
				case strings.HasPrefix(value, "Zinsen und ähnliche Aufwendungen"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.InterestExpenses = v * 1000
					}
				case strings.HasPrefix(value, "Steuern"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.Taxes = v * 1000
					}
				case strings.HasPrefix(value, "Aufwand aus Verlustübernahme"):
					fallthrough
				case strings.HasPrefix(value, "Aufwendungen aus Verlustübernahme"):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.LossAbsorption = v * 1000
					}
				case strings.HasPrefix(value, "a.o. Ergebnis  "):
					for i := 0; i < len(sheetbuffer[r])-1; i++ {
						v, _ := strconv.ParseInt(sheetbuffer[r][i+1], 10, 64)
						comp.AnnualReports[i].Loss.ExtraordinaryResult = v * 1000
					}

				default:
					_, err := strconv.ParseFloat(value, 64)
					if err != nil {
						if contains(value) != true {
							fmt.Printf("ERROR: ---------  %v\n", value)
						}

					}

				}

			}
		}

		companies = append(companies, comp)
		//reset slice
		sheetbuffer = nil
		rowCount = 0
		cellCount = 0

	}

	file, _ := json.MarshalIndent(companies, "", " ")

	_ = ioutil.WriteFile("./out/Beteiligungsbilanzen-2018.json", file, 0644)
}

func contains(str string) bool {
	knownStrings := []string{"Jahresergebnis", "Aktiva in T €", "Passiva in T €", "EK Quote", "Jahresüberschuss / Fehlbetrag", "Gewinn- und Verlustrechnung in T€", "Gewinn- und Verlustrechnung in T€:"}

	for _, a := range knownStrings {
		if a == str {
			return true
		}
	}
	return false
}
