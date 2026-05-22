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

package metricsschema

import (
	_ "embed"
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"time"
)

//go:embed metrics-pmm2.ddl
var metricsPMM2DDL string

//go:embed metrics-pmm3.ddl
var metricsPMM3DDL string

var (
	scanTypesCache sync.Map // Version -> []reflect.Type
)

// ScanTypes returns reflect types for pmm.metrics columns in CREATE TABLE order.
func ScanTypes(version Version) ([]reflect.Type, error) {
	if version != PMM2 && version != PMM3 {
		return nil, fmt.Errorf("unsupported PMM version %v", version)
	}

	if cached, ok := scanTypesCache.Load(version); ok {
		return cached.([]reflect.Type), nil
	}

	ddl := metricsPMM3DDL
	if version == PMM2 {
		ddl = metricsPMM2DDL
	}

	scanTypes, err := parseMetricsDDL(ddl)
	if err != nil {
		return nil, err
	}

	scanTypesCache.Store(version, scanTypes)
	return scanTypes, nil
}

var columnLineRE = regexp.MustCompile("`([^`]+)`\\s+([^,\\n]+)")

func parseMetricsDDL(ddl string) ([]reflect.Type, error) {
	lines := make([]string, 0)
	for _, line := range regexp.MustCompile("\n").Split(ddl, -1) {
		line = regexp.MustCompile(`^\s+`).ReplaceAllString(line, "")
		if len(line) > 0 && line[0] == '`' {
			lines = append(lines, line)
		}
	}

	out := make([]reflect.Type, 0, len(lines))
	for _, line := range lines {
		line = regexp.MustCompile(`,$`).ReplaceAllString(line, "")
		m := columnLineRE.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("parse column line %q", line)
		}
		name, chType := m[1], m[2]
		scan, err := clickHouseTypeToScanType(chType)
		if err != nil {
			return nil, fmt.Errorf("column %q: %w", name, err)
		}
		out = append(out, scan)
	}

	return out, nil
}

func clickHouseTypeToScanType(chType string) (reflect.Type, error) {
	switch {
	case chType == "LowCardinality(String)", chType == "String":
		return reflect.TypeFor[string](), nil
	case chType == "DateTime":
		return reflect.TypeFor[time.Time](), nil
	case chType == "UInt32":
		return reflect.TypeFor[uint32](), nil
	case chType == "UInt8":
		return reflect.TypeFor[uint8](), nil
	case chType == "Float32":
		return reflect.TypeFor[float32](), nil
	case chType == "Array(LowCardinality(String))", chType == "Array(String)":
		return reflect.TypeFor[[]string](), nil
	case chType == "Array(UInt32)":
		return reflect.TypeFor[[]uint32](), nil
	case chType == "Array(Float32)":
		return reflect.TypeFor[[]float32](), nil
	case chType == "Array(UInt64)":
		return reflect.TypeFor[[]uint64](), nil
	case len(chType) >= 5 && chType[:5] == "Enum8":
		// pmm-dump TSV stores enum labels as strings.
		return reflect.TypeFor[string](), nil
	default:
		return nil, fmt.Errorf("unsupported ClickHouse type %q", chType)
	}
}
