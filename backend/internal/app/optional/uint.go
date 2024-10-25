package optional

type Uint struct {
	Val   uint
	IsSet bool
}

func (u Uint) Read() (uint, bool) {
	return u.Val, u.IsSet
}

func NewUint(val uint) Uint {
	return Uint{Val: val, IsSet: true}
}
