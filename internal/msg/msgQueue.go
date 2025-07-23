package msg

import (
	"fmt"
	"strings"
	"sync"
)

type PackagedQueue struct {
	msgs []PackagedMessage
	mux  *sync.Mutex
}

func NewQueue() *PackagedQueue {
	pkgQueue := make([]PackagedMessage, 0)
	mux := &sync.Mutex{}

	return &PackagedQueue{msgs: pkgQueue, mux: mux}
}

func (q *PackagedQueue) Size() int {
	q.mux.Lock()
	qLen := len(q.msgs)
	q.mux.Unlock()

	return qLen
}

func (q *PackagedQueue) IsEmpty() bool {
	q.mux.Lock()
	isEmpty := len(q.msgs) == 0
	q.mux.Unlock()

	return isEmpty
}

func (q *PackagedQueue) Enqueue(newMsg PackagedMessage) {
	// naive implementation, reallocs a slice every call
	q.mux.Lock()
	q.msgs = append(q.msgs, newMsg)
	q.mux.Unlock()
}

func (q *PackagedQueue) Dequeue() (PackagedMessage, bool) {
	// in order to increase efficiency look into:
	// - ring buffers or linked lists
	if q.IsEmpty() {
		return PackagedMessage{}, false
	}

	q.mux.Lock()
	nextMsg := q.msgs[0]

	// reslicing the queue, not efficient but functional
	q.msgs = q.msgs[1:]
	q.mux.Unlock()

	return nextMsg, true
}

func (q *PackagedQueue) QueueSummary() string {
	if q.IsEmpty() {
		return "Queue is empty.\n"
	}

	q.mux.Lock()
	var subjects []string
	for _, msg := range q.msgs {
		subjects = append(subjects, msg.Subject)
	}
	q.mux.Unlock()

	subjectSummary := strings.Join(subjects, " || ")
	return fmt.Sprintf("%d messages in queue\nMessage subjects: %s\n", len(q.msgs), subjectSummary)
}
