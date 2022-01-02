package utility

//DivideBySeparators will "split" string by separators, but will leave them in
// result slice
// Example:
//  DivideBySeparators("this(is test)", []string{"(",")"," "})
//	This will return []string{"this", "(", "is", " ", "test", ")"}
func DivideBySeparators(line string, separators []string) []string {
	var (
		result       []string
		character    string
		tempPosition = 0
	)
	lineLeft := line
	for len(line) != 0 {
		tempPosition, character = FindFirstOf(lineLeft, separators)
		//TODO: Fix this. Works correctly but ugly as f
		if tempPosition == -1 {
			if lineLeft != "" {
				result = append(result, lineLeft)
			}
			return result
		} else {
			if character != "" {
				result = append(result, lineLeft[:tempPosition], character)
			} else {
				result = append(result, lineLeft[:tempPosition])
			}
			lineLeft = lineLeft[tempPosition+1:]
			//result = append(result, string(lineLeft[0]))
		}
	}
	return nil
}

//FindFirstOf will return first position of one of elements provided
// in catch slice or -1 if there are no such elements in line
func FindFirstOf(line string, catch []string) (int, string) {
	var currentPosition int
	for currentPosition = 0; currentPosition < (len(line)); currentPosition++ {
		if StringIsIn(string(line[currentPosition]), catch) {
			return currentPosition, string(line[currentPosition])
		}
	}
	return -1, ""
}
