package domain

type Merch struct {
	Id    int    `json:"-" db:"id"`
	Name  string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
}

type Transactions struct {
	Id int `json:"-" db:"id"`
	//Source      *int              `json:"source"`
	Destination int `json:"destination" binding:"required"`
	Amount      int `json:"amount" binding:"required"`
	//Timestamp   *pgtype.Timestamp `json:"timestamp"`
}
