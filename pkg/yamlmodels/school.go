package yamlmodels

// TODO: Change date string to time.Time using custom unmarshaller

type SchoolData struct {
	Grades      []Grade      `yaml:"grades"`
	Schools     []School     `yaml:"schools"`
	SchoolYears []SchoolYear `yaml:"school_years"`
	Invoices    []Invoice    `yaml:"invoices"`
}

type Grade struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
}

type School struct {
	Id        string  `yaml:"id"`
	Name      string  `yaml:"name"`
	Address   *string `yaml:"address"`
	Phone     *string `yaml:"phone"`
	Principal *string `yaml:"principal"`
}

type SchoolYear struct {
	Year      string  `yaml:"school_year"`
	SchoolId  string  `yaml:"school"`
	GradeId   string  `yaml:"grade"`
	Teacher   *string `yaml:"teacher"`
	Education *string `yaml:"education"`
}

type Invoice struct {
	DueDate         string         `yaml:"due_date"`
	SchoolYear      string         `yaml:"school_year"`
	School          string         `yaml:"school"`
	Grade           string         `yaml:"grade"`
	Description     string         `yaml:"description"`
	Amount          int            `yaml:"amount"`
	Event           Event          `yaml:"event"`
	DatePaid        *string        `yaml:"date_paid"`
	EventMarkedPaid bool           `yaml:"event_marked_paid"`
	Reimbursement   *Reimbursement `yaml:"reimbursement"`
}

type Event struct {
	Id string `yaml:"id"`
}

type Reimbursement struct {
	Split  string  `yaml:"split"`
	Amount *int    `yaml:"amount"`
	Date   *string `yaml:"date"`
}
