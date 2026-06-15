// Command openapiv3-merge merges multiple OpenAPI 3.1 JSON documents into
// a single document and writes the result to stdout.
//
// The OpenAPI 3.1 spec describes one API as a single document, but
// protoc-gen-openapiv3 follows the protobuf convention of one output file
// per input file. openapiv3-merge bridges the two: pipe `protoc-gen-openapiv3`
// output through it to obtain a single combined spec.
//
// Usage:
//
//	openapiv3-merge FILE [FILE ...] > merged.openapi.json
//
// The merger is strict: path collisions, conflicting component definitions,
// and conflicting tag metadata are errors. `info`, `servers`, `externalDocs`,
// and unknown top-level fields are taken from the first input; later inputs'
// values for those fields are ignored. See the package documentation in
// internal/merge for the full ruleset.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/openapiv3-merge/internal/merge"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "openapiv3-merge:", err)
		os.Exit(1)
	}
}

// run is the testable entry point. It accepts the program's arguments (no
// program name) and the writer to emit the merged document to.
func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: openapiv3-merge FILE [FILE ...]")
	}
	inputs := make([]merge.Input, 0, len(args))
	for _, path := range args {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		inputs = append(inputs, merge.Input{Name: path, Data: data})
	}
	merged, err := merge.Merge(inputs)
	if err != nil {
		return err
	}
	if _, err := out.Write(merged); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n"); err != nil {
		return err
	}
	return nil
}
