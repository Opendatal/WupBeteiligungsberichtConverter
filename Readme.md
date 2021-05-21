# Beteiligungsdaten Wuppertal 2018 Datenaufbereitung

## F ups

- Mehrere bilanzen auf einer Seite (Lokalfunk)
- Buchungskonten unterschiedlich benannt
- Jahreszahlen in unterschiedlichen formaten
- Buchungskontennummen währen echt nice gewesen!

## The json file

Das Ergebnis liegt in `out` [Beteiligungsbilanzen-2018.json](./out/Beteiligungsbilanzen-2018.json)

``` go

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
```





## Sourcen

Originaldatei: [offenedaten-wuppertal](https://www.offenedaten-wuppertal.de/dataset/beteiligungsmanagement)
Bericht: [PDF](https://www.wuppertal.de/vv/produkte/Finanzen/Beteiligungsmanagement.php.media/311227/Beteiligungsbericht_2018.pdf)