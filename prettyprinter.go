package prettyprinter

import "os"

type PrettyPrinter struct {
	spool  interface{}
	stderr Writer
	stdout Writer
	err    error
}

func NewPrettyPrinter() *PrettyPrinter {
	return NewPrettyPrinterFromRaw(os.Stderr, os.Stdout)
}

func NewPrettyPrinterFromRaw(
	stderr, stdout Writer,
) *PrettyPrinter {
	return &PrettyPrinter{nil, stderr, stdout, nil}
}

func (p *PrettyPrinter) Add(
	value interface{},
) Printer {
	p.spool = value
	return p
}

func (p *PrettyPrinter) StderrDump() Printer {
	return p.Dump(p.stderr)
}

func (p *PrettyPrinter) StdoutDump() Printer {
	return p.Dump(p.stdout)
}

func (p *PrettyPrinter) Dump(w Writer) Printer {
	bytes, err := prettyJSON(p.spool)
	p.Flush()
	if err != nil {
		p.err = err
		return p
	}
	p.err = print(bytes, w)
	return p
}

func (p *PrettyPrinter) Flush() {
	p.spool = nil
}

func (p *PrettyPrinter) Error() error {
	return p.err
}

func (p *PrettyPrinter) StderrDumpOnError() error {
	if p.Error() != nil {
		p.Flush()
		kve := MakeKeyValueError(p.Error())
		p.Add(kve)
		return p.StderrDump().Error()
	}
	return nil
}
