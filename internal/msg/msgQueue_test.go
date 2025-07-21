package msg

import "testing"

var queueSecretKey = []byte("GgfY0UssupyYBlFy92/ENsq5/Qy8dq3bh3Mp8hZcPMDEdSnxMgi5E1TPzJuHVHzRs60aq6r7gKyLGwbauaUn1Q==")

func TestSizeFunction(t *testing.T) {
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

	// should be empty on creation
	queue := NewQueue()
	if queue.IsEmpty() != true {
		t.Errorf("Queue should be empty on creation")
	}
	if queue.Size() != 0 {
		t.Errorf("Queue should be emtpy on creation")
	}

	queue.Enqueue(packagedMsgs[0])
	if queue.IsEmpty() != false {
		t.Errorf("Queue should not be empty after enqueue")
	}
	if queue.Size() != 1 {
		t.Errorf("Queue should contain one item after enqueue")
	}
}
