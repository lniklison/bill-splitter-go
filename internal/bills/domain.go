package bills

type Cents int64

type Bill struct {
	ID          int64
	Description string
	TotalAmount Cents
	Shares      []Share
}

type Share struct {
	Position       int
	PersonName     string
	NormalizedName string
	Percentage     Percentage
}

type AllocatedShare struct {
	PersonName string
	Percentage Percentage
	Amount     Cents
}

type Allocation struct {
	BillID      int64
	Description string
	TotalAmount Cents
	Shares      []AllocatedShare
}
