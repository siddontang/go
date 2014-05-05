package leveldb

import (
	"github.com/siddontang/golib/hack"
	"strconv"
)

func Int(v []byte, err error) (int64, error) {
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseInt(hack.String(v), 10, 64)
}

func Uint(v []byte, err error) (uint64, error) {
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseUint(hack.String(v), 10, 64)
}

func Float(v []byte, err error) (float64, error) {
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseFloat(hack.String(v), 64)
}

func String(v []byte, err error) (string, error) {
	if err != nil {
		return "", err
	} else if v == nil {
		return "", nil
	}

	return hack.String(v), nil
}

func Slice(v []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	} else if v == nil {
		return []byte{}, nil
	}

	return v, nil
}
