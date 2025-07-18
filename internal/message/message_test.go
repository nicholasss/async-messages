package message

import (
	"fmt"
	"testing"
)

var secretKey = []byte("GgfY0UssupyYBlFy92/ENsq5/Qy8dq3bh3Mp8hZcPMDEdSnxMgi5E1TPzJuHVHzRs60aq6r7gKyLGwbauaUn1Q==")

// test that the messages is being composed and that signatures are being written
func TestCreateMessage(t *testing.T) {
	tt := []struct {
		caseName string
		to       string
		from     string
		subject  string
		body     string
	}{
		{
			caseName: "All fields present",
			to:       "Bob",
			from:     "Kevin",
			subject:  "Re: Tuesday",
			body:     "This is super important that we work on this tuesday.",
		},
		{
			caseName: "All fields present",
			to:       "Kevin",
			from:     "Bob",
			subject:  "Re: Re: Tuesday",
			body:     "I am not sure that I can do tuesday though.",
		},
		{
			caseName: "All fields present",
			to:       "Bob",
			from:     "Kevin",
			subject:  "Re: Re: Re: Tuesday",
			body:     "You can do tue, just let me know when on tue. I am free after lunch.",
		},
		{
			caseName: "Missing 'to' field",
			to:       "",
			from:     "Kevin",
			subject:  "To anyone out there...",
			body:     "Please pass the salt.",
		},
	}

	for _, tc := range tt {
		msg, err := CreateMessage(tc.to, tc.from, tc.subject, tc.body, secretKey)
		if err != nil {
			t.Errorf("case: %q\nwas not able to create message due to: %q", tc.caseName, err)
		}
		if msg.To != tc.to {
			t.Errorf("case: %q\ndid not copy 'to' field. got %q want %q", tc.caseName, msg.To, tc.to)
		}
		if msg.From != tc.from {
			t.Errorf("case: %q\ndid not copy 'from' field. got %q want %q", tc.caseName, msg.From, tc.from)
		}
		if msg.Subject != tc.subject {
			t.Errorf("case: %q\ndid not copy 'subject' field. got %q want %q", tc.caseName, msg.Subject, tc.subject)
		}
		if msg.Body != tc.body {
			t.Errorf("case: %q\n did not copy 'body' field. got %q want %q", tc.caseName, msg.Body, tc.body)
		}
	}
}

func TestVerifySignatures(t *testing.T) {
	tt := []struct {
		name      string
		to        string
		from      string
		subject   string
		body      string
		signature string
	}{
		{
			name: "testing that it works correctly",
			to:   "Bob", from: "Kevin",
			subject:   "Re: Tuesday",
			body:      "This is super important that we work on this tuesday.",
			signature: "6786dea088cb7648a016f5d3736e9ee3521d7f42a4b402e063da64d51ba4d4c0",
		},
		{
			name: "testing that it works correctly",
			to:   "Kevin", from: "Bob",
			subject:   "Re: Re: Tuesday",
			body:      "I am not sure that I can do tuesday though.",
			signature: "dd5c02b3e521b5a855e7fe85806f2ad2a66948e5569c31a32bff50843530113b",
		},
		{
			name: "testing that it works correctly",
			to:   "Bob", from: "Kevin",
			subject:   "Re: Re: Re: Tuesday",
			body:      "You can do tue, just let me know when on tue. I am free after lunch.",
			signature: "bcf00b7f415f54ff1391ba5b899d39b4405e3d9bb04293d62c907bf961592c73",
		},
	}

	for i, tc := range tt {
		msg, err := CreateMessage(tc.to, tc.from, tc.subject, tc.body, secretKey)
		if err != nil {
			fmt.Printf("- [%d] Should work, but didn't. Got error: %q", i, err)
			t.Fail()
		}

		_, err = verifySignature(msg, secretKey)
		if err != nil {
			fmt.Printf("- [%d] Should work, but didn't. Got error: %q", i, err)
			t.Fail()
		}
	}
}
