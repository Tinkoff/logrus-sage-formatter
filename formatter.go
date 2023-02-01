package sageformatter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

const (
	SageSeverityTrace = "TRACE"
	SageSeverityDebug = "DEBUG"
	SageSeverityFatal = "FATAL"
	SageSeverityError = "ERROR"
	SageSeverityWarn  = "WARN"
	SageSeverityInfo  = "INFO"
)

type EnvironmentMetadata struct {
	// DC is datacenter
	DC string `json:"dc,omitempty" env:"SAGE_DC"`
	// Env is a runtime environment
	Env string `json:"env,omitempty" env:"SAGE_ENV"`
	// Group is a name of project, team, business line or similar
	Group string `json:"group,omitempty" env:"SAGE_GROUP"`
	// System is a job or application name
	System string `json:"system,omitempty" env:"SAGE_SYSTEM"`
	// Inst is an instance ID, can be hostname or Kubernetes pod ID
	Inst string `json:"inst,omitempty" env:"SAGE_INST"`
	// TimeFormat is format for @timestamp
	TimeFormat string `json:"time_format,omitempty" env:"SAGE_TIME_FORMAT"`
}

type Formatter struct {
	EnvironmentMetadata
}

type SageLogEntry struct {
	// Time is datetime in ISO8601 with timezone
	Time string `json:"@timestamp"`
	// Msg is a log entry main message
	Msg string `json:"msg"`
	// Level severity level of the entry
	Level string `json:"level"`
	// Extra fields given with WithFields
	Extra map[string]interface{} `json:"extra,omitempty"`
	EnvironmentMetadata
}

type FormatterOption = func(formatter *Formatter) error

func MetadataFromEnv(formatter *Formatter) error {
	if err := env.Parse(&formatter.EnvironmentMetadata); err != nil {
		return err
	}

	return nil
}

func Metadata(metadata EnvironmentMetadata) FormatterOption {
	return func(formatter *Formatter) error {
		formatter.EnvironmentMetadata = metadata

		return nil
	}
}

func NewFormatter(opts ...FormatterOption) *Formatter {
	f := Formatter{}

	for _, apply := range opts {
		if err := apply(&f); err != nil {
			panic(err)
		}
	}

	return &f
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	metadata := EnvironmentMetadata{
		Env:    f.Env,
		Inst:   f.Inst,
		DC:     f.DC,
		Group:  f.Group,
		System: f.System,
	}

	data := make(logrus.Fields, len(entry.Data))

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	timeFormat := time.RFC3339
	if f.TimeFormat != "" {
		timeFormat = f.TimeFormat
	}

	sageEntry := SageLogEntry{
		EnvironmentMetadata: metadata,
		Time:                entry.Time.UTC().Format(timeFormat),
		Msg:                 entry.Message,
		Level:               toSageLevel(entry.Level),
		Extra:               data,
	}

	serialized, err := json.Marshal(sageEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return append(serialized, '\n'), nil
}

func toSageLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return SageSeverityTrace
	case logrus.DebugLevel:
		return SageSeverityDebug
	case logrus.InfoLevel:
		return SageSeverityInfo
	case logrus.WarnLevel:
		return SageSeverityWarn
	case logrus.ErrorLevel:
		return SageSeverityError
	case logrus.FatalLevel:
		return SageSeverityFatal
	case logrus.PanicLevel:
		return SageSeverityFatal
	}

	return SageSeverityInfo
}
