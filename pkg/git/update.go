package git

func Update(url string, directory string) error {
	args := []string{
		"-C",
		directory,
		"fetch",
	}
	_, err := execute(args)
	return err
}
