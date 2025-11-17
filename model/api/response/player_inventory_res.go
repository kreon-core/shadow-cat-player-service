package response

type PlayerInventory struct {
	Skins []int        `json:"skins"`
	Props []PlayerProp `json:"props"`
	// Consumables    []int        `json:"consumables"`
	// SpecialItems   []int        `json:"special_items"`
	// CollectedItems []int        `json:"collected_items"`
}

type PlayerProp struct {
	// PropID is the unique identifier for this specific prop instance (e.g., a UUID or string).
	PropID       string `json:"prop_id"`
	// ConfigPropID references the configuration/template ID for the prop (from config data).
	ConfigPropID int    `json:"config_prop_id"`
	Quantity     int    `json:"quantity"`
	Level        int    `json:"level"`
}
