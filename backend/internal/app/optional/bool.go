package optional

type Bool int8

var NoneBool = Boole(false, false)

func Boole(value, isValue bool) Bool {
	if !isValue {
		return -1
	} else if value {
		return 1
	} else {
		return 0
	}
}

func (b Bool) Val() bool {
	return b == 1
}

func (b Bool) None() bool {
	return b == -1
}
