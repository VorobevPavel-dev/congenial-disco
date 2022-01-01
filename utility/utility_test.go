package utility

import "testing"

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
}
