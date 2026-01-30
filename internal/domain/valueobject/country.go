package valueobject

import "github.com/shopspring/decimal"

type Country string

const (
	CountryIndia        Country = "India"
	CountryUnitedStates Country = "United States"
)

var taxRates = map[Country]decimal.Decimal{
	CountryIndia:        decimal.NewFromFloat(0.10), // 10%
	CountryUnitedStates: decimal.NewFromFloat(0.12), // 12%
}

func (c Country) GetTaxRate() decimal.Decimal {
	if rate, exists := taxRates[c]; exists {
		return rate
	}
	return decimal.Zero
}

func (c Country) String() string {
	return string(c)
}

func (c Country) IsValid() bool {
	return c != ""
}
