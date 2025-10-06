package valueobject

// PlatformType 平台类型值对象
type PlatformType uint8

const (
	PlatformWxMinioProgram PlatformType = iota + 1 // 微信小程序
	PlatformWxOfficial                             // 微信公众号
)

func (p PlatformType) String() string {
	switch p {
	case PlatformWxMinioProgram:
		return "WxMinioProgram"
	case PlatformWxOfficial:
		return "WxOfficial"
	}
	return ""
}
