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
	"reflect"
	"testing"
	"time"
)

func TestQuoteString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"plain", "'plain'"},
		{"it's", `'it\'s'`},
		{"a\\b", `'a\\b'`},
		{"line\nbreak", `'line\nbreak'`},
	}

	for _, tt := range tests {
		if got := quoteString(tt.in); got != tt.want {
			t.Errorf("quoteString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFormatSlice(t *testing.T) {
	t.Parallel()

	got, err := formatSlice([]any{"a", "b"}, reflect.TypeFor[string]())
	if err != nil {
		t.Fatal(err)
	}
	if got != "['a', 'b']" {
		t.Fatalf("got %q", got)
	}

	got, err = formatSlice([]any{}, reflect.TypeFor[string]())
	if err != nil {
		t.Fatal(err)
	}
	if got != "[]" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatTime(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 5, 17, 2, 18, 32, 0, time.UTC)
	got, err := formatValue(ts, reflect.TypeFor[time.Time]())
	if err != nil {
		t.Fatal(err)
	}
	if got != "'2026-05-17 02:18:32'" {
		t.Fatalf("got %q", got)
	}
}
