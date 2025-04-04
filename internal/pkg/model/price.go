package model

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Price 金额字段
type Price int64

const (
	CNYFen  Price = 1
	CNYJiao       = 10 * CNYFen
	CNYYuan       = 10 * CNYJiao
)

func (p Price) Yuan() float64 {
	return float64(p) / float64(CNYYuan)
}

func (p Price) Jiao() float64 {
	return float64(p) / float64(CNYJiao)
}

func (p Price) Fen() int64 {
	return int64(p)
}

// Format 金额格式化
// @param format string 格式
// @param isThousand bool 按千分位格式化
func (p Price) Format(format string, isThousand bool) string {
	if isThousand {
		printer := message.NewPrinter(language.English)
		return printer.Sprintf(format, p.Yuan())
	}
	return fmt.Sprintf(format, p.Yuan())
}

// PriceColor 价格颜色
const PriceColor = `#08AF5D`

// FormatColor 带颜色标签的格式化
// @param format string 格式
// @param color string 颜色值
// @param isThousand bool 按千分位格式化
func (p Price) FormatColor(format, color string, isThousand bool) string {
	return p.Format("<span style=\"color: "+color+"\">"+format+"</span>", isThousand)
}

func PriceYuan(yuan float64) Price {
	return Price(yuan * float64(CNYYuan))
}

// Percent 百分比字段
type Percent uint

func (p Percent) Valid() bool {
	return p <= 10000
}

// Numeric 换算为百分比形式 eg. 10%==>10
func (p Percent) Numeric() float64 {
	return float64(p) / 100
}

// Decimal 换算为小数形式 例如 10%==>0.1
func (p Percent) Decimal() float64 {
	return float64(p) / 10000
}

// FormatNumeric 格式化百分比 例如 10%==>10
func (p Percent) FormatNumeric() string {
	return fmt.Sprintf("%.2f", p.Numeric())
}

// FormatDecimal 格式化百分比小数 例如 10%==>0.1
func (p Percent) FormatDecimal() string {
	return fmt.Sprintf("%.2f", p.Decimal())
}

// FormatInt 格式化整数 例如 10%==>1
func (p Percent) FormatInt() string {
	return fmt.Sprintf("%.0f", p.Numeric())
}

// NumericPercent 根据百分比形式创建Percent
func NumericPercent(per float64) Percent {
	return Percent(uint(per * 100))
}

// DecimalPercent 根据百分比小数形式创建Percent
func DecimalPercent(per float64) Percent {
	return Percent(uint(per * 10000))
}

// CalculatePricePercent 计算金额的%
func CalculatePricePercent(price Price, percent Percent) Price {
	if price <= 0 || percent <= 0 {
		return Price(0)
	}
	return Price(price.Fen() * int64(percent) / 10000)
}

// Score 评分字段
type Score uint8

func (s Score) Valid() bool {
	return s <= 50
}

func (s Score) Score() float32 {
	return float32(s) / 10
}

func (s Score) Format() string {
	return fmt.Sprintf("%.1f", float64(s)/10)
}

func (s Score) FormatFloat() float32 {
	score := fmt.Sprintf("%.1f", float64(s)/10)
	sc, _ := strconv.ParseFloat(score, 64)
	return float32(sc)
}

// NewScore 新建评分
func NewScore(score float32) Score {
	s := Score(score * 10)
	if !s.Valid() {
		s = Score(0)
	}
	return s
}
