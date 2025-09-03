package shop

type ShopService struct {
	shop Shop
}

func NewShopService(shop Shop) *ShopService {
	return &ShopService{shop: shop}
}
