package git

func Push(url string, directory string) error {
	args := []string{
		"push",
		url,
		"-C",
		directory,
	}
	_, r := execute(args)
	return r
}
