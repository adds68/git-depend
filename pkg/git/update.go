package git

func Update(url string, directory string) error {
	args := []string{"fetch"}
	_, err := execute(directory, args)
	return err
}
