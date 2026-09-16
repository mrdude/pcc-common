package wideevt

import (
	"encoding/json"

	"go.uber.org/zap"
)

type Sink interface {
	Write(evt *Event) error
}

type logSink struct {
	logger *zap.Logger
}

func NewLogSink(logger *zap.Logger) Sink {
	return &logSink{logger: logger}
}

func (s *logSink) Write(evt *Event) error {
	var fields []zap.Field

	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}

	s.logger.Info("WideEvent", fields...)
	return nil
}
