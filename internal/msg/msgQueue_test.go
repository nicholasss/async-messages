package msg

import "testing"

var queueSecretKey = []byte("GgfY0UssupyYBlFy92/ENsq5/Qy8dq3bh3Mp8hZcPMDEdSnxMgi5E1TPzJuHVHzRs60aq6r7gKyLGwbauaUn1Q==")

func TestEmptyQueue(t *testing.T) {
	// should be empty on creation
	queue := NewQueue()
	if queue.IsEmpty() != true {
		t.Errorf("Queue should be empty on creation")
	}
	if queue.Size() != 0 {
		t.Errorf("Queue should be emtpy on creation")
	}
}

func TestQueueSummary(t *testing.T) {
	tt := []struct {
		rawMsg RawMessage
	}{
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Tuesday",
				Body:       "I am planning on proceeding on tuesday since there is a break in the weather",
			},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Kevin",
				ToVessel:   "Liberty",
				FromName:   "Bob",
				FromVessel: "Snow",
				Subject:    "Re: Tuesday",
				Body:       "I will need to wait longer because of needed repair work. Hope to catch up.",
			},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
		},
	}

	wantQueueSummary := "3 messages in queue\nMessage subjects: tuesday || re: tuesday || re: re: tuesday\n"

	// packaging messages
	msgs := make([]PackagedMessage, 0)
	for _, tc := range tt {
		msg, err := tc.rawMsg.ToPackagedMessage(queueSecretKey)
		if err != nil {
			t.Errorf("Unable to make new message due to: %q", err)
		}

		msgs = append(msgs, *msg)
	}

	// adding to queue
	queue := NewQueue()
	for _, packagedMsg := range msgs {
		queue.Enqueue(packagedMsg)
	}

	// performing test
	gotQueueSummary := queue.QueueSummary()
	if wantQueueSummary != gotQueueSummary {
		t.Errorf("queue dump does not match. got=%q want=%q", gotQueueSummary, wantQueueSummary)
	}
}

func TestSizeAndEnqueueDequeueFunction(t *testing.T) {
	tt := []struct {
		rawMsg RawMessage
	}{
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Tuesday",
				Body:       "I am planning on proceeding on tuesday since there is a break in the weather",
			},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Kevin",
				ToVessel:   "Liberty",
				FromName:   "Bob",
				FromVessel: "Snow",
				Subject:    "Re: Tuesday",
				Body:       "I will need to wait longer because of needed repair work. Hope to catch up.",
			},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
		},
	}

	// packaging messages
	msgs := make([]PackagedMessage, 0)
	for _, tc := range tt {
		msg, err := tc.rawMsg.ToPackagedMessage(queueSecretKey)
		if err != nil {
			t.Errorf("Unable to make new message due to: %q", err)
		}

		msgs = append(msgs, *msg)
	}

	// create queue and enqueue first
	queue := NewQueue()
	queue.Enqueue(msgs[0])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 1 {
		t.Error("Queue should contain one item")
	}

	// enqueue second
	queue.Enqueue(msgs[1])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 2 {
		t.Error("Queue should contain two items")
	}

	// enqueue third
	queue.Enqueue(msgs[2])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 3 {
		t.Error("Queue should contain three items")
	}

	// dequeue first
	msg1, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg1 != msgs[0] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg1.String(), msgs[0].String())
	}
	if queue.Size() != 2 {
		t.Error("Queue should contain two items")
	}

	// dequeue second
	msg2, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg2 != msgs[1] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg2.String(), msgs[2].String())
	}
	if queue.Size() != 1 {
		t.Error("Queue should contain one item")
	}

	// dequeue third
	msg3, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg3 != msgs[2] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg3.String(), msgs[3].String())
	}

	// attempting to dequeue empty queue
	_, ok = queue.Dequeue()
	if ok {
		t.Error("Dequeue-ing an empty queue should not be ok")
	}
	if queue.IsEmpty() != true {
		t.Error("Unable to fully dequeue queue")
	}
	if queue.Size() != 0 {
		t.Error("Unable to return size of emtpy queue as 0")
	}
}

// old benchmark code
//
// func BenchmarkQueueEnqueue(b *testing.B) {
// 	msgExample, err := NewMessage("kevin", "bob", "important", "this is an example body", queueSecretKey)
// 	if err != nil {
// 		b.Errorf("Unable to create msg example due to: %q", err)
// 	}
//
// 	for b.Loop() {
// 		queue := NewQueue()
// 		for range 1000 {
// 			queue.Enqueue(*msgExample)
// 		}
// 	}
// }
//
// func BenchmarkQueueDequeue(b *testing.B) {
// 	msgExample, err := NewMessage("kevin", "bob", "important", "this is an example body", queueSecretKey)
// 	if err != nil {
// 		b.Errorf("Unable to create msg example due to: %q", err)
// 	}
//
// 	for b.Loop() {
// 		// not measuring queue setup
// 		b.StopTimer()
// 		queue := NewQueue()
// 		for range 1000 {
// 			queue.Enqueue(*msgExample)
// 		}
//
// 		// measuring dequeue explicitly
// 		b.StartTimer()
// 		for range 1000 {
// 			_, ok := queue.Dequeue()
// 			if !ok {
// 				b.Error("Unable to dequeue")
// 			}
// 		}
// 	}
// }
