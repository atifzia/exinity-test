package common

type (
	Gateway struct {
		ID                  int    `json:"id"`
		Name                string `json:"name"`
		DataFormatSupported string
		Priority            int `json:"priority"`
		CountryID           int `json:"country_id"`
	}
)
