package request

type (
	// Transaction is a standard request structure for the transactions
	Transaction struct {
		Amount    float64 `json:"amount" xml:"amount"`
		UserID    int     `json:"user_id" xml:"user_id"`
		GatewayID int     `json:"gateway_id" xml:"gateway_id"`
		CountryID int     `json:"country_id" xml:"country_id"`
		Currency  string  `json:"currency" xml:"currency"`
	}
)
