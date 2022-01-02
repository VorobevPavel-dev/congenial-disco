package utility

import (
	"strings"
	"testing"
)

func TestSliceMethods(t *testing.T) {
	t.Run("Test StringIsIn method", func(t *testing.T) {
		slice := []string{"test", "test1", "test2"}
		elementIn := "test"
		elementOut := "test4"
		if !StringIsIn(elementIn, slice) {
			t.Errorf("String was not found in slice: element: %v, slice: %v",
				elementIn, slice)
		}
		if StringIsIn(elementOut, slice) {
			t.Errorf("String was found in slice: element: %v, slice: %v",
				elementOut, slice)
		}
	})
	t.Run("Test StringSliceEqual method", func(t *testing.T) {
		slice := []string{"", "", "1"}
		otherSlice := []string{"", "", "1"}
		if !StringSliceEqual(slice, otherSlice) {
			t.Errorf("Incoorect responce. Slices are equal")
		}
	})
}

func TestStringMethods(t *testing.T) {
	t.Run("Test DivideBySeparators method", func(t *testing.T) {
		inputs := []string{
			"this",
			"this(",
			"this(is",
			"this(is ",
			"this(is test",
			"this(is test)",
		}
		outputs := [][]string{
			{
				"this",
			},
			{
				"this", "(",
			},
			{
				"this", "(", "is",
			},
			{
				"this", "(", "is", " ",
			},
			{
				"this", "(", "is", " ", "test",
			},
			{
				"this", "(", "is", " ", "test", ")",
			},
		}
		separators := []string{"(", ")", " "}
		for testIndex := range inputs {
			actualResult := DivideBySeparators(inputs[testIndex], separators)
			if !StringSliceEqual(outputs[testIndex], actualResult) {
				t.Errorf("Error on test case #%d: expected:%s (size:%d), got:%s (size: %d) from line \"%s\"",
					testIndex,
					PrettyPrintStringsSlice(outputs[testIndex]),
					len(outputs[testIndex]),
					PrettyPrintStringsSlice(actualResult),
					len(actualResult),
					inputs[testIndex],
				)
			}
		}
	})
	t.Run("Test FindFirstOf method", func(t *testing.T) {
		input := "this(is test)"
		separators := []string{"(", ")", " ", "}"}
		for _, element := range separators {
			expectedPosition := strings.Index(input, element)
			actualPosition, _ := FindFirstOf(input, []string{element})
			if expectedPosition != actualPosition {
				t.Errorf("Incorrect responce from method: expected: %d, got: %d, line: %s, element: %s",
					expectedPosition, actualPosition, input, element)
			}
		}
	})
}
