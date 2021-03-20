package git

// PushNotes will push the requested note to the remote.
func PushNotes(remote string, directory string, ref string) error {
	args := []string{
		"push",
		remote,
		"refs/notes/", //TODO: Check these slashes are OK for Windows.
		ref,
	}
	_, err := execute(directory, args)
	return err
}
