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
	"encoding/csv"
	"fmt"
	"os"
)

func DetectFromTSV(path string) (Version, error) {
	f, err := os.Open(path)
	if err != nil {
		return VersionUnknown, fmt.Errorf("open TSV: %w", err)
	}
	defer f.Close() //nolint:errcheck

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 0

	record, err := r.Read()
	if err != nil {
		return VersionUnknown, fmt.Errorf("read first TSV row: %w", err)
	}

	return VersionFromFieldCount(len(record))
}

func ParseVersionFlag(v string) (Version, error) {
	switch v {
	case "", "auto":
		return VersionUnknown, nil
	case "pmm2", "2":
		return PMM2, nil
	case "pmm3", "3":
		return PMM3, nil
	default:
		return VersionUnknown, fmt.Errorf("invalid --pmm-version %q (use auto, pmm2, or pmm3)", v)
	}
}
