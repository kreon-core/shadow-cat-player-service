package request

type ClaimSignInReq struct {
	WeekID string `json:"week_id" binding:"required"`
}
