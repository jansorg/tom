package tristate

import (
	"fmt"

	"github.com/jansorg/tom/go-tom/util"
)

type Tristate int8

const (
	False Tristate = iota + 1
	True
	Inherited
)

func FalseP() *Tristate {
	v := False
	return &v
}

func TrueP() *Tristate {
	v := True
	return &v
}

func InheritedP() *Tristate {
	v := Inherited
	return &v
}

func FromBool(value *bool) Tristate {
	if value == nil {
		return Inherited
	}
	if *value == true {
		return True
	}
	return False
}

func FromString(value string) (Tristate, error) {
	if value == "true" {
		return True, nil
	}
	if value == "false" {
		return False, nil
	}
	if value == "" || value == "inherit" {
		return Inherited, nil
	}
	return False, fmt.Errorf("unable to parse %s, valid values: true, false, inherit, or empty value", value)
}

func (t Tristate) ToBool() *bool {
	if t == Inherited {
		return nil
	}
	if t == True {
		return util.TrueP()
	}
	return util.FalseP()
}

func (t Tristate) IsInherited() bool {
	return t == Inherited
}

func (t Tristate) IsTrue() bool {
	return t == True
}

func (t Tristate) IsFalse() bool {
	return t != Inherited && t != True
}
