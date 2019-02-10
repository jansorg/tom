package properties

var CurrencyType = &currencyType{}

func FindType(id string) PropertyType {
	switch id {
	case CurrencyType.ID():
		return CurrencyType
	default:
		return nil
	}
}
