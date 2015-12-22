/* Benchmarking against several packages

We compare against:
* Built in append
* container/list
* github.com/eapache/queue

*/
package deque

import (
	"container/list"
	"sync"
	"testing"
	"time"

	"github.com/eapache/queue"
	"github.com/tebeka/deque"
)

type HistRec struct {
	id        int64
	hash      string
	version   string
	timestamp time.Time // when was committed
}

type History interface {
	Add(rec *HistRec)
	Len() int
}

func newHistRec(adId int) *HistRec {
	return &HistRec{
		id:        int64(adId),
		hash:      "hash",
		version:   "some version",
		timestamp: time.Now(),
	}
}

func runBenchmark(b *testing.B, hist History) {
	for i := 0; i < b.N; i++ {
		hist.Add(newHistRec(i))
	}
	if hist.Len() != b.N {
		b.Fatalf("wrong length - %d (should be %d)", hist.Len(), b.N)
	}
}

type AHist struct {
	sync.Mutex
	hist []*HistRec
}

func (ah *AHist) Add(rec *HistRec) {
	ah.Lock()
	defer ah.Unlock()

	ah.hist = append(ah.hist, rec)
}

func (ah *AHist) Len() int {
	return len(ah.hist)
}

func BenchmarkHistAppend(b *testing.B) {
	ah := &AHist{}
	runBenchmark(b, ah)
}

type LHist struct {
	sync.Mutex
	hist *list.List
}

func NewLHist() *LHist {
	return &LHist{
		hist: list.New(),
	}
}

func (lh *LHist) Add(rec *HistRec) {
	lh.Lock()
	defer lh.Unlock()
	lh.hist.PushBack(rec)
}

func (lh *LHist) Len() int {
	return lh.hist.Len()
}

func BenchmarkHistList(b *testing.B) {
	lh := NewLHist()
	runBenchmark(b, lh)
}

type QHist struct {
	sync.Mutex
	hist *queue.Queue
}

func NewQHist() *QHist {
	return &QHist{
		hist: queue.New(),
	}
}

func (qh *QHist) Add(rec *HistRec) {
	qh.Lock()
	defer qh.Unlock()
	qh.hist.Add(rec)
}

func (qh *QHist) Len() int {
	return qh.hist.Length()
}

func BenchmarkHistQueue(b *testing.B) {
	qh := NewQHist()
	runBenchmark(b, qh)
}

type DQHist struct {
	sync.Mutex
	hist *deque.Deque
}

func NewDQHist() *DQHist {
	return &DQHist{
		hist: deque.New(),
	}
}

func (dqh *DQHist) Add(rec *HistRec) {
	dqh.Lock()
	defer dqh.Unlock()
	dqh.hist.Append(rec)
}

func (dqh *DQHist) Len() int {
	return dqh.hist.Len()
}

func BenchmarkHistDeque(b *testing.B) {
	dqh := NewDQHist()
	runBenchmark(b, dqh)
}
