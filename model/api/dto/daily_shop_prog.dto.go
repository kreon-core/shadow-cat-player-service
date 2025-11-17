package dto

type DailyShopProgress struct {
	PurchasedItems []PurchasedItem `json:"purchased_items"`
}

type PurchasedItem struct {
	ItemID int `json:"item_id"`
	Amount int `json:"amount"`
}
