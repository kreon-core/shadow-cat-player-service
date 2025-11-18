package request

type UnlockSignInReq struct {
	WeekID string `json:"week_id" binding:"required"`
}
