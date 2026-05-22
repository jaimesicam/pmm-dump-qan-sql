// Copyright (C) 2026 pmm-dump-load authors
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <https://www.gnu.org/licenses/>.


package sqlexport

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func formatRow(values []any, scanTypes []reflect.Type) (string, error) {
	if len(values) != len(scanTypes) {
		return "", fmt.Errorf("column count mismatch: %d values, %d types", len(values), len(scanTypes))
	}

	parts := make([]string, len(values))
	for i, v := range values {
		lit, err := formatValue(v, scanTypes[i])
		if err != nil {
			return "", fmt.Errorf("column %d: %w", i, err)
		}
		parts[i] = lit
	}

	return "(" + strings.Join(parts, ", ") + ")", nil
}

func formatValue(v any, st reflect.Type) (string, error) {
	if v == nil {
		return defaultLiteral(st), nil
	}

	switch val := v.(type) {
	case string:
		return quoteString(val), nil
	case time.Time:
		return quoteString(val.UTC().Format("2006-01-02 15:04:05")), nil
	case []any:
		return formatSlice(val, st.Elem())
	case bool:
		if val {
			return "1", nil
		}
		return "0", nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'g', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported value type %T", v)
	}
}

func formatSlice(values []any, elemType reflect.Type) (string, error) {
	if len(values) == 0 {
		return "[]", nil
	}

	parts := make([]string, len(values))
	for i, v := range values {
		lit, err := formatValue(v, elemType)
		if err != nil {
			return "", err
		}
		parts[i] = lit
	}

	return "[" + strings.Join(parts, ", ") + "]", nil
}

func defaultLiteral(st reflect.Type) string {
	switch st.Kind() {
	case reflect.String:
		return "''"
	case reflect.Slice:
		return "[]"
	default:
		return "0"
	}
}

func quoteString(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('\'')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\x00':
			b.WriteString(`\0`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('\'')
	return b.String()
}
