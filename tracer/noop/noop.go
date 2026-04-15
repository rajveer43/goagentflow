package noop

type Tracer struct{}

func New() Tracer { return Tracer{} }

func (Tracer) StartSpan(string) func() { return func() {} }
