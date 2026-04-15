package goagentflow

import "goagentflow/internal/stream"

type streamSink struct {
	stream *stream.EventStream
}

func (s *streamSink) Emit(event RuntimeEvent) {
	s.stream.TryEmit(event)
}

