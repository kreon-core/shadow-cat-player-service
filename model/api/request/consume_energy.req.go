package request

type ChangeEnergy struct {
	Amount int32 `json:"amount" validate:"required,gt=0"`
}
