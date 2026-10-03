package genopenapi

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestSplitSummaryDescription(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		in          string
		wantSummary string
		wantDesc    string
	}{
		{
			name: "empty",
		},
		{
			name:        "summary only",
			in:          "Creates a book.",
			wantSummary: "Creates a book.",
		},
		{
			name:        "summary and description",
			in:          "Creates a book.\n\nLong description.\nSecond line.",
			wantSummary: "Creates a book.",
			wantDesc:    "Long description.\nSecond line.",
		},
		{
			name:        "leading and trailing whitespace is trimmed",
			in:          "  Summary.  \n\n  Body.  \n",
			wantSummary: "Summary.",
			wantDesc:    "Body.",
		},
		{
			// Only a blank line (\n\n) starts the description. A hard line
			// break inside the summary should stay part of the summary.
			name:        "single newline stays in summary",
			in:          "Line one\nLine two",
			wantSummary: "Line one\nLine two",
		},
		{
			// SplitN splits at the first blank line only; the body keeps
			// any internal blank lines verbatim.
			name:        "first blank line wins",
			in:          "Summary.\n\nFirst para.\n\nSecond para.",
			wantSummary: "Summary.",
			wantDesc:    "First para.\n\nSecond para.",
		},
		{
			// A comment that's only whitespace still produces empty strings,
			// not a lone space for summary.
			name: "whitespace only",
			in:   "   \n\t  ",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			summary, desc := splitSummaryDescription(tc.in)
			if summary != tc.wantSummary {
				t.Errorf("summary = %q, want %q", summary, tc.wantSummary)
			}
			if desc != tc.wantDesc {
				t.Errorf("description = %q, want %q", desc, tc.wantDesc)
			}
		})
	}
}

func TestExtractCommentsLineEndings(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		leading string
	}{
		{
			name:    "LF",
			leading: " Creates a book.\n\n Long description.\n Second line.\n",
		},
		{
			name:    "CRLF",
			leading: " Creates a book.\r\n\r\n Long description.\r\n Second line.\r\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := []int32{serviceProtoPath, 0, methodProtoPath, 0}
			file := &descriptor.File{
				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{
						Location: []*descriptorpb.SourceCodeInfo_Location{
							{Path: path, LeadingComments: &tc.leading},
						},
					},
				},
			}
			summary, desc := splitSummaryDescription(extractComments(file, path))
			if want := "Creates a book."; summary != want {
				t.Errorf("summary = %q, want %q", summary, want)
			}
			if want := "Long description.\nSecond line."; desc != want {
				t.Errorf("description = %q, want %q", desc, want)
			}
		})
	}
}
