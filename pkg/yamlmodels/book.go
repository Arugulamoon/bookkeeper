package yamlmodels

type BookData struct {
	Currencies []Currency `yaml:"currencies"`
	Accounts   []Account  `yaml:"accounts"`
	Assigners  []Assigner `yaml:"assigners"`
}

type Currency struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Account struct {
	AccountType string  `yaml:"account_type"`
	Name        string  `yaml:"name"`
	BankAccount *string `yaml:"bank_account"`
	SortOrder   *int    `yaml:"sort_order"`
}

type Assigner struct {
	Name         string   `yaml:"name"`
	AccountType  string   `yaml:"account_type"`
	Account      string   `yaml:"account"`
	Descriptions []string `yaml:"descriptions"`
}
