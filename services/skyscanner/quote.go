package skyscanner

// Quote ressource
type Quote struct {
	QuoteId       int
	MinPrice      float64
	Direct        bool
	OutboundLeg   OutBoundLeg
	QuoteDateTime string
}
