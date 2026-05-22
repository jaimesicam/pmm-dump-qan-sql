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

import "strconv"

// ClickHouseImportCommand returns a clickhouse-client command to load SQL on the PMM server.
func ClickHouseImportCommand(version Version, queriesFile string) string {
	file := strconv.Quote(queriesFile)
	switch version {
	case PMM2:
		return "clickhouse-client --database=pmm --queries-file " + file
	case PMM3:
		return "clickhouse-client --database=pmm --password=clickhouse --queries-file " + file
	default:
		return ""
	}
}
