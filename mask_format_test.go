package cpuset

import (
	"testing"
)

func TestParseMask(t *testing.T) {
	for _, params := range []struct {
		name string
		s    string
		want CPUSet
		err  error
	}{
		{
			name: `invalid 32-bit word ""`,
			s:    ",",
			err:  formatParseError(",", `invalid 32-bit word ""`),
		},
		{
			name: `invalid 32-bit word "100000000"`,
			s:    "100000000",
			err:  formatParseError("100000000", `invalid 32-bit word "100000000"`),
		},
		{
			name: `invalid 32-bit word "xxxxxxxx"`,
			s:    "xxxxxxxx",
			err:  formatParseError("xxxxxxxx", `invalid 32-bit word "xxxxxxxx"`),
		},
		{
			name: "no bit set",
			s:    "",
			want: CPUSet{},
		},
		{
			name: "bit 0 set",
			s:    "00000001",
			want: Of(0),
		},
		{
			name: "bit 94 set",
			s:    "40000000,00000000,00000000",
			want: Of(94),
		},
		{
			name: "bit 64 set",
			s:    "00000001,00000000,00000000",
			want: Of(64),
		},
		{
			name: "bits 32, 33, 34, 35, 36, 37, 38, and 39 set",
			s:    "000000ff,00000000",
			want: Of(32, 33, 34, 35, 36, 37, 38, 39),
		},
		{
			name: "bits 1, 5, 6, 11, 12, 13, 17, 18, and 19 set",
			s:    "000e3862",
			want: Of(1, 5, 6, 11, 12, 13, 17, 18, 19),
		},
		{
			name: "bits 0, 1, 2, 4, 8, 16, 32, and 64 set",
			s:    "00000001,00000001,00010117",
			want: Of(0, 1, 2, 4, 8, 16, 32, 64),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			switch got, err := ParseMask(params.s); {
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

func TestMaskString(t *testing.T) {
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
			name: "bit 0 set",
			s:    Of(0),
			want: "00000001",
		},
		{
			name: "bit 94 set",
			s:    Of(94),
			want: "40000000,00000000,00000000",
		},
		{
			name: "bit 64 set",
			s:    Of(64),
			want: "00000001,00000000,00000000",
		},
		{
			name: "bits 32, 33, 34, 35, 36, 37, 38, and 39 set",
			s:    Of(32, 33, 34, 35, 36, 37, 38, 39),
			want: "000000ff,00000000",
		},
		{
			name: "bits 1, 5, 6, 11, 12, 13, 17, 18, and 19 set",
			s:    Of(1, 5, 6, 11, 12, 13, 17, 18, 19),
			want: "000e3862",
		},
		{
			name: "bits 0, 1, 2, 4, 8, 16, 32, and 64 set",
			s:    Of(0, 1, 2, 4, 8, 16, 32, 64),
			want: "00000001,00000001,00010117",
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.MaskString(); got != params.want {
				t.Errorf("unexpected mask string: got %q, want %q", got, params.want)
			}
		})
	}
}
