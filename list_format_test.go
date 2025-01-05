package cpuset

import (
	"testing"
)

func TestParseList(t *testing.T) {
	for _, params := range []struct {
		name string
		s    string
		want CPUSet
		err  error
	}{
		{
			name: `invalid element ""`,
			s:    ",",
			err:  formatParseError(",", `invalid element ""`),
		},
		{
			name: `invalid element "a"`,
			s:    "a",
			err:  formatParseError("a", `invalid element "a"`),
		},
		{
			name: `invalid lower bound "" in range "-1"`,
			s:    "-1",
			err:  formatParseError("-1", `invalid lower bound "" in range "-1"`),
		},
		{
			name: `invalid lower bound "a" in range "a-1"`,
			s:    "a-1",
			err:  formatParseError("a-1", `invalid lower bound "a" in range "a-1"`),
		},
		{
			name: `invalid upper bound "" in range "1-"`,
			s:    "1-",
			err:  formatParseError("1-", `invalid upper bound "" in range "1-"`),
		},
		{
			name: `invalid upper bound "b" in range "1-b"`,
			s:    "1-b",
			err:  formatParseError("1-b", `invalid upper bound "b" in range "1-b"`),
		},
		{
			name: `negative range "1-0"`,
			s:    "1-0",
			err:  formatParseError("1-0", `negative range "1-0"`),
		},
		{
			name: `invalid element "0-1-2"`,
			s:    "0-1-2",
			err:  formatParseError("0-1-2", `invalid element "0-1-2"`),
		},
		{
			name: "no bit set",
			s:    "",
			want: CPUSet{},
		},
		{
			name: "bits 0, 1, 2, 3, 4, and 9 set",
			s:    "0-4,9",
			want: Of(0, 1, 2, 3, 4, 9),
		},
		{
			name: "bits 0, 1, 2, 7, 12, 13, and 14 set",
			s:    "0-2,7,12-14",
			want: Of(0, 1, 2, 7, 12, 13, 14),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			switch got, err := ParseList(params.s); {
			case err == nil && params.err != nil:
				t.Error("expected error")
			case err != nil && params.err == nil:
				t.Errorf("unexpected error: %v", err)
			case err != nil && params.err != nil && err.Error() != params.err.Error():
				t.Errorf("unexpected error: got %v, want %v", err, params.err)
			case err == nil && params.err == nil && !got.Equal(params.want):
				t.Errorf("unexpected cpuset: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestListString(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		want string
	}{
		{
			name: "no bit set",
			s:    CPUSet{},
			want: "",
		},
		{
			name: "bits 0, 1, 2, 3, 4, and 9 set",
			s:    Of(0, 1, 2, 3, 4, 9),
			want: "0-4,9",
		},
		{
			name: "bits 0, 1, 2, 7, 12, 13, and 14 set",
			s:    Of(0, 1, 2, 7, 12, 13, 14),
			want: "0-2,7,12-14",
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.ListString(); got != params.want {
				t.Errorf("unexpected list string: got %q, want %q", got, params.want)
			}
		})
	}
}
