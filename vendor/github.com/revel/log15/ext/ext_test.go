package ext

import (
	"errors"
	log "github.com/inconshreveable/log15"
	"math"
	"testing"
)

func testHandler() (log.Handler, *log.Record) {
	rec := new(log.Record)
	return log.FuncHandler(func(r *log.Record) error {
		*rec = *r
		return nil
	}), rec
}

func TestHotSwapHandler(t *testing.T) {
	t.Parallel()

	h1, r1 := testHandler()

	l := log.New()
	h := HotSwapHandler(h1)
	l.SetHandler(h)

	l.Info("to h1")
	if r1.Msg != "to h1" {
		t.Fatalf("didn't get expected message to h1")
	}

	h2, r2 := testHandler()
	h.Swap(h2)
	l.Info("to h2")
	if r2.Msg != "to h2" {
		t.Fatalf("didn't get expected message to h2")
	}
}

func TestSpeculativeHandler(t *testing.T) {
	t.Parallel()

	// test with an even multiple of the buffer size, less than full buffer size
	// and not a multiple of the buffer size
	for _, count := range []int{10000, 50, 432} {
		recs := make(chan *log.Record)
		done := make(chan int)
		spec := SpeculativeHandler(100, log.ChannelHandler(recs))

		go func() {
			defer close(done)
			expectedCount := int(math.Min(float64(count), float64(100)))
			expectedIdx := count - expectedCount
			for r := range recs {
				if r.Ctx[1] != expectedIdx {
					t.Errorf("Bad ctx 'i', got %d expected %d", r.Ctx[1], expectedIdx)
					return
				}
				expectedIdx++
				expectedCount--

				if expectedCount == 0 {
					// got everything we expected
					break
				}
			}

			select {
			case <-recs:
				t.Errorf("got an extra record we shouldn't have!")
			default:
			}
		}()

		lg := log.New()
		lg.SetHandler(spec)
		for i := 0; i < count; i++ {
			lg.Debug("test speculative", "i", i)
		}

		go spec.Flush()

		// wait for the go routine to finish
		<-done
	}
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	h, r := testHandler()
	lg := log.New()
	lg.SetHandler(EscalateErrHandler(
		log.LvlFilterHandler(log.LvlError, h)))

	lg.Debug("some function result", "err", nil)
	if r.Msg != "" {
		t.Fatalf("Expected debug level message to be filtered")
	}

	lg.Debug("some function result", "err", errors.New("failed operation"))
	if r.Msg != "some function result" {
		t.Fatalf("Expected debug level message to be escalated and pass lvlfilter")
	}

	if r.Lvl != log.LvlError {
		t.Fatalf("Expected debug level message to be escalated to LvlError")
	}
}
