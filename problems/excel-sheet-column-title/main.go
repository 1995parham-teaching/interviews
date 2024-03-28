package main

import "fmt"

func convertToTitle(columnNumber int) string {
	if columnNumber == 0 {
		return ""
	}
	r := columnNumber % 26
	r = (r + 25) % 26
	return convertToTitle((columnNumber-1)/26) + string(byte('A'+r))
}

func main() {
	tests := []struct {
		columnNumber int
		columnName   string
	}{
		{
			columnNumber: 701,
			columnName:   "ZY",
		},
		{
			columnNumber: 1,
			columnName:   "A",
		},
		{
			columnNumber: 26,
			columnName:   "Z",
		},
		{
			columnNumber: 27,
			columnName:   "AA",
		},
	}

	for _, t := range tests {
		if t.columnName != convertToTitle(t.columnNumber) {
			fmt.Printf("Column number %d should be converted to %s but it converted to %s",
				t.columnNumber,
				t.columnName,
				convertToTitle(t.columnNumber),
			)
		}
	}
}
