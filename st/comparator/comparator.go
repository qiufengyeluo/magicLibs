package comparator

import "time"

//Comparator desc
//@type Comparator desc: Comparator function
type Comparator func(a, b interface{}) int

//StringComparator desc
//@Method StringComparator desc: default a fast comparison on strings
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func StringComparator(a, b interface{}) int {
	s1 := a.(string)
	s2 := b.(string)

	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

//IntComparator desc
//@Method IntComparator desc: default a fast comparison on int
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func IntComparator(a, b interface{}) int {
	aAss := a.(int)
	bAss := b.(int)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int8Comparator desc
//@Method Int8Comparator desc: default a fast comparison on int8
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Int8Comparator(a, b interface{}) int {
	aAss := a.(int8)
	bAss := b.(int8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int16Comparator desc
//@Method Int16Comparator desc: default a fast comparison on int16
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Int16Comparator(a, b interface{}) int {
	aAss := a.(int16)
	bAss := b.(int16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int32Comparator desc
//@Method Int32Comparator desc: default a fast comparison on int32
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Int32Comparator(a, b interface{}) int {
	aAss := a.(int32)
	bAss := b.(int32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int64Comparator desc
//@Method Int64Comparator desc: default a fast comparison on int64
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Int64Comparator(a, b interface{}) int {
	aAss := a.(int64)
	bAss := b.(int64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UIntComparator desc
//@Method UIntComparator desc: default a fast comparison on uint
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func UIntComparator(a, b interface{}) int {
	aAss := a.(uint)
	bAss := b.(uint)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt8Comparator desc
//@Method UInt8Comparator desc: default a fast comparison on uint8
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func UInt8Comparator(a, b interface{}) int {
	aAss := a.(uint8)
	bAss := b.(uint8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt16Comparator desc
//@Method UInt16Comparator desc: default a fast comparison on uint16
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func UInt16Comparator(a, b interface{}) int {
	aAss := a.(uint16)
	bAss := b.(uint16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt32Comparator desc
//@Method UInt32Comparator desc: default a fast comparison on uint32
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func UInt32Comparator(a, b interface{}) int {
	aAss := a.(uint32)
	bAss := b.(uint32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt64Comparator desc
//@Method UInt64Comparator desc: default a fast comparison on uint64
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func UInt64Comparator(a, b interface{}) int {
	aAss := a.(uint64)
	bAss := b.(uint64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Float32Comparator desc
//@Method Float32Comparator desc: default a fast comparison on float32
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Float32Comparator(a, b interface{}) int {
	aAss := a.(float32)
	bAss := b.(float32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Float64Comparator desc
//@Method Float64Comparator desc: default a fast comparison on float64
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func Float64Comparator(a, b interface{}) int {
	aAss := a.(float64)
	bAss := b.(float64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//ByteComparator desc
//@Method ByteComparator desc: default a fast comparison on byte
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func ByteComparator(a, b interface{}) int {
	aAss := a.(byte)
	bAss := b.(byte)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//RuneComparator desc
//@Method RuneComparator desc: default a fast comparison on  time.Time
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func RuneComparator(a, b interface{}) int {
	aAss := a.(rune)
	bAss := b.(rune)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//TimeComparator desc
//@Method TimeComparator desc: default a fast comparison on rune
//@Param  (interface{}) a
//@Param  (interface{}) b
//@Return (int) comparator result
func TimeComparator(a, b interface{}) int {
	aAss := a.(time.Time)
	bAss := b.(time.Time)
	switch {
	case aAss.After(bAss):
		return 1
	case aAss.Before(bAss):
		return -1
	default:
		return 0
	}
}
