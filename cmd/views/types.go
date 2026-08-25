package views

type SelectOption struct {
	Value          string
	Label          string
	Hint           string
	Disabled       bool
	DisabledReason string
}

type SummaryItem struct {
	Label string
	Value string
}
