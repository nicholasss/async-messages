package msg

type queue struct {
	msgs []Message
}

func NewQueue() *queue {
	rawQueue := make([]Message, 0)
	return &queue{msgs: rawQueue}
}

func (q *queue) Size() int {
	return len(q.msgs)
}

func (q *queue) IsEmpty() bool {
	return len(q.msgs) == 0
}

// naive implementation, reallocs a slice every call
func (q *queue) Enqueue(newMsg Message) {
	q.msgs = append(q.msgs, newMsg)
}

// in order to increase efficiency look into:
// - ring buffers or linked lists
func (q *queue) Dequeue() (Message, bool) {
	if q.IsEmpty() {
		return Message{}, false
	}
	nextMsg := q.msgs[0]

	// reslicing the queue, not efficient but functional
	q.msgs = q.msgs[1:]

	return nextMsg, true
}
