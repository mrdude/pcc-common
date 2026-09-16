package wideevt

import (
	"encoding/json"
	"time"
)

type Event struct {
	eventData
	Ts       time.Time
	Duration time.Duration
	data     map[string]any
}

type eventData struct {
	// TODO add more standard fields
	ServerName string `json:"server"`

	TraceId      string `json:"traceId,omitempty"`
	SpanId       string `json:"spanId,omitempty"`
	ParentSpanId string `json:"parentSpanId,omitempty"`

	RequestMethod string `json:"method"`
	RequestPath   string `json:"path"`
	RequestURL    string `json:"url"`

	StatusCode int `json:"statusCode"`
}

func NewEvent() *Event {
	return &Event{
		Ts:   time.Now().UTC(),
		data: make(map[string]any),
	}
}

func (evt *Event) Update(key string, value any) {
	evt.data[key] = value
}

func (evt *Event) MarshalJSON() ([]byte, error) {
	toMap := func(obj any) (map[string]any, error) {
		b, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}

		var v map[string]any
		err = json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		return v, nil
	}

	// marshal extra fields
	fields, err := toMap(evt.data)
	if err != nil {
		return nil, err
	}

	// marshal standard fields
	stdFields, err := toMap(evt.eventData)
	if err != nil {
		return nil, err
	}

	stdFields["duration_ms"] = evt.Duration.Milliseconds()
	stdFields["ts"] = evt.Ts.Format(time.DateTime)

	// merge
	for k, v := range stdFields {
		fields[k] = v
	}

	return json.Marshal(fields)
}
