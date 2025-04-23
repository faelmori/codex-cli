package utils

type AutoApprovalMode string

const (
	Suggest   AutoApprovalMode = "suggest"
	AutoEdit  AutoApprovalMode = "auto-edit"
	FullAuto  AutoApprovalMode = "full-auto"
)

type FullAutoErrorMode string

const (
	AskUser            FullAutoErrorMode = "ask-user"
	IgnoreAndContinue  FullAutoErrorMode = "ignore-and-continue"
)
