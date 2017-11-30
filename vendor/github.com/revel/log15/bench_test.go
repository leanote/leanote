// +build go1.3

package log15

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
	"time"
)

func BenchmarkStreamNoCtx(b *testing.B) {
	lg := New()

	buf := bytes.Buffer{}
	lg.SetHandler(StreamHandler(&buf, LogfmtFormat()))

	for i := 0; i < b.N; i++ {
		lg.Info("test message")
		buf.Reset()
	}
}

func BenchmarkDiscard(b *testing.B) {
	lg := New()
	lg.SetHandler(DiscardHandler())

	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkCallerFileHandler(b *testing.B) {
	lg := New()
	lg.SetHandler(CallerFileHandler(DiscardHandler()))

	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkCallerFuncHandler(b *testing.B) {
	lg := New()
	lg.SetHandler(CallerFuncHandler(DiscardHandler()))

	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkLogfmtNoCtx(b *testing.B) {
	r := Record{
		Time: time.Now(),
		Lvl:  LvlInfo,
		Msg:  "test message",
		Ctx:  []interface{}{},
	}

	logfmt := LogfmtFormat()
	for i := 0; i < b.N; i++ {
		logfmt.Format(&r)
	}
}

func BenchmarkJsonNoCtx(b *testing.B) {
	r := Record{
		Time: time.Now(),
		Lvl:  LvlInfo,
		Msg:  "test message",
		Ctx:  []interface{}{},
	}

	jsonfmt := JsonFormat()
	for i := 0; i < b.N; i++ {
		jsonfmt.Format(&r)
	}
}

func BenchmarkMultiLevelFilter(b *testing.B) {
	handler := MultiHandler(
		LvlFilterHandler(LvlDebug, DiscardHandler()),
		LvlFilterHandler(LvlError, DiscardHandler()),
	)

	lg := New()
	lg.SetHandler(handler)
	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkDescendant1(b *testing.B) {
	lg := New()
	lg.SetHandler(DiscardHandler())
	lg = lg.New()
	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkDescendant2(b *testing.B) {
	lg := New()
	lg.SetHandler(DiscardHandler())
	for i := 0; i < 2; i++ {
		lg = lg.New()
	}
	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkDescendant4(b *testing.B) {
	lg := New()
	lg.SetHandler(DiscardHandler())
	for i := 0; i < 4; i++ {
		lg = lg.New()
	}
	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

func BenchmarkDescendant8(b *testing.B) {
	lg := New()
	lg.SetHandler(DiscardHandler())
	for i := 0; i < 8; i++ {
		lg = lg.New()
	}
	for i := 0; i < b.N; i++ {
		lg.Info("test message")
	}
}

// Copied from https://github.com/uber-go/zap/blob/master/benchmarks/log15_bench_test.go
// (MIT License)
func newLog15() Logger {
	logger := New()
	logger.SetHandler(StreamHandler(ioutil.Discard, JsonFormat()))
	return logger
}

var errExample = errors.New("fail")

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var _jane = user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}

func BenchmarkLog15AddingFields(b *testing.B) {
	logger := newLog15()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.",
				"int", 1,
				"int64", int64(1),
				"float", 3.0,
				"string", "four!",
				"bool", true,
				"time", time.Unix(0, 0),
				"error", errExample.Error(),
				"duration", time.Second,
				"user-defined type", _jane,
				"another string", "done!",
			)
		}
	})
}

func BenchmarkLog15WithAccumulatedContext(b *testing.B) {
	logger := newLog15().New(
		"int", 1,
		"int64", int64(1),
		"float", 3.0,
		"string", "four!",
		"bool", true,
		"time", time.Unix(0, 0),
		"error", errExample.Error(),
		"duration", time.Second,
		"user-defined type", _jane,
		"another string", "done!",
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go really fast.")
		}
	})
}

func BenchmarkLog15WithoutFields(b *testing.B) {
	logger := newLog15()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.")
		}
	})
}
