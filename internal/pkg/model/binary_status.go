package model

// BinaryStatus 二元状态 0-否 1-是
type BinaryStatus uint8

const (
	BinaryStatusYes BinaryStatus = 1
	BinaryStatusNo  BinaryStatus = 0
)

func (s BinaryStatus) Bool() bool {
	return s > 0
}

func (s BinaryStatus) Uint() uint8 {
	return uint8(s)
}

func BinaryStatusByBool(status bool) BinaryStatus {
	if status {
		return BinaryStatusYes
	}
	return BinaryStatusNo
}

func BinaryStatusByUint(status uint8) BinaryStatus {
	if status > 0 {
		return BinaryStatusYes
	}
	return BinaryStatusNo
}
