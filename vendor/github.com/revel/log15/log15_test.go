package log15

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"runtime"
	"sync"
	"testing"
	"time"
)

func testHandler() (Handler, *Record) {
	rec := new(Record)
	return FuncHandler(func(r *Record) error {
		*rec = *r
		return nil
	}), rec
}

func testLogger() (Logger, Handler, *Record) {
	l := New()
	h, r := testHandler()
	l.SetHandler(LazyHandler(h))
	return l, h, r
}

func TestLazy(t *testing.T) {
	t.Parallel()

	x := 1
	lazy := func() int {
		return x
	}

	l, _, r := testLogger()
	l.Info("", "x", Lazy{lazy})
	if r.Ctx[1] != 1 {
		t.Fatalf("Lazy function not evaluated, got %v, expected %d", r.Ctx[1], 1)
	}

	x = 2
	l.Info("", "x", Lazy{lazy})
	if r.Ctx[1] != 2 {
		t.Fatalf("Lazy function not evaluated, got %v, expected %d", r.Ctx[1], 1)
	}
}

func TestInvalidLazy(t *testing.T) {
	t.Parallel()

	l, _, r := testLogger()
	validate := func() {
		if len(r.Ctx) < 4 {
			t.Fatalf("Invalid lazy, got %d args, expecting at least 4", len(r.Ctx))
		}

		if r.Ctx[2] != errorKey {
			t.Fatalf("Invalid lazy, got key %s expecting %s", r.Ctx[2], errorKey)
		}
	}

	l.Info("", "x", Lazy{1})
	validate()

	l.Info("", "x", Lazy{func(x int) int { return x }})
	validate()

	l.Info("", "x", Lazy{func() {}})
	validate()
}

func TestCtx(t *testing.T) {
	t.Parallel()

	l, _, r := testLogger()
	l.Info("", Ctx{"x": 1, "y": "foo", "tester": t})
	if len(r.Ctx) != 6 {
		t.Fatalf("Expecting Ctx tansformed into %d ctx args, got %d: %v", 6, len(r.Ctx), r.Ctx)
	}
}

func testFormatter(f Format) (Logger, *bytes.Buffer) {
	l := New()
	var buf bytes.Buffer
	l.SetHandler(StreamHandler(&buf, f))
	return l, &buf
}

func TestJson(t *testing.T) {
	t.Parallel()

	l, buf := testFormatter(JsonFormat())
	l.Error("some message", "x", 1, "y", 3.2)

	var v map[string]interface{}
	decoder := json.NewDecoder(buf)
	if err := decoder.Decode(&v); err != nil {
		t.Fatalf("Error decoding JSON: %v", v)
	}

	validate := func(key string, expected interface{}) {
		if v[key] != expected {
			t.Fatalf("Got %v expected %v for %v", v[key], expected, key)
		}
	}

	validate("msg", "some message")
	validate("x", float64(1)) // all numbers are floats in JSON land
	validate("y", 3.2)
	validate("lvl", "eror")
}

type testtype struct {
	name string
}

func (tt testtype) String() string {
	return tt.name
}

func TestLogfmt(t *testing.T) {
	t.Parallel()

	var nilVal *testtype

	l, buf := testFormatter(LogfmtFormat())
	l.Error("some message", "x", 1, "y", 3.2, "equals", "=", "quote", "\"",
		"nil", nilVal, "carriage_return", "bang"+string('\r')+"foo", "tab", "bar	baz", "newline", "foo\nbar")

	// skip timestamp in comparison
	got := buf.Bytes()[27:buf.Len()]
	expected := []byte(`lvl=eror msg="some message" x=1 y=3.200 equals="=" quote="\"" nil=nil carriage_return="bang\rfoo" tab="bar\tbaz" newline="foo\nbar"` + "\n")
	if !bytes.Equal(got, expected) {
		t.Fatalf("Got %s, expected %s", got, expected)
	}
}

func TestMultiHandler(t *testing.T) {
	t.Parallel()

	h1, r1 := testHandler()
	h2, r2 := testHandler()
	l := New()
	l.SetHandler(MultiHandler(h1, h2))
	l.Debug("clone")

	if r1.Msg != "clone" {
		t.Fatalf("wrong value for h1.Msg. Got %s expected %s", r1.Msg, "clone")
	}

	if r2.Msg != "clone" {
		t.Fatalf("wrong value for h2.Msg. Got %s expected %s", r2.Msg, "clone")
	}

}

type waitHandler struct {
	ch chan Record
}

func (h *waitHandler) Log(r *Record) error {
	h.ch <- *r
	return nil
}

func TestBufferedHandler(t *testing.T) {
	t.Parallel()

	ch := make(chan Record)
	l := New()
	l.SetHandler(BufferedHandler(0, &waitHandler{ch}))

	l.Debug("buffer")
	if r := <-ch; r.Msg != "buffer" {
		t.Fatalf("wrong value for r.Msg. Got %s expected %s", r.Msg, "")
	}
}

