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
	"io"
	"reflect"

	"pmm-dump-qan-sql/internal/qantsv"
)

type Options struct {
	Database  string
	Table     string
	BatchSize int
}

func Generate(r io.Reader, w io.Writer, scanTypes []reflect.Type, opts Options) (int, error) {
	if opts.BatchSize <= 0 {
		return 0, fmt.Errorf("batch size must be positive, got %d", opts.BatchSize)
	}
	if opts.Table == "" {
		return 0, fmt.Errorf("table name is required")
	}
	if len(scanTypes) == 0 {
		return 0, fmt.Errorf("column scan types are required")
	}

	reader := qantsv.NewReader(r, scanTypes)
	batch := make([]string, 0, opts.BatchSize)
	rows := 0

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := writeInsert(w, opts, batch); err != nil {
			return err
		}
		batch = batch[:0]
		return nil
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return rows, fmt.Errorf("read TSV row %d: %w", rows+1, err)
		}

		tuple, err := formatRow(record, scanTypes)
		if err != nil {
			return rows, fmt.Errorf("format row %d: %w", rows+1, err)
		}

		batch = append(batch, tuple)
		rows++

		if len(batch) >= opts.BatchSize {
			if err := flush(); err != nil {
				return rows, err
			}
		}
	}

	if err := flush(); err != nil {
		return rows, err
	}

	return rows, nil
}

func writeInsert(w io.Writer, opts Options, tuples []string) error {
	target := opts.Table
	if opts.Database != "" {
		target = opts.Database + "." + opts.Table
	}

	if _, err := fmt.Fprintf(w, "INSERT INTO %s VALUES\n", target); err != nil {
		return err
	}

	for i, tuple := range tuples {
		sep := ",\n"
		if i == len(tuples)-1 {
			sep = "\n"
		}
		if _, err := fmt.Fprint(w, tuple, sep); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, ";\n\n"); err != nil {
		return err
	}

	return nil
}
