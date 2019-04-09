package gomybatis

import (
	"gohelper/convert"
	"reflect"
)

func compare(value interface{}, op, testValue string) bool {
	bl := false

	switch reflect.TypeOf(value).Name() {
	case "int":
		bl = compareInt(value.(int), convert.ToInt(testValue), op)
	case "int8":
		bl = compareInt8(value.(int8), convert.ToInt8(testValue), op)
	case "int16":
		bl = compareInt16(value.(int16), convert.ToInt16(testValue), op)
	case "int32":
		bl = compareInt32(value.(int32), convert.ToInt32(testValue), op)
	case "int64":
		bl = compareInt64(value.(int64), convert.ToInt64(testValue), op)
	case "uint":
		bl = compareUint(value.(uint), convert.ToUint(testValue), op)
	case "uint8":
		bl = compareUint8(value.(uint8), convert.ToUint8(testValue), op)
	case "uint16":
		bl = compareUint(value.(uint), convert.ToUint(testValue), op)
	case "uint32":
		bl = compareUint32(value.(uint32), convert.ToUint32(testValue), op)
	case "uint64":
		bl = compareUint64(value.(uint64), convert.ToUint64(testValue), op)
	case "float32":
		bl = compareFloat32(value.(float32), convert.ToFloat32(testValue), op)
	case "float64":
		bl = compareFloat64(value.(float64), convert.ToFloat64(testValue), op)
	case "string":
		bl = compareString(value.(string), testValue, op)
	}

	return bl
}

func compareInt(value, target int, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt8(value, target int8, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt16(value, target int16, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt32(value, target int32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareInt64(value, target int64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint(value, target uint, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint8(value, target uint8, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint16(value, target uint16, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint32(value, target uint32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareUint64(value, target uint64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareFloat32(value, target float32, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareFloat64(value, target float64, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target

	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}

func compareString(value, target string, op string) bool {
	bl := false
	switch op {
	case "!=":
		bl = value != target
	case "==":
		fallthrough
	case "=":
		bl = value == target

	case ">":
		bl = value > target

	case ">=":
		bl = value >= target

	case "<":
		bl = value < target

	case "<=":
		bl = value <= target

	}

	return bl
}
