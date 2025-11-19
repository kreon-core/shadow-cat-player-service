package request

type UnlockNewSkin struct {
	SkinIDs []int32 `json:"skin_ids" validate:"required"`
}
