package msg

import (
	"fmt"
	"strings"
)

type Queue struct {
	msgs []Message
}

func NewQueue() *Queue {
	rawQueue := make([]Message, 0)
	return &Queue{msgs: rawQueue}
}

func (q *Queue) Size() int {
	return len(q.msgs)
}

func (q *Queue) IsEmpty() bool {
	return len(q.msgs) == 0
}

func (q *Queue) Enqueue(newMsg Message) {
	// naive implementation, reallocs a slice every call
	q.msgs = append(q.msgs, newMsg)
}

func (q *Queue) Dequeue() (Message, bool) {
	// in order to increase efficiency look into:
	// - ring buffers or linked lists
	if q.IsEmpty() {
		return Message{}, false
	}
	nextMsg := q.msgs[0]

	// reslicing the queue, not efficient but functional
	q.msgs = q.msgs[1:]

	return nextMsg, true
}

func (q *Queue) DumpToString() string {
	var subjects []string
	for _, msg := range q.msgs {
		subjects = append(subjects, msg.Subject)
	}

	subjectSummary := strings.Join(subjects, " || ")
	return fmt.Sprintf("Number of messages in queue: %d\nSubjects: %s\n", len(q.msgs), subjectSummary)
}
