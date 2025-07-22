package msg

import (
	"fmt"
	"strings"
)

type PackagedQueue struct {
	msgs []PackagedMessage
}

func NewQueue() *PackagedQueue {
	pkgQueue := make([]PackagedMessage, 0)
	return &PackagedQueue{msgs: pkgQueue}
}

func (q *PackagedQueue) Size() int {
	return len(q.msgs)
}

func (q *PackagedQueue) IsEmpty() bool {
	return len(q.msgs) == 0
}

func (q *PackagedQueue) Enqueue(newMsg PackagedMessage) {
	// naive implementation, reallocs a slice every call
	q.msgs = append(q.msgs, newMsg)
}

func (q *PackagedQueue) Dequeue() (PackagedMessage, bool) {
	// in order to increase efficiency look into:
	// - ring buffers or linked lists
	if q.IsEmpty() {
		return PackagedMessage{}, false
	}
	nextMsg := q.msgs[0]

	// reslicing the queue, not efficient but functional
	q.msgs = q.msgs[1:]

	return nextMsg, true
}

func (q *PackagedQueue) QueueSummary() string {
	if q.IsEmpty() {
		return "Queue is empty.\n"
	}

	var subjects []string
	for _, msg := range q.msgs {
		subjects = append(subjects, msg.Subject)
	}

	subjectSummary := strings.Join(subjects, " || ")
	return fmt.Sprintf("%d messages in queue\nMessage subjects: %s\n", len(q.msgs), subjectSummary)
}
