package output_test

import (
	"reflect"
	"testing"

	"github.com/sketch-hq/gql-lint/output"
)

func TestCompareFiles(t *testing.T) {
	type args struct {
		fileA string
		fileB string
	}
	tests := []struct {
		name    string
		args    args
		want    output.Data
		wantErr bool
	}{
		{
			name: "returns an error if file A has a parse error",
			args: args{
				fileA: "testdata/parse_error.json",
				fileB: "testdata/b.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns an error if file B has a parse error",
			args: args{
				fileA: "testdata/parse_error.json",
				fileB: "testdata/b.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns any new fields in `b`",
			args: args{
				fileA: "testdata/a.json",
				fileB: "testdata/b.json",
			},
			want: output.Data{
				output.Field{
					Field:             "Article.view",
					File:              "somefile.graphql",
					Line:              13,
					DeprecationReason: "Please migrate to Article.permissions",
				},
			},
			wantErr: false,
		},
		{
			name: "ignores any fields that was removed",
			args: args{
				fileA: "testdata/b.json",
				fileB: "testdata/a.json",
			},
			want:    output.Data{},
			wantErr: false,
		},
		{
			name: "returns an empty response if no changes are found",
			args: args{
				fileA: "testdata/a.json",
				fileB: "testdata/a.json",
			},
			want:    output.Data{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := output.CompareFiles(tt.args.fileA, tt.args.fileB)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
