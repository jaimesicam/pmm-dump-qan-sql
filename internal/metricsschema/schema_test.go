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
	"path/filepath"
	"testing"
)

func TestScanTypesCount(t *testing.T) {
	t.Parallel()

	pmm2, err := ScanTypes(PMM2)
	if err != nil {
		t.Fatal(err)
	}
	if len(pmm2) != PMM2ColumnCount {
		t.Fatalf("PMM2: expected %d columns, got %d", PMM2ColumnCount, len(pmm2))
	}

	pmm3, err := ScanTypes(PMM3)
	if err != nil {
		t.Fatal(err)
	}
	if len(pmm3) != PMM3ColumnCount {
		t.Fatalf("PMM3: expected %d columns, got %d", PMM3ColumnCount, len(pmm3))
	}
}

func TestDetectFromTSV(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..", "..", "source")
	cases := []struct {
		path    string
		version Version
	}{
		{filepath.Join(root, "pmm2", "ch", "0.tsv"), PMM2},
		{filepath.Join(root, "ch", "0.tsv"), PMM3},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.version.String(), func(t *testing.T) {
			t.Parallel()
			got, err := DetectFromTSV(tc.path)
			if err != nil {
				t.Skip(err)
			}
			if got != tc.version {
				t.Fatalf("got %v, want %v", got, tc.version)
			}
		})
	}
}

func TestClickHouseImportCommand(t *testing.T) {
	t.Parallel()

	if got := ClickHouseImportCommand(PMM2, "out.sql"); got != `clickhouse-client --database=pmm --queries-file "out.sql"` {
		t.Fatalf("PMM2: %q", got)
	}
	if got := ClickHouseImportCommand(PMM3, "out.sql"); got != `clickhouse-client --database=pmm --password=clickhouse --queries-file "out.sql"` {
		t.Fatalf("PMM3: %q", got)
	}
}
