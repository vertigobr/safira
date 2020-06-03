package stack

func GetAllFunctions() (map[string]Function, error) {
	stack, err := LoadStackFile(stackFileName)
	if err != nil {
		return nil, err
	}

	return stack.Functions, nil
}
