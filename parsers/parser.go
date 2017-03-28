package parsers

type Parser interface {
	Marshal([]byte, int) []byte
	Unmarshal([]byte, int) (int, []byte, error)
}
