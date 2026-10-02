package internal

import (
	"cmp"
	"flag"
	"reflect"
	"slices"
)

// GroupFlags groups flags by [flag.Value] equivalence and returns them.
//
// The flags are grouped by [flag.Value] equivalence. This allows flags to be grouped together in the rendered
// usage text when two flags are aliases of each other. This is often the case for short flags which are aliases of
// longer flags (e.g. '-a' is an alias of '--all').
//
//	-a <string>, --addr=<string>
//	-s <string>, --serial-number=<string>
//
// The resulting map entries are keyed by the flag group name, which is the longest flag name in the group. The map
// values are slices of (one or more) flags in the flag group, sorted by flag name length ('-a' before '--all').
//
// Hidden flags are excluded from the resulting map.
func GroupFlags(collected []*flag.Flag) map[string][]*flag.Flag {
	// sort flags by name length in descending order to ensure that keys in resulting map will use long names first
	slices.SortFunc(collected, func(a, b *flag.Flag) int {
		return cmp.Compare(len(b.Name), len(a.Name))
	})

	groups := map[string][]*flag.Flag{}

	for len(collected) > 0 {
		var flg *flag.Flag

		// pop the head of the slice
		flg, collected = collected[0], collected[1:]

		// update groups
		groups[flg.Name] = []*flag.Flag{flg}

		// traverse the flags again and find (and remove) any which match flg
		for i := range slices.Backward(collected) {
			other := collected[i]

			if areSame(flg.Value, other.Value) {
				groups[flg.Name] = append(groups[flg.Name], other)
				collected = append(collected[:i], collected[i+1:]...)
			}
		}

		// sort by length (then lexical order), this time ascending (-a before --all)
		slices.SortFunc(groups[flg.Name], func(a, b *flag.Flag) int {
			return cmp.Or(
				cmp.Compare(len(a.Name), len(b.Name)),
				cmp.Compare(a.Name, b.Name),
			)
		})
	}

	return groups
}

// areSame check if f1 and f2 have the same underlying [flag.Value].
func areSame(f1, f2 flag.Value) bool {
	var (
		ref1 = reflect.ValueOf(f1)
		ref2 = reflect.ValueOf(f2)
	)

	if ref1.Comparable() && ref2.Comparable() && f1 == f2 {
		return true
	}

	if ref1.Kind() != ref2.Kind() {
		return false
	}

	if !slices.Contains([]reflect.Kind{reflect.Map, reflect.Pointer, reflect.Func, reflect.Slice}, ref1.Kind()) {
		return false
	}

	return ref1.Pointer() == ref2.Pointer()
}
