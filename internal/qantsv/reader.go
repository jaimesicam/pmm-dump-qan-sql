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


package qantsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Reader struct {
	*csv.Reader
	scanTypes []reflect.Type
}

func NewReader(r io.Reader, scanTypes []reflect.Type) *Reader {
	reader := csv.NewReader(r)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 0
	return &Reader{Reader: reader, scanTypes: scanTypes}
}

func (r *Reader) Read() ([]any, error) {
	records, err := r.Reader.Read()
	if err != nil {
		return nil, err
	}
	if len(r.scanTypes) != len(records) {
		return nil, fmt.Errorf("amount of columns mismatch: expected %d, got %d", len(r.scanTypes), len(records))
	}

	values := make([]any, 0, len(records))
	for i, record := range records {
		value, err := parseElement(record, r.scanTypes[i])
		if err != nil {
			return nil, fmt.Errorf("parsing error: %w", err)
		}
		values = append(values, value)
	}

	return values, nil
}

func parseSlice(slice string, st reflect.Type) (any, error) {
	slice = strings.TrimSpace(slice[1 : len(slice)-1])
	elements := strings.Split(slice, ",")
	result := make([]any, 0, len(elements))
	if slice == "" {
		return result, nil
	}
	for _, v := range elements {
		value, err := parseElement(v, st)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func parseElement(record string, st reflect.Type) (any, error) {
	var value any
	var err error
	switch st.Kind() {
	case reflect.Slice:
		value, err = parseSlice(record, st.Elem())
		if err != nil {
			return nil, err
		}
	case reflect.Int8:
		result, err := strconv.ParseInt(record, 10, 8)
		if err != nil {
			return nil, err
		}
		value = int8(result)
	case reflect.Int16:
		result, err := strconv.ParseInt(record, 10, 16)
		if err != nil {
			return nil, err
		}
		value = int16(result)
	case reflect.Int32:
		result, err := strconv.ParseInt(record, 10, 32)
		if err != nil {
			return nil, err
		}
		value = int32(result)
	case reflect.Int64:
		value, err = strconv.ParseInt(record, 10, 64)
		if err != nil {
			return nil, err
		}
	case reflect.Uint8:
		result, err := strconv.ParseUint(record, 10, 8)
		if err != nil {
			return nil, err
		}
		value = uint8(result)
	case reflect.Uint16:
		result, err := strconv.ParseUint(record, 10, 16)
		if err != nil {
			return nil, err
		}
		value = uint16(result)
	case reflect.Uint32:
		result, err := strconv.ParseUint(record, 10, 32)
		if err != nil {
			return nil, err
		}
		value = uint32(result)
	case reflect.Uint64:
		value, err = strconv.ParseUint(record, 10, 64)
		if err != nil {
			return nil, err
		}
	case reflect.Float32:
		result, err := strconv.ParseFloat(record, 32)
		if err != nil {
			return nil, err
		}
		value = float32(result)
	case reflect.Float64:
		value, err = strconv.ParseFloat(record, 64)
		if err != nil {
			return nil, err
		}
	case reflect.String:
		value = record
	default:
		switch st.Name() {
		case "Time":
			value, err = time.Parse("2006-01-02 15:04:05 -0700 UTC", record)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown type")
		}
	}
	return value, nil
}
