package git

// PushNotes will push the requested note to the remote.
func PushNotes(remote string, directory string, ref string) error {
	args := []string{
		"push",
		remote,
		"refs/notes/" + ref,
	}
	_, err := execute(directory, args)
	return err
}
