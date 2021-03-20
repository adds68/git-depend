package git

func Clone(url string, directory string) error {
	args := []string{"clone", url, directory}
	_, err := execute("", args)
	return err
}
