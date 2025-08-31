package shop

// not db model, because shop items are static and infinite
type Shop map[string]int

func NewShop() Shop {
	return Shop{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}
}

// func (s Shop) BuyItem(ctx context.Context, emp *employee.EmployeeDto, item string) error {
// 	const op = "shop.model.BuyItem"

// 	price, ok := s[item]
// 	if !ok {
// 		return fmt.Errorf("%s: item %s not found", op, item)
// 	}

// 	if emp.Coins-price < 0 {
// 		return fmt.Errorf("%s: not enough coins to buy %s", op, item)
// 	}

// 	emp.Coins -= price
// 	emp.Items = append(emp.Items, item)

// 	return nil
// }
