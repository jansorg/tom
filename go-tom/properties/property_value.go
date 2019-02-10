package properties

import (
	"encoding/json"
	"fmt"

	"github.com/rhymond/go-money"
)

var ErrUnsupportedType = fmt.Errorf("unsupported property type")

type PropertyValue interface {
	json.Marshaler
	json.Unmarshaler
	fmt.Stringer

	Type() PropertyType
	PropertyID() string

	IsString() bool
	IsFloat() bool
	IsCurrency() bool

	AsString() (string, error)
	AsFloat() (float64, error)
	AsCurrency() (*money.Money, error)
}
