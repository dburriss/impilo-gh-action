package domain

type Command interface {
	Execute() []Report
	Title() string
}

type Report interface {
	Run()
	Title() string
}
