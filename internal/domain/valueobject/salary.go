package valueobject

import "github.com/shopspring/decimal"

type Salary struct {
	GrossSalary decimal.Decimal
	TaxRate     decimal.Decimal
	TaxAmount   decimal.Decimal
	NetSalary   decimal.Decimal
}

func CalculateNetSalary(grossSalary decimal.Decimal, country Country) Salary {
	taxRate := country.GetTaxRate()
	taxAmount := grossSalary.Mul(taxRate)
	netSalary := grossSalary.Sub(taxAmount)

	return Salary{
		GrossSalary: grossSalary,
		TaxRate:     taxRate,
		TaxAmount:   taxAmount.Round(2),
		NetSalary:   netSalary.Round(2),
	}
}

type SalaryStats struct {
	MinSalary decimal.Decimal
	MaxSalary decimal.Decimal
	AvgSalary decimal.Decimal
	Count     int64
}

type JobTitleSalaryStats struct {
	JobTitle  string
	AvgSalary decimal.Decimal
	Count     int64
}
