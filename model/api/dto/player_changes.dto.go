package dto

type PlayerChanges struct {
	Level         *int32   `json:"level,omitempty"`
	Exp           *int32   `json:"exp,omitempty"`
	Coins         *int32   `json:"coins,omitempty"`
	Gems          *int32   `json:"gems,omitempty"`
	BestMap       *BestMap `json:"best_map,omitempty"`
	CurrentSkin   *int32   `json:"current_skin,omitempty"`
	EquippedProps *[]int32 `json:"equipped_props,omitempty"`
	Props         []Prop   `json:"props,omitempty"`
}
