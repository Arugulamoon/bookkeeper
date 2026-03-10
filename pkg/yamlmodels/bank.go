package yamlmodels

type BankData struct {
	Currencies []Currency `yaml:"currencies"`
	Banks      []Bank     `yaml:"banks"`
}

type Bank struct {
	Id       string        `yaml:"id"`
	Name     string        `yaml:"name"`
	Accounts []BankAccount `yaml:"accounts"`
}

type BankAccount struct {
	Name  string   `yaml:"name"`
	Type  string   `yaml:"type"`
	Files []string `yaml:"files"`

	PaymentMadeDescs []string `yaml:"payment_made_descriptions"`
	PaymentRecdDescs []string `yaml:"payment_received_descriptions"`
}
