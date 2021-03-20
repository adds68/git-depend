package git

func Push(url string, directory string) error {
	args := []string{
		"push",
		url,
	}
	_, err := execute(directory, args)
	return err
}
