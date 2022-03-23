package users

import "testing"

func Test_EmailToName(t *testing.T) {
	testCases := map[string]string{
		"nikita@getsturdy.com":          "Nikita",
		"gustav.westling@getsturdy.com": "Gustav Westling",
		"pro_HaCkEr123@gmail.com":       "Pro Hacker",
		"i2like3digits@gmail.com":       "I Like Digits",
		"sturdy@iamverysmart.com":       "Iamverysmart",
		"kiril.v.99@getsturdy.com":      "Kiril V",
	}
	for in, out := range testCases {
		t.Run(in, func(t *testing.T) {
			if got := EmailToName(in); got != out {
				t.Errorf("EmailToName(%q) = %q, want %q", in, got, out)
			}
		})
	}
}
