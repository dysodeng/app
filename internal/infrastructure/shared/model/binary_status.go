package model

// BinaryStatus 二元状态 0-否 1-是
type BinaryStatus uint8

const (
	BinaryStatusTrue  BinaryStatus = 1
	BinaryStatusFalse BinaryStatus = 0
)

func (s BinaryStatus) Bool() bool {
	return s > 0
}

func (s BinaryStatus) Uint() uint8 {
	return uint8(s)
}

func BinaryStatusByBool(status bool) BinaryStatus {
	if status {
		return BinaryStatusTrue
	}
	return BinaryStatusFalse
}

func BinaryStatusByUint(status uint8) BinaryStatus {
	if status > 0 {
		return BinaryStatusTrue
	}
	return BinaryStatusFalse
}
