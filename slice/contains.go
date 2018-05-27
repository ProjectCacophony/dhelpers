package slice

import "strings"

// ContainsLowerExclude returns true if args contains a string in include, false if not
// newArgs will return args without the include string if there is a match, or the original slice if not
// each match is case insensitive
func ContainsLowerExclude(args, include []string) (includes bool, newArgs []string) {
	for i, arg := range args {
		arg = strings.ToLower(arg)
		for _, includeItem := range include {
			if strings.ToLower(includeItem) == arg {
				newArgs = append(args[:i], args[i+1:]...)
				return true, newArgs
			}
		}
	}
	return false, args
}
