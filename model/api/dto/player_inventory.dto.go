package response

type PlayerInventory struct {
	Skins []int  `json:"skins"`
	Props []Prop `json:"props"`
}

type Prop struct {
	PropID       string `json:"prop_id"`
	ConfigPropID int    `json:"config_prop_id"`
	Level        int    `json:"level"`
	Quantity     int    `json:"quantity"`
}
