package main

import "fmt"

func main() {
	fmt.Println(longestpalindrome("babad"))
}

func longestpalindrome(s string) string {
	lognest := 0
	sub := ""

	for i := 0; i < len(s); i++ {
		for j := i; j <= len(s); j++ {
			if j-i > lognest && isPalindrome(s[i:j]) {
				lognest = j - i
				sub = s[i:j]
			}
		}
	}

	return sub
}

func isPalindrome(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i := 0; i <= len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}
