package agent

type ReviewDecision string

const (
	ReviewDecisionYes        ReviewDecision = "yes"
	ReviewDecisionNoContinue ReviewDecision = "no-continue"
	ReviewDecisionNoExit     ReviewDecision = "no-exit"
	ReviewDecisionAlways     ReviewDecision = "always"
	ReviewDecisionExplain    ReviewDecision = "explain"
)
