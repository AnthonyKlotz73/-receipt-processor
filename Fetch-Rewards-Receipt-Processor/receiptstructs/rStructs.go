package receiptstructs

/*
a Post response from a JSON file
*/
type PostResponse struct {
	Id string `json:"id"`
}

/*
a GET response from a JSON file
*/
type GetResponse struct {
	Points string `json:"points"`
}

/*
An item from a receipt
*/
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"` //might change this to a float
}

/*
Represents a Receipt from a struct given from teh fetch Github
*/
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}
