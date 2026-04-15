package runtime

import "github.com/rajveer43/goagentflow/internal/stream"

type streamSink struct {
	stream *stream.EventStream
}

func (s *streamSink) Emit(event RuntimeEvent) {
	s.stream.TryEmit(event)
}

