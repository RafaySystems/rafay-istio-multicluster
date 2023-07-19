package utils

import "testing"

func TestStringSliceContainsDuplicates(t *testing.T) {
	type args struct {
		strSlice []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test without duplicates",
			args: args{
				strSlice: []string{"a", "b", "c"},
			},
			want:  "",
			want1: false,
		},
		{
			name: "Test without duplicates",
			args: args{
				strSlice: []string{"a"},
			},
			want:  "",
			want1: false,
		},
		{
			name: "Test with duplicates",
			args: args{
				strSlice: []string{"a", "a"},
			},
			want:  "a",
			want1: true,
		},
		{
			name: "Test with multiple duplicates",
			args: args{
				strSlice: []string{"a", "b", "a", "b"},
			},
			want:  "a",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StringSliceContainsDuplicates(tt.args.strSlice)
			if got != tt.want {
				t.Errorf("StringSliceContainsDuplicates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("StringSliceContainsDuplicates() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
