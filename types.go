package rendezvoushash

type (
	ordered interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
			~float32 | ~float64 | ~string
	}
)

func compareOrdered[T ordered](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
