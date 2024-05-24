package main
type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}