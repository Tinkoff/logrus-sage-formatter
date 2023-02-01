package sageformatter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	sageformatter "github.com/Tinkoff/logrus-sage-formatter"
	"github.com/sirupsen/logrus"
)

func TestFormatter(t *testing.T) {
	for _, tt := range formatterTests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			var out bytes.Buffer

			logger := logrus.New()
			logger.Out = &out
			logger.SetFormatter(sageformatter.NewFormatter(sageformatter.Metadata(sageformatter.EnvironmentMetadata{
				Env:    "prod",
				Group:  "test",
				System: "system",
			})))

			mTime := time.Now()
			test.run(logger, mTime)
			m := map[string]interface{}{}
			if err := json.Unmarshal(out.Bytes(), &m); err != nil {
				t.Error(err)
			}

			matchOut := test.out
			matchOut["@timestamp"] = mTime.UTC().Format(time.RFC3339)

			if !reflect.DeepEqual(test.out, m) {
				correct, _ := json.MarshalIndent(&test.out, "", "  ")
				t.Log(out.String())
				t.Log("expected:")
				t.Log(string(correct))
				t.Error("invalid format")
			}
		})
	}
}

var formatterTests = []struct {
	run  func(*logrus.Logger, time.Time)
	out  map[string]interface{}
	name string
}{
	{
		name: "With Field",
		run: func(logger *logrus.Logger, t time.Time) {
			logger.WithTime(t).WithField("foo", "bar").Info("my log entry")
		},
		out: map[string]interface{}{
			"level": "INFO",
			"msg":   "my log entry",
			"extra": map[string]interface{}{
				"foo": "bar",
			},
			"env":    "prod",
			"group":  "test",
			"system": "system",
		},
	},
	{
		name: "With Field, skip empty msg",
		run: func(logger *logrus.Logger, t time.Time) {
			logger.WithTime(t).WithField("foo", "bar").Info()
		},
		out: map[string]interface{}{
			"level": "INFO",
			"msg":   "",
			"extra": map[string]interface{}{
				"foo": "bar",
			},
			"env":    "prod",
			"group":  "test",
			"system": "system",
		},
	},
	{
		name: "WithField and WithError",
		run: func(logger *logrus.Logger, t time.Time) {
			logger.
				WithTime(t).
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Info("my log entry")
		},
		out: map[string]interface{}{
			"level": "INFO",
			"msg":   "my log entry",
			"extra": map[string]interface{}{
				"foo":   "bar",
				"error": "test error",
			},
			"env":    "prod",
			"group":  "test",
			"system": "system",
		},
	},
	{
		name: "WithField and Error",
		run: func(logger *logrus.Logger, t time.Time) {
			logger.WithTime(t).WithField("foo", "bar").Error("my log entry")
		},
		out: map[string]interface{}{
			"level": "ERROR",
			"msg":   "my log entry",
			"extra": map[string]interface{}{
				"foo": "bar",
			},
			"env":    "prod",
			"group":  "test",
			"system": "system",
		},
	},
	{
		name: "WithField, WithError and Error",
		run: func(logger *logrus.Logger, t time.Time) {
			logger.
				WithTime(t).
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"level": "ERROR",
			"msg":   "my log entry",
			"extra": map[string]interface{}{
				"foo":   "bar",
				"error": "test error",
			},
			"env":    "prod",
			"group":  "test",
			"system": "system",
		},
	},
}
