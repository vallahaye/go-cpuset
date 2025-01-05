package cpuset

import (
	"fmt"
	"slices"
	"strings"
)

const wordBitSize = 32

// ParseMask decodes s into a [CPUSet]. It returns an error if s is not a valid
// mask string, as specified in the [Linux cpuset(7) man page] (see "Mask
// Format").
//
// [Linux cpuset(7) man page]: https://man7.org/linux/man-pages/man7/cpuset.7.html
func ParseMask(s string) (CPUSet, error) {
	if s == "" {
		return CPUSet{}, nil
	}

	var cpus []uint
	words := strings.Split(s, ",")
	for i, word := range words {
		var cpuMask uint32
		if _, err := fmt.Sscanf(word, "%x", &cpuMask); err != nil {
			return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid 32-bit word %q", word))
		}

		offset := (len(words) - i - 1) * wordBitSize
		for pos := range wordBitSize {
			if cpuMask&(1<<pos) != 0 {
				cpus = append(cpus, uint(pos+offset))
			}
		}
	}

	return Of(cpus...), nil
}

// MaskString encodes s into a mask string.
func (s *CPUSet) MaskString() string {
	if len(s.m) == 0 {
		return ""
	}

	cpus := s.UnsortedList()
	cpuMasks := make([]uint32, slices.Max(cpus)/wordBitSize+1)
	for _, cpu := range cpus {
		i, pos := int(cpu/wordBitSize), cpu%wordBitSize
		cpuMasks[len(cpuMasks)-i-1] |= 1 << pos
	}

	words := make([]string, len(cpuMasks))
	for i, cpuMask := range cpuMasks {
		words[i] = fmt.Sprintf("%08x", cpuMask)
	}

	return strings.Join(words, ",")
}
