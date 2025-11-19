package request

type UpdatePlayer struct {
	CurrentSkin   int32   `json:"current_skin"   validate:"required"`
	EquippedProps []int32 `json:"equipped_props" validate:"required,dive,min=0"`
}
