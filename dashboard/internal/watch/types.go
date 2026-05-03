package watch

type Wallet struct {
	Address  string
	Nickname string
	Grade    string
}

type Signal struct {
	Title          string
	Slug           string
	Side           string
	Entry          float64
	Current        float64
	Score          float64
	Recommendation string
	SuggestedSize  float64
	EndDate        string
}

type WalletReport struct {
	Wallet            Wallet
	Date              string
	OutputPath        string
	PortfolioValueUSD float64
	OpenPositions     int
	Signals           []Signal
	StrongCount       int
	ClosedRecentCount int
	Warnings          []string
}

type RunSummary struct {
	Reports []WalletReport
	Errors  []string
}
