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

func TestDumpToString(t *testing.T) {
	rawMsgs := []struct {
		to      string
		from    string
		subject string
		body    string
	}{
		{
			to:      "Bob",
			from:    "Kevin",
			subject: "Re: Tuesday",
			body:    "This is super important that we work on this tuesday.",
		},
		{
			to:      "Kevin",
			from:    "Bob",
			subject: "Re: Re: Tuesday",
			body:    "I am not sure that I can do tuesday though.",
		},
		{
			to:      "Bob",
			from:    "Kevin",
			subject: "Re: Re: Re: Tuesday",
			body:    "You can do tue, just let me know when on tue. I am free after lunch.",
		},
	}

	wantQueueDump := "Number of messages in queue: 3\nSubjects: Re: Tuesday || Re: Re: Tuesday || Re: Re: Re: Tuesday\n"

	packagedMsgs := make([]Message, 0)
	for _, rawMsg := range rawMsgs {
		msg, err := NewMessage(rawMsg.to, rawMsg.from, rawMsg.subject, rawMsg.body, queueSecretKey)
		if err != nil {
			t.Errorf("Unable to make new message due to: %q", err)
		}

		packagedMsgs = append(packagedMsgs, *msg)
	}

	queue := NewQueue()
	for _, packagedMsg := range packagedMsgs {
		queue.Enqueue(packagedMsg)
	}

	gotQueueDump := queue.DumpToString()

	if wantQueueDump != gotQueueDump {
		t.Errorf("queue dump does not match. got=%q want=%q", gotQueueDump, wantQueueDump)
	}
}

func TestSizeAndEnqueueDequeueFunction(t *testing.T) {
	rawMsgs := []struct {
		to      string
		from    string
		subject string
		body    string
	}{
		{
			to:      "Bob",
			from:    "Kevin",
			subject: "Re: Tuesday",
			body:    "This is super important that we work on this tuesday.",
		},
		{
			to:      "Kevin",
			from:    "Bob",
			subject: "Re: Re: Tuesday",
			body:    "I am not sure that I can do tuesday though.",
		},
		{
			to:      "Bob",
			from:    "Kevin",
			subject: "Re: Re: Re: Tuesday",
			body:    "You can do tue, just let me know when on tue. I am free after lunch.",
		},
	}
	packagedMsgs := make([]Message, 0)
	for _, rawMsg := range rawMsgs {
		msg, err := NewMessage(rawMsg.to, rawMsg.from, rawMsg.subject, rawMsg.body, queueSecretKey)
		if err != nil {
			t.Errorf("Unable to make new message due to: %q", err)
		}

		packagedMsgs = append(packagedMsgs, *msg)
	}

	queue := NewQueue()

	queue.Enqueue(packagedMsgs[0])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 1 {
		t.Error("Queue should contain one item")
	}

	queue.Enqueue(packagedMsgs[1])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 2 {
		t.Error("Queue should contain two items")
	}

	queue.Enqueue(packagedMsgs[2])
	if queue.IsEmpty() != false {
		t.Error("Queue should not be empty after enqueue")
	}
	if queue.Size() != 3 {
		t.Error("Queue should contain three items")
	}

	msg1, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg1 != packagedMsgs[0] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg1.ToString(), packagedMsgs[0].ToString())
	}
	if queue.Size() != 2 {
		t.Error("Queue should contain two items")
	}

	msg2, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg2 != packagedMsgs[1] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg2.ToString(), packagedMsgs[2].ToString())
	}
	if queue.Size() != 1 {
		t.Error("Queue should contain one item")
	}

	msg3, ok := queue.Dequeue()
	if !ok {
		t.Error("Unable to dequeue message")
	}
	if msg3 != packagedMsgs[2] {
		t.Errorf("Unable to dequeue correct message. got=%q want=%q", msg3.ToString(), packagedMsgs[3].ToString())
	}

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

func BenchmarkQueueEnqueue(b *testing.B) {
	msgExample, err := NewMessage("kevin", "bob", "important", "this is an example body", queueSecretKey)
	if err != nil {
		b.Errorf("Unable to create msg example due to: %q", err)
	}

	for b.Loop() {
		queue := NewQueue()
		for range 1000 {
			queue.Enqueue(*msgExample)
		}
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	msgExample, err := NewMessage("kevin", "bob", "important", "this is an example body", queueSecretKey)
	if err != nil {
		b.Errorf("Unable to create msg example due to: %q", err)
	}

	for b.Loop() {
		// not measuring queue setup
		b.StopTimer()
		queue := NewQueue()
		for range 1000 {
			queue.Enqueue(*msgExample)
		}

		// measuring dequeue explicitly
		b.StartTimer()
		for range 1000 {
			_, ok := queue.Dequeue()
			if !ok {
				b.Error("Unable to dequeue")
			}
		}
	}
}
