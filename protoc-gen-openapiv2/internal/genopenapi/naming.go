package genopenapi

import (
	"reflect"
	"strings"
)

// LookupNamingStrategy looks up the given naming strategy and returns the naming
// strategy function for it. The naming strategy function takes in the list of all
// fully-qualified proto message names, and returns a mapping from fully-qualified
// name to OpenAPI name.
func LookupNamingStrategy(strategyName string) func([]string) map[string]string {
	switch strategyName {
	case "fqn":
		return resolveNamesFQN
	case "legacy":
		return resolveNamesLegacy
	case "simple":
		return resolveNamesSimple
	}
	return nil
}

// resolveNamesFQN uses the fully-qualified proto message name as the
// OpenAPI name, stripping the leading dot.
func resolveNamesFQN(messages []string) map[string]string {
	uniqueNames := make(map[string]string, len(messages))
	for _, p := range messages {
		// strip leading dot from proto fqn
		uniqueNames[p] = p[1:]
	}
	return uniqueNames
}

// resolveNamesLegacy takes the names of all protos and uniq-ifies them applying the legacy
// heuristics for deriving unique names: starting from the bottom of the name hierarchy, it
// determines the minimum number of components necessary to yield a unique name, adds one
// to that number, and then concatenates those last components with no separator in between
// to form a unique name.
//
// E.g., if the fully qualified name is `.a.b.C.D`, and there are other messages with fully
// qualified names ending in `.D` but not in `.C.D`, it assigns the unique name `bCD`.
func resolveNamesLegacy(messages []string) map[string]string {
	return resolveNamesUniqueWithContext(messages, 1, "")
}

// resolveNamesSimple takes the names of all protos and uniq-ifies them using a simple
// heuristic: starting from the bottom of the name hierarchy, it determines the minimum
// number of components necessary to yield a unique name, and then concatenates those last
// components with a "." separator in between to form a unique name.
//
// E.g., if the fully qualified name is `.a.b.C.D`, and there are other messages with
// fully qualified names ending in `.D` but not in `.C.D`, it assigns the unique name `C.D`.
func resolveNamesSimple(messages []string) map[string]string {
	return resolveNamesUniqueWithContext(messages, 0, ".")
}

// Take the names of every proto and "uniq-ify" them. The idea is to produce a
// set of names that meet a couple of conditions. They must be stable, they
// must be unique, and they must be shorter than the FQN.
//
// This likely could be made better. This will always generate the same names
// but may not always produce optimal names. This is a reasonably close
// approximation of what they should look like in most cases.
func resolveNamesUniqueWithContext(messages []string, extraContext int, componentSeparator string) map[string]string {
	packagesByDepth := make(map[int][][]string)
	uniqueNames := make(map[string]string)

	hierarchy := func(pkg string) []string {
		return strings.Split(pkg, ".")
	}

	for _, p := range messages {
		h := hierarchy(p)
		for depth := range h {
			if _, ok := packagesByDepth[depth]; !ok {
				packagesByDepth[depth] = make([][]string, 0)
			}
			packagesByDepth[depth] = append(packagesByDepth[depth], h[len(h)-depth:])
		}
	}

	count := func(list [][]string, item []string) int {
		i := 0
		for _, element := range list {
			if reflect.DeepEqual(element, item) {
				i++
			}
		}
		return i
	}

	for _, p := range messages {
		h := hierarchy(p)
		depth := 0
		for ; depth < len(h); depth++ {
			// depth + extraContext > 0 ensures that we only break for values of depth when the
			// resulting slice of name components is non-empty. Otherwise, we would return the
			// empty string as the concise unique name is len(messages) == 1 (which is
			// technically correct).
			if depth+extraContext > 0 && count(packagesByDepth[depth], h[len(h)-depth:]) == 1 {
				break
			}
		}
		start := len(h) - depth - extraContext
		if start < 0 {
			start = 0
		}
		uniqueNames[p] = strings.Join(h[start:], componentSeparator)
	}
	return uniqueNames
}
