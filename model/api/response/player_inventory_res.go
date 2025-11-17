package response

type PlayerInventory struct {
	Skins []int        `json:"skins"`
	Props []PlayerProp `json:"props"`
	// Consumables    []int        `json:"consumables"`
	// SpecialItems   []int        `json:"special_items"`
	// CollectedItems []int        `json:"collected_items"`
}

type PlayerProp struct {
	PropID       string `json:"prop_id"`
	ConfigPropID int    `json:"config_prop_id"`
	Quantity     int    `json:"quantity"`
	Level        int    `json:"level"`
}