func TestLogContext(t *testing.T) {
	t.Parallel()

	l, _, r := testLogger()
	l = l.New("foo", "bar")
	l.Crit("baz")

	if len(r.Ctx) != 2 {
		t.Fatalf("Expected logger context in record context. Got length %d, expected %d", len(r.Ctx), 2)
	}

	if r.Ctx[0] != "foo" {
		t.Fatalf("Wrong context key, got %s expected %s", r.Ctx[0], "foo")
	}

	if r.Ctx[1] != "bar" {
		t.Fatalf("Wrong context value, got %s expected %s", r.Ctx[1], "bar")
	}
}

func TestMapCtx(t *testing.T) {
	t.Parallel()

	l, _, r := testLogger()
	l.Crit("test", Ctx{"foo": "bar"})

	if len(r.Ctx) != 2 {
		t.Fatalf("Wrong context length, got %d, expected %d", len(r.Ctx), 2)
	}

	if r.Ctx[0] != "foo" {
		t.Fatalf("Wrong context key, got %s expected %s", r.Ctx[0], "foo")
	}

	if r.Ctx[1] != "bar" {
		t.Fatalf("Wrong context value, got %s expected %s", r.Ctx[1], "bar")
	}
}

func TestLvlFilterHandler(t *testing.T) {
	t.Parallel()

	l := New()
	h, r := testHandler()
	l.SetHandler(LvlFilterHandler(LvlWarn, h))
	l.Info("info'd")

	if r.Msg != "" {
		t.Fatalf("Expected zero record, but got record with msg: %v", r.Msg)
	}

	l.Warn("warned")
	if r.Msg != "warned" {
		t.Fatalf("Got record msg %s expected %s", r.Msg, "warned")
	}

	l.Warn("error'd")
	if r.Msg != "error'd" {
		t.Fatalf("Got record msg %s expected %s", r.Msg, "error'd")
	}
}

