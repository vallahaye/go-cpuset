package cpuset

import (
	"fmt"
	"maps"
)

// A CPUSet defines a list of CPUs to restrict processes to.
// The zero value of a CPUSet is ready to use.
type CPUSet struct {
	m map[uint]struct{}
}

// Of returns a new [CPUSet] containing the CPUs listed.
func Of(cpus ...uint) CPUSet {
	s := CPUSet{
		m: make(map[uint]struct{}),
	}

	for _, cpu := range cpus {
		s.m[cpu] = struct{}{}
	}

	return s
}

// Add adds a CPU to s.
// It reports whether the CPU was not present before.
func (s *CPUSet) Add(cpu uint) bool {
	if s.m == nil {
		s.m = make(map[uint]struct{})
	}

	n := len(s.m)
	s.m[cpu] = struct{}{}
	return len(s.m) > n
}

// Delete removes a CPU from s.
// It reports whether the CPU was present.
func (s *CPUSet) Delete(cpu uint) bool {
	n := len(s.m)
	delete(s.m, cpu)
	return len(s.m) < n
}

// Contains reports whether a CPU is present in s.
func (s *CPUSet) Contains(cpu uint) bool {
	_, ok := s.m[cpu]
	return ok
}

// UnsortedList returns a slice of all the CPUs in s, in an unpredictable
// order.
func (s *CPUSet) UnsortedList() []uint {
	cpus := make([]uint, len(s.m))
	i := 0
	for cpu, _ := range s.m {
		cpus[i] = cpu
		i++
	}

	return cpus
}

// Equal reports whether s and s2 contain exactly the same CPUs.
func (s *CPUSet) Equal(s2 CPUSet) bool {
	if len(s.m) != len(s2.m) {
		return false
	}

	for cpu := range s2.m {
		if _, ok := s.m[cpu]; !ok {
			return false
		}
	}

	return true
}

// Clear removes all CPUs from s, leaving it empty.
func (s *CPUSet) Clear() {
	clear(s.m)
}

// Clone returns a copy of s.
func (s *CPUSet) Clone() CPUSet {
	return CPUSet{
		m: maps.Clone(s.m),
	}
}

// Len returns the number of CPUs in s.
func (s *CPUSet) Len() int {
	return len(s.m)
}

// String is an alias for [CPUSet.ListString].
func (s *CPUSet) String() string {
	return s.ListString()
}

// Difference returns a new [CPUSet] containing the CPUs of s1 that are not
// in s2.
func Difference(s1, s2 CPUSet) CPUSet {
	var s CPUSet
	for cpu := range s1.m {
		if _, ok := s2.m[cpu]; !ok {
			s.Add(cpu)
		}
	}

	return s
}

// Intersection returns a new [CPUSet] containing the CPUs of s1 that are
// in s2.
func Intersection(s1, s2 CPUSet) CPUSet {
	var s CPUSet
	m1, m2 := s1.m, s2.m
	if len(m1) > len(m2) {
		m2, m1 = m1, m2
	}

	for cpu := range m1 {
		if _, ok := m2[cpu]; ok {
			s.Add(cpu)
		}
	}

	return s
}

// Union returns a new [CPUSet] containing the CPUs of s1 and s2.
func Union(s1, s2 CPUSet) CPUSet {
	var s CPUSet
	for cpu := range s1.m {
		s.Add(cpu)
	}

	for cpu := range s2.m {
		s.Add(cpu)
	}

	return s
}

func formatParseError(s string, text string) error {
	return fmt.Errorf("cpuset: parsing %q: %s", s, text)
}
