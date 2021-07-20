package flagutil

import "flag"

// IsFlagPassed determines if a flag was explicitly passed in
func IsFlagPassed(name string) bool {
	if !flag.Parsed() {
		flag.Parse()
	}

	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
