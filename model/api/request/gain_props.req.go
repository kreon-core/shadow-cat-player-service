package request

type GainProps struct {
	Props []GainProp `json:"props" validate:"required,dive"`
}

type GainProp struct {
	ConfigPropID int `json:"config_prop_id" validate:"required"`
	Level        int `json:"level"          validate:"required"`
	Quantity     int `json:"quantity"       validate:"required,min=1"`
}
