package models

import "time"

type SportsRegistration struct {
	Id                  string
	Name                string
	PriceInCents        int
	RegularPriceInCents int
	DiscountInCents     int
	TaxInCents          int
	Location            string
	Day                 string
	StartTime           string
	StartTimeRange      string
	EndTimeRange        string
	Duration            int
	StartDate           time.Time
	EndDate             time.Time
	Sessions            int
}

func (reg SportsRegistration) IsUpcoming() bool {
	now := time.Now()
	if now.Before(reg.StartDate) &&
		now.AddDate(0, 1, 0).After(reg.StartDate) {

		return true
	}
	return false
}

func (reg SportsRegistration) IsCurrent() bool {
	now := time.Now()
	if now.After(reg.StartDate) && now.Before(reg.EndDate) {
		return true
	}
	return false
}

func (reg SportsRegistration) Price() int {
	return reg.PriceInCents / 100
}

func (reg SportsRegistration) PricePerSession() int {
	return reg.PriceInCents / reg.Sessions / 100
}
