package views

type SelectOption struct {
	Value          string
	Label          string
	Hint           string
	Disabled       bool
	DisabledReason string
	IsHeader       bool
	Checked        bool
}

type SummaryItem struct {
	Label string
	Value string
}
