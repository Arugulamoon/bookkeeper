package yamlmodels

type ConfigData struct {
	Bank   BankData   `yaml:"bank"`
	Book   BookData   `yaml:"book"`
	Sports SportsData `yaml:"sports"`
	School SchoolData `yaml:"school"`
}
