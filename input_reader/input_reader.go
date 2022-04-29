package input_reader

type InputReaderData struct {
	eventType string
	args []string
}

type InputReader interface {
	Read() InputReaderData
}

func NewInputReader() InputReader {
	return nil
}
