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

import "testing"

func TestVersionFromFieldCount(t *testing.T) {
	t.Parallel()

	v, err := VersionFromFieldCount(228)
	if err != nil || v != PMM2 {
		t.Fatalf("228: got %v err %v", v, err)
	}

	v, err = VersionFromFieldCount(269)
	if err != nil || v != PMM3 {
		t.Fatalf("269: got %v err %v", v, err)
	}

	_, err = VersionFromFieldCount(100)
	if err == nil {
		t.Fatal("expected error for unknown width")
	}
}

func TestParseVersionFlag(t *testing.T) {
	t.Parallel()

	v, err := ParseVersionFlag("auto")
	if err != nil || v != VersionUnknown {
		t.Fatalf("auto: %v %v", v, err)
	}

	v, err = ParseVersionFlag("pmm2")
	if err != nil || v != PMM2 {
		t.Fatalf("pmm2: %v %v", v, err)
	}
}
