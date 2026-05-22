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

import "fmt"

type Version int

const (
	VersionUnknown Version = 0
	PMM2           Version = 2
	PMM3           Version = 3
)

const (
	PMM2ColumnCount = 228
	PMM3ColumnCount = 269
)

func (v Version) String() string {
	switch v {
	case PMM2:
		return "pmm2"
	case PMM3:
		return "pmm3"
	default:
		return "unknown"
	}
}

func VersionFromFieldCount(n int) (Version, error) {
	switch n {
	case PMM2ColumnCount:
		return PMM2, nil
	case PMM3ColumnCount:
		return PMM3, nil
	default:
		return VersionUnknown, fmt.Errorf("unsupported TSV width %d columns (expected %d for PMM2 or %d for PMM3)", n, PMM2ColumnCount, PMM3ColumnCount)
	}
}
