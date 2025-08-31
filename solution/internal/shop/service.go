package shop

type ShopService struct {
	shop Shop
}

func NewShopService(shop Shop) *ShopService {
	return &ShopService{shop: shop}
}

// func (s *ShopService) BuyItem(ctx context.Context, emp *employee.Employee, item string) error {
// 	return s.shop.BuyItem(ctx, emp, item)
// }
