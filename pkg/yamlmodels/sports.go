package yamlmodels

import "time"

// TODO: Change date string to time.Time using custom unmarshaller

type SportsData struct {
	Registrations []Registration `yaml:"registrations"`
	Memberships   []Membership   `yaml:"memberships"`
}

type Registration struct {
	Name     string `yaml:"name"`
	Price    Price  `yaml:"price"`
	Location string `yaml:"location"`
	Time     Time   `yaml:"time"`
	Date     Date   `yaml:"date"`
	Sessions int    `yaml:"sessions"`
	Notes    string `yaml:"notes"`
}

type Price struct {
	Regular  int `yaml:"regular"`
	Discount int `yaml:"discount"`
	Tax      int `yaml:"tax"`
	Total    int `yaml:"total"`
}

type Time struct {
	Day      string `yaml:"day"`
	Start    string `yaml:"start"`
	End      string `yaml:"end"`
	Range    Range  `yaml:"range"`
	Duration int    `yaml:"duration"`
}

type Range struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type Date struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type Membership struct {
	Name     string `yaml:"name"`
	Season   Season `yaml:"season"`
	Location string `yaml:"location"`
	Games    []Game `yaml:"games"`
}

type Season struct {
	Year string `yaml:"year"`
	Type string `yaml:"type"`
}

type Game struct {
	Date     time.Time `yaml:"date"`
	Time     Time      `yaml:"time"`
	Opponent string    `yaml:"opponent"`
	Type     string    `yaml:"type"`
	Notes    string    `yaml:"notes"`
	Location string    `yaml:"location"`
	Event    Event     `yaml:"event"`
}
