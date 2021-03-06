package depend

// Lock allows us to safely write to a note.
type Lock struct {
	ID        string `json: "Id"`
	Timestamp string `json: "Timestamp"`
	Status    string `json: "Status"`
}
