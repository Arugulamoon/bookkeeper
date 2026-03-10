package yamlmodels

// TODO: Change date string to time.Time using custom unmarshaller

type SportsData struct {
	Registrations []Registration `yaml:"registrations"`
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
