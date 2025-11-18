package controller

import "net/http"

type ShopH struct{}

func NewShopH() *ShopH {
	return &ShopH{}
}

func (ctrl *ShopH) GetShopItems(w http.ResponseWriter, r *http.Request) {}

func (ctrl *ShopH) PurchaseItem(w http.ResponseWriter, r *http.Request) {}
