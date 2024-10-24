package typ

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"
)

type Uints []uint

func (s *Uints) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan Uints field - source is not a string")
	}
	strVals := strings.Split(str, ",")
	*s = make([]uint, len(strVals))
	for i, v := range strVals {
		if val, err := strconv.ParseUint(v, 10, 32); err != nil {
			return err
		} else {
			(*s)[i] = uint(val)
		}
	}

	return nil
}

func (s Uints) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	strVals := make([]string, len(s))
	for i, v := range s {
		strVals[i] = strconv.FormatUint(uint64(v), 10)
	}

	return strings.Join(strVals, ","), nil
}
