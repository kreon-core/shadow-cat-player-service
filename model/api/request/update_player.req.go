package request

type UpdatePlayer struct {
	CurrentSkin   int32   `json:"current_skin"`
	EquippedProps []int32 `json:"equipped_props"`
}
