package prettyprinter

type Printer interface {
	Add(interface{}) Printer
	Dump(Writer) Printer
	StderrDump() Printer
	StdoutDump() Printer
	Error() error
	StderrDumpOnError() error
}

type Writer interface {
	Write([]byte) (int, error)
}
