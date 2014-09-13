package num

import (
	"strconv"
)

func StringToUint(s string) (uint, error) {
	if v, err := strconv.ParseUint(s, 10, 0); err != nil {
		return 0, err
	} else {
		return uint(v), nil
	}
}

func StringToUint8(s string) (uint8, error) {
	if v, err := strconv.ParseUint(s, 10, 8); err != nil {
		return 0, err
	} else {
		return uint8(v), nil
	}
}

func StringToUint16(s string) (uint16, error) {
	if v, err := strconv.ParseUint(s, 10, 16); err != nil {
		return 0, err
	} else {
		return uint16(v), nil
	}
}

func StringToUint32(s string) (uint32, error) {
	if v, err := strconv.ParseUint(s, 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(v), nil
	}
}

func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func StringToInt(s string) (int, error) {
	if v, err := strconv.ParseInt(s, 10, 0); err != nil {
		return 0, err
	} else {
		return int(v), nil
	}
}

func StringToInt8(s string) (int8, error) {
	if v, err := strconv.ParseInt(s, 10, 8); err != nil {
		return 0, err
	} else {
		return int8(v), nil
	}
}

func StringToInt16(s string) (int16, error) {
	if v, err := strconv.ParseInt(s, 10, 16); err != nil {
		return 0, err
	} else {
		return int16(v), nil
	}
}

func StringToInt32(s string) (int32, error) {
	if v, err := strconv.ParseInt(s, 10, 32); err != nil {
		return 0, err
	} else {
		return int32(v), nil
	}
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func IntToString(v int) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int8ToString(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int16ToString(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int64ToString(v int64) string {
	return strconv.FormatInt(int64(v), 10)
}

func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint8ToString(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint16ToString(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint32ToString(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(uint64(v), 10)
}

func IntToSlice(v int) []byte {
	return strconv.AppendInt(nil, int64(v), 10)
}

func Int8ToSlice(v int8) []byte {
	return strconv.AppendInt(nil, int64(v), 10)
}

func Int16Slice(v int16) []byte {
	return strconv.AppendInt(nil, int64(v), 10)
}

func Int32ToSlice(v int32) []byte {
	return strconv.AppendInt(nil, int64(v), 10)
}

func Int64ToSlice(v int64) []byte {
	return strconv.AppendInt(nil, int64(v), 10)
}

func UintToSlice(v uint) []byte {
	return strconv.AppendUint(nil, uint64(v), 10)
}

func Uint8ToSlice(v uint8) []byte {
	return strconv.AppendUint(nil, uint64(v), 10)
}

func Uint16ToSlice(v uint16) []byte {
	return strconv.AppendUint(nil, uint64(v), 10)
}

func Uint32ToSlice(v uint32) []byte {
	return strconv.AppendUint(nil, uint64(v), 10)
}

func Uint64ToSlice(v uint64) []byte {
	return strconv.AppendUint(nil, uint64(v), 10)
}
