package model

type Transaction struct {
	Id             string  `json:"-" bson:"_id,omitempty"`
	PurchaseAmount float64 `json:"purchase_amount,omitempty" bson:"purchase_amount,omitempty" validate:"required"`
	Description    string  `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	CreatedAt      int64   `json:"-" bson:"created_at,omitempty"`
	PurchaseDate   string  `json:"purchase_date" bson:"purchase_date,omitempty"`
}

type TransactionResponse struct {
	Id                      string  `json:"id,omitempty" bson:"_id,omitempty"`
	PurchaseAmount          float64 `json:"purchase_amount,omitempty" bson:"purchase_amount,omitempty" validate:"required"`
	Description             string  `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	CreatedAt               int64   `json:"-" bson:"created_at,omitempty"`
	PurchaseDate            string  `json:"purchase_date" bson:"purchase_date,omitempty"`
	ExchangeRate            float64 `json:"exchange_rate" bson:"-"`
	ConvertedPurchaseAmount float64 `json:"converted_purchase_amount" bson:"-"`
}

type GetTransactionParams struct {
	Ids      string `query:"ids" validate:"required"`
	Currency string `query:"currency" validate:"required"`
}

type GetTransactionParamsByPeriod struct {
	StartDate string `query:"startDate" validate:"required"`
	EndDate   string `query:"endDate" validate:"required"`
	Currency  string `query:"currency" validate:"required"`
}

type GetTransactionParamsByPeriodEpoch struct {
	StartDate int64  `query:"startDate" validate:"required"`
	EndDate   int64  `query:"endDate" validate:"required"`
	Currency  string `query:"currency" validate:"required"`
}
