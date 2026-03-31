package config

type BankConfig struct {
	TransactionPattern string            `yaml:"transaction_pattern"`
	DateInputFormat    string            `yaml:"date_input_format"`
	DateOutputFormat   string            `yaml:"date_output_format"`
	GroupDate          string            `yaml:"group_date"`
	GroupDescription   string            `yaml:"group_description"`
	GroupAmount        string            `yaml:"group_amount"`
	MonthMap           map[string]string `yaml:"month_map"`
	ExcludePatterns    []string          `yaml:"exclude_patterns"`
	SectionStart       string            `yaml:"section_start"`
	SectionEnd         string            `yaml:"section_end"`
	DebitIsPositive    bool              `yaml:"debit_is_positive"`
	AccountAssets      string            `yaml:"account_assets"`
	AccountExpenses    string            `yaml:"account_expenses"`
	AccountIncome      string            `yaml:"account_income"`
	Currency           string            `yaml:"currency"`
}

func (c *BankConfig) SetDefaults() {
	if c.GroupDate == "" {
		c.GroupDate = "date"
	}
	if c.GroupDescription == "" {
		c.GroupDescription = "description"
	}
	if c.GroupAmount == "" {
		c.GroupAmount = "amount"
	}
	if c.DateOutputFormat == "" {
		c.DateOutputFormat = "2006-01-02"
	}
	if c.AccountAssets == "" {
		c.AccountAssets = "assets:bank:checking"
	}
	if c.AccountExpenses == "" {
		c.AccountExpenses = "expenses:unknown"
	}
	if c.AccountIncome == "" {
		c.AccountIncome = "income:unknown"
	}
	if c.Currency == "" {
		c.Currency = "$"
	}
}
