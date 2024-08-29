package model

import "fmt"

// Weight 重量(克)
type Weight uint64

const (
	weightGram Weight = 1                 // 克
	weightKG          = 1000 * weightGram // 千克
	weightTon         = 1000 * weightKG   // 吨
)

func (w Weight) Gram() uint64 {
	return uint64(w)
}

func (w Weight) Kg() float64 {
	return float64(w) / float64(weightKG)
}

func (w Weight) Ton() float64 {
	return float64(w) / float64(weightTon)
}

func (w Weight) FormatKg(format string) string {
	return fmt.Sprintf(format, float64(w)/float64(weightKG))
}

func (w Weight) FormatTon(format string) string {
	return fmt.Sprintf(format, float64(w)/float64(weightTon))
}

// WeightByKg 以千克创建重量
func WeightByKg(weight float64) Weight {
	return Weight(weight * float64(weightKG))
}

// WeightByTon 以吨创建重量
func WeightByTon(weight float64) Weight {
	return Weight(weight * float64(weightTon))
}
