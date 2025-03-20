package cpuset

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
	"strings"
)

const (
	partBase    = 10
	partBitSize = bits.UintSize
)

// ParseList decodes s into a [CPUSet]. It returns an error if s is not a valid
// list string, as specified in the [Linux cpuset(7) man page] (see "List
// Format").
//
// [Linux cpuset(7) man page]: https://man7.org/linux/man-pages/man7/cpuset.7.html
func ParseList(s string) (cset CPUSet, _ error) {
	if s == "" {
		return CPUSet{}, nil
	}

	for _, elem := range strings.Split(s, ",") {
		switch parts := strings.Split(elem, "-"); len(parts) {
		case 1:
			part, exclude := strings.CutPrefix(parts[0], "^")

			ui64, err := strconv.ParseUint(part, partBase, partBitSize)
			if err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid element %q", elem))
			}

			if cpu := uint(ui64); exclude {
				cset.Delete(cpu)
			} else {
				cset.Add(cpu)
			}

		case 2:
			ui64, err := strconv.ParseUint(parts[0], partBase, partBitSize)
			if err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid lower bound %q in range %q", parts[0], elem))
			}

			lowerBound := uint(ui64)

			ui64, err = strconv.ParseUint(parts[1], partBase, partBitSize)
			if err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid upper bound %q in range %q", parts[1], elem))
			}

			upperBound := uint(ui64)

			if upperBound < lowerBound {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("negative range %q", elem))
			}

			for i := range upperBound - lowerBound + 1 {
				cset.Add(lowerBound + i)
			}

		default:
			return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid element %q", elem))
		}
	}

	return cset, nil
}

// ListString encodes s into a list string.
func (s *CPUSet) ListString() string {
	if s.Len() == 0 {
		return ""
	}

	cpus := s.UnsortedList()
	slices.Sort(cpus)

	var (
		elems      []string
		lowerBound = cpus[0]
		upperBound = cpus[0]
	)

	for _, cpu := range cpus[1:] {
		if cpu != upperBound+1 {
			elems = append(elems, formatListElem(lowerBound, upperBound))
			lowerBound = cpu
		}

		upperBound = cpu
	}

	elems = append(elems, formatListElem(lowerBound, upperBound))

	return strings.Join(elems, ",")
}

func formatListElem(lowerBound, upperBound uint) string {
	if lowerBound == upperBound {
		return fmt.Sprint(lowerBound)
	}

	return fmt.Sprint(lowerBound, "-", upperBound)
}
