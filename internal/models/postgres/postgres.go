package postgres

import "time"

type (
	User struct {
		ID        int
		Username  string
		Email     string
		CountryID int       `db:"country_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	Gateway struct {
		ID                  int
		Name                string
		DataFormatSupported string    `db:"data_format_supported"`
		CreatedAt           time.Time `db:"created_at"`
		UpdatedAt           time.Time `db:"updated_at"`
	}

	Country struct {
		ID        int
		Name      string
		Code      string
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	Transaction struct {
		ID        int64
		Amount    float64
		Type      string
		Status    string
		UserID    int       `db:"user_id"`
		GatewayID int       `db:"gateway_id"`
		CountryID int       `db:"country_id"`
		CreatedAt time.Time `db:"created_at"`
	}
)
