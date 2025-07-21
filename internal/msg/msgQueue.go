package msg

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

// naive implementation, reallocs a slice every call
func (q *Queue) Enqueue(newMsg Message) {
	q.msgs = append(q.msgs, newMsg)
}

// in order to increase efficiency look into:
// - ring buffers or linked lists
func (q *Queue) Dequeue() (Message, bool) {
	if q.IsEmpty() {
		return Message{}, false
	}
	nextMsg := q.msgs[0]

	// reslicing the queue, not efficient but functional
	q.msgs = q.msgs[1:]

	return nextMsg, true
}
