package processor

type RawRequest struct {
	Retailer     string  `json:"Retailer" description:"The name of the Retailer or store the receipt is from."`
	PurchaseDate string  `json:"purchaseDate" description:"The date of the purchase printed on the receipt."`
	PurchaseTime string  `json:"purchaseTime" description:"The time of the purchase printed on the receipt. 24-hour time expected."`
	Items        []Items `json:"Items" description:"List of Items purchased as in the receipt."`
	Total        string  `json:"Total" description:"The Total amount paid on the receipt."`
}

type Items struct {
	Description string `json:"shortDescription" description:"The Short Product Description for the item."`
	Price       string `json:"price" description:"The Total price payed for this item."`
}

type Response1 struct {
	Id string `json:"id"`
}

type Response2 struct {
	Points int64 `json:"points"`
}

type Request struct {
	Retailer     string
	PurchaseDay  int
	PurchaseHour int
	Items        []Item
	Total        float64
}

type Item struct {
	description string
	price       float64
}
