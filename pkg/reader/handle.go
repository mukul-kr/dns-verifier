package reader

type InputHandler interface {
	Handle(content string) func() ([]string, error)
}

func HandlerFactory(inputType string) InputHandler {
	switch inputType {
	case "csv":
		return CSVHandler{}
	case "json":
		return JSONHandler{}
	case "txt":
		return TXTHandler{}
	case "terminal":
		return TerminalHandler{}
	default:
		return nil
	}
}
