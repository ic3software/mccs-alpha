package constant

// Trading Status decides whether a business can perform transactions (already in accepted status).
var Trading = struct {
	Pending  string
	Accepted string
	Rejected string
}{
	Pending:  "tradingPending",
	Accepted: "tradingAccepted",
	Rejected: "tradingRejected",
}