func TestNetHandler(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", l)
	}

	errs := make(chan error)
	go func() {
		c, err := l.Accept()
		if err != nil {
			errs <- fmt.Errorf("Failed to accept connection: %v", err)
			return
		}

		rd := bufio.NewReader(c)
		s, err := rd.ReadString('\n')
		if err != nil {
			errs <- fmt.Errorf("Failed to read string: %v", err)
			return
		}

		got := s[27:]
		expected := "lvl=info msg=test x=1\n"
		if got != expected {
			t.Errorf("Got log line %s, expected %s", got, expected)
		}

		errs <- nil
	}()

	lg := New()
	h, err := NetHandler("tcp", l.Addr().String(), LogfmtFormat())
	if err != nil {
		t.Fatal(err)
	}
	lg.SetHandler(h)
	lg.Info("test", "x", 1)

	select {
	case <-time.After(time.Second):
		t.Fatalf("Test timed out!")
	case err := <-errs:
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestMatchFilterHandler(t *testing.T) {
	t.Parallel()

	l, h, r := testLogger()
	l.SetHandler(MatchFilterHandler("err", nil, h))

	l.Crit("test", "foo", "bar")
	if r.Msg != "" {
		t.Fatalf("expected filter handler to discard msg")
	}

	l.Crit("test2", "err", "bad fd")
	if r.Msg != "" {
		t.Fatalf("expected filter handler to discard msg")
	}

	l.Crit("test3", "err", nil)
	if r.Msg != "test3" {
		t.Fatalf("expected filter handler to allow msg")
	}
}

func TestMatchFilterBuiltin(t *testing.T) {
	t.Parallel()

	l, h, r := testLogger()
	l.SetHandler(MatchFilterHandler("lvl", LvlError, h))
	l.Info("does not pass")

	if r.Msg != "" {
		t.Fatalf("got info level record that should not have matched")
	}

	l.Error("error!")
	if r.Msg != "error!" {
		t.Fatalf("did not get error level record that should have matched")
	}

	r.Msg = ""
	l.SetHandler(MatchFilterHandler("msg", "matching message", h))
	l.Info("doesn't match")
	if r.Msg != "" {
		t.Fatalf("got record with wrong message matched")
	}

	l.Debug("matching message")
	if r.Msg != "matching message" {
		t.Fatalf("did not get record which matches")
	}
}

type failingWriter struct {
	fail bool
}

func (w *failingWriter) Write(buf []byte) (int, error) {
	if w.fail {
		return 0, errors.New("fail")
	} else {
		return len(buf), nil
	}
}

func TestFailoverHandler(t *testing.T) {
	t.Parallel()

	l := New()
	h, r := testHandler()
	w := &failingWriter{false}

	l.SetHandler(FailoverHandler(
		StreamHandler(w, JsonFormat()),
		h))

	l.Debug("test ok")
	if r.Msg != "" {
		t.Fatalf("expected no failover")
	}

	w.fail = true
	l.Debug("test failover", "x", 1)
	if r.Msg != "test failover" {
		t.Fatalf("expected failover")
	}

	if len(r.Ctx) != 4 {
		t.Fatalf("expected additional failover ctx")
	}

	got := r.Ctx[2]
	expected := "failover_err_0"
	if got != expected {
		t.Fatalf("expected failover ctx. got: %s, expected %s", got, expected)
	}
}

// https://github.com/inconshreveable/log15/issues/16
func TestIndependentSetHandler(t *testing.T) {
	t.Parallel()

	parent, _, r := testLogger()
	child := parent.New()
	child.SetHandler(DiscardHandler())
	parent.Info("test")
	if r.Msg != "test" {
		t.Fatalf("parent handler affected by child")
	}
}

// https://github.com/inconshreveable/log15/issues/16
func TestInheritHandler(t *testing.T) {
	t.Parallel()

	parent, _, r := testLogger()
	child := parent.New()
	parent.SetHandler(DiscardHandler())
	child.Info("test")
	if r.Msg == "test" {
		t.Fatalf("child handler affected not affected by parent")
	}
}

func TestCallerFileHandler(t *testing.T) {
	t.Parallel()

	l := New()
	h, r := testHandler()
	l.SetHandler(CallerFileHandler(h))

	l.Info("baz")
	_, _, line, _ := runtime.Caller(0)

	if len(r.Ctx) != 2 {
		t.Fatalf("Expected caller in record context. Got length %d, expected %d", len(r.Ctx), 2)
	}

	const key = "caller"

	if r.Ctx[0] != key {
		t.Fatalf("Wrong context key, got %s expected %s", r.Ctx[0], key)
	}

	s, ok := r.Ctx[1].(string)
	if !ok {
		t.Fatalf("Wrong context value type, got %T expected string", r.Ctx[1])
	}

	exp := fmt.Sprint("log15_test.go:", line-1)
	if s != exp {
		t.Fatalf("Wrong context value, got %s expected string matching %s", s, exp)
	}
}

func TestCallerFuncHandler(t *testing.T) {
	t.Parallel()

	l := New()
	h, r := testHandler()
	l.SetHandler(CallerFuncHandler(h))

	l.Info("baz")

	if len(r.Ctx) != 2 {
		t.Fatalf("Expected caller in record context. Got length %d, expected %d", len(r.Ctx), 2)
	}

	const key = "fn"

	if r.Ctx[0] != key {
		t.Fatalf("Wrong context key, got %s expected %s", r.Ctx[0], key)
	}

	const regex = ".*\\.TestCallerFuncHandler"

	s, ok := r.Ctx[1].(string)
	if !ok {
		t.Fatalf("Wrong context value type, got %T expected string", r.Ctx[1])
	}

	match, err := regexp.MatchString(regex, s)
	if err != nil {
		t.Fatalf("Error matching %s to regex %s: %v", s, regex, err)
	}

	if !match {
		t.Fatalf("Wrong context value, got %s expected string matching %s", s, regex)
	}
}

// https://github.com/inconshreveable/log15/issues/27
func TestCallerStackHandler(t *testing.T) {
	t.Parallel()

	l := New()
	h, r := testHandler()
	l.SetHandler(CallerStackHandler("%#v", h))

	lines := []int{}

	func() {
		l.Info("baz")
		_, _, line, _ := runtime.Caller(0)
		lines = append(lines, line-1)
	}()
	_, file, line, _ := runtime.Caller(0)
	lines = append(lines, line-1)

	if len(r.Ctx) != 2 {
		t.Fatalf("Expected stack in record context. Got length %d, expected %d", len(r.Ctx), 2)
	}

	const key = "stack"

	if r.Ctx[0] != key {
		t.Fatalf("Wrong context key, got %s expected %s", r.Ctx[0], key)
	}

	s, ok := r.Ctx[1].(string)
	if !ok {
		t.Fatalf("Wrong context value type, got %T expected string", r.Ctx[1])
	}

	exp := "["
	for i, line := range lines {
		if i > 0 {
			exp += " "
		}
		exp += fmt.Sprint(file, ":", line)
	}
	exp += "]"

	if s != exp {
		t.Fatalf("Wrong context value, got %s expected string matching %s", s, exp)
	}
}

// tests that when logging concurrently to the same logger
// from multiple goroutines that the calls are handled independently
// this test tries to trigger a previous bug where concurrent calls could
// corrupt each other's context values
//
// this test runs N concurrent goroutines each logging a fixed number of
// records and a handler that buckets them based on the index passed in the context.
// if the logger is not concurrent-safe then the values in the buckets will not all be the same
//
// https://github.com/inconshreveable/log15/pull/30
func TestConcurrent(t *testing.T) {
	root := New()
	// this was the first value that triggered
	// go to allocate extra capacity in the logger's context slice which
	// was necessary to trigger the bug
	const ctxLen = 34
	l := root.New(make([]interface{}, ctxLen)...)
	const goroutines = 8
	var res [goroutines]int
	l.SetHandler(SyncHandler(FuncHandler(func(r *Record) error {
		res[r.Ctx[ctxLen+1].(int)]++
		return nil
	})))
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				l.Info("test message", "goroutine_idx", idx)
			}
		}(i)
	}
	wg.Wait()
	for _, val := range res[:] {
		if val != 10000 {
			t.Fatalf("Wrong number of messages for context: %+v", res)
		}
	}
}
