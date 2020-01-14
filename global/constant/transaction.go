package constant

// Transaction Status
var Transaction = struct {
	Initiated string
	Completed string
	Cancelled string
}{
	Initiated: "transactionInitiated",
	Completed: "transactionCompleted",
	Cancelled: "transactionCancelled",
}
