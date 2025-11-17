package request

type UpdatePlayer struct {
	CurrentSkin   int   `json:"current_skin"`
	EquippedProps []int `json:"equipped_props"`
}
