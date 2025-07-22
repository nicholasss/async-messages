package msg

import (
	"testing"
)

var pkgMsgSecretKey = []byte("GgfY0UssupyYBlFy92/ENsq5/Qy8dq3bh3Mp8hZcPMDEdSnxMgi5E1TPzJuHVHzRs60aq6r7gKyLGwbauaUn1Q==")

// tests the String() method of a packaged message
// a packaged message should already have valid fields
func TestStringOfPackagedMessage(t *testing.T) {
	tt := []struct {
		msg        PackagedMessage
		wantString string
	}{
		{
			msg: PackagedMessage{
				To: UserVessel{
					Name:   "Bob",
					Vessel: "Snow",
				},
				From: UserVessel{
					Name:   "Kevin",
					Vessel: "Liberty",
				},
				Subject:   "Strong currents ahead",
				Body:      "There are strong currents ahead, ~5nm NE of our position.",
				Signature: "",
			},
			wantString: "To: Bob@Snow\nFrom: Kevin@Liberty\nSubject: Strong currents ahead\nBody: There are strong currents ahead, ~5nm NE of our position.\nSignature: \n",
		},
		{
			msg: PackagedMessage{
				To: UserVessel{
					Name:   "Kevin",
					Vessel: "Liberty",
				},
				From: UserVessel{
					Name:   "Bob",
					Vessel: "Snow",
				},
				Subject:   "Roger",
				Body:      "We will slow to below 4kt and watch our drifting. Thanks.",
				Signature: "",
			},
			wantString: "To: Kevin@Liberty\nFrom: Bob@Snow\nSubject: Roger\nBody: We will slow to below 4kt and watch our drifting. Thanks.\nSignature: \n",
		},
	}

	for _, tc := range tt {
		gotString := tc.msg.String()
		if gotString != tc.wantString {
			t.Errorf("stringified version of packaged message not equal. got=%q want=%q", gotString, tc.wantString)
		}
	}
}

// test that the function is successfully creating the message
// func TestCreateMessage(t *testing.T) {
// 	tt := []struct {
// 		caseName string
// 		to       string
// 		from     string
// 		subject  string
// 		body     string
// 	}{
// 		{
// 			caseName: "All fields present",
// 			to:       "Bob",
// 			from:     "Kevin",
// 			subject:  "Re: Tuesday",
// 			body:     "This is super important that we work on this tuesday.",
// 		},
// 		{
// 			caseName: "All fields present",
// 			to:       "Kevin",
// 			from:     "Bob",
// 			subject:  "Re: Re: Tuesday",
// 			body:     "I am not sure that I can do tuesday though.",
// 		},
// 		{
// 			caseName: "All fields present",
// 			to:       "Bob",
// 			from:     "Kevin",
// 			subject:  "Re: Re: Re: Tuesday",
// 			body:     "You can do tue, just let me know when on tue. I am free after lunch.",
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		msg, err := NewMessage(tc.to, tc.from, tc.subject, tc.body, msgSecretKey)
// 		if err != nil {
// 			t.Errorf("case: %q\nwas not able to create message due to: %q", tc.caseName, err)
// 		}
// 		if msg.To != tc.to {
// 			t.Errorf("case: %q\ndid not copy 'to' field. got %q want %q", tc.caseName, msg.To, tc.to)
// 		}
// 		if msg.From != tc.from {
// 			t.Errorf("case: %q\ndid not copy 'from' field. got %q want %q", tc.caseName, msg.From, tc.from)
// 		}
// 		if msg.Subject != tc.subject {
// 			t.Errorf("case: %q\ndid not copy 'subject' field. got %q want %q", tc.caseName, msg.Subject, tc.subject)
// 		}
// 		if msg.Body != tc.body {
// 			t.Errorf("case: %q\n did not copy 'body' field. got %q want %q", tc.caseName, msg.Body, tc.body)
// 		}
// 	}
// }
//
// // tests whether the ToString() function works as intended
// func TestToString(t *testing.T) {
// 	tt := []struct {
// 		to           string
// 		from         string
// 		subject      string
// 		body         string
// 		toStringWant string
// 	}{
// 		{
// 			to:           "Bob",
// 			from:         "Kevin",
// 			subject:      "Re: Tuesday",
// 			body:         "This is super important that we work on this tuesday.",
// 			toStringWant: "To: Bob\nFrom: Kevin\nSubject: Re: Tuesday\nBody: This is super important that we work on this tuesday.\nSignature: 6786dea088cb7648a016f5d3736e9ee3521d7f42a4b402e063da64d51ba4d4c0\n",
// 		},
// 		{
// 			to:           "Kevin",
// 			from:         "Bob",
// 			subject:      "Re: Re: Tuesday",
// 			body:         "I am not sure that I can do tuesday though.",
// 			toStringWant: "To: Kevin\nFrom: Bob\nSubject: Re: Re: Tuesday\nBody: I am not sure that I can do tuesday though.\nSignature: dd5c02b3e521b5a855e7fe85806f2ad2a66948e5569c31a32bff50843530113b\n",
// 		},
// 		{
// 			to:           "Bob",
// 			from:         "Kevin",
// 			subject:      "Re: Re: Re: Tuesday",
// 			body:         "You can do tue, just let me know when on tue. I am free after lunch.",
// 			toStringWant: "To: Bob\nFrom: Kevin\nSubject: Re: Re: Re: Tuesday\nBody: You can do tue, just let me know when on tue. I am free after lunch.\nSignature: bcf00b7f415f54ff1391ba5b899d39b4405e3d9bb04293d62c907bf961592c73\n",
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		msg, err := NewMessage(tc.to, tc.from, tc.subject, tc.body, msgSecretKey)
// 		if err != nil {
// 			t.Errorf("was not able to create message due to: %q", err)
// 		}
//
// 		toStringGot := msg.ToString()
// 		if toStringGot != tc.toStringWant {
// 			t.Errorf("was not able to convert message to string successfully. got=%q, want=%q", toStringGot, tc.toStringWant)
// 		}
// 	}
// }
//
// // tests whether the function errors when there is a missing field
// func TestEmptyPropertyMessages(t *testing.T) {
// 	tt := []struct {
// 		caseName string
// 		to       string
// 		from     string
// 		subject  string
// 		body     string
// 	}{
// 		{
// 			caseName: "Missing 'to' field",
// 			to:       "",
// 			from:     "Kevin",
// 			subject:  "To anyone out there...",
// 			body:     "Please pass the salt.",
// 		},
// 		{
// 			caseName: "Missing 'from' field",
// 			to:       "Bob",
// 			from:     "",
// 			subject:  "To anyone out there...",
// 			body:     "Please pass the salt.",
// 		},
// 		{
// 			caseName: "Missing 'subject' field",
// 			to:       "Bob",
// 			from:     "Kevin",
// 			subject:  "",
// 			body:     "Please pass the salt.",
// 		},
// 		{
// 			caseName: "Missing 'body' field",
// 			to:       "Bob",
// 			from:     "Kevin",
// 			subject:  "To anyone out there...",
// 			body:     "",
// 		},
// 		{
// 			caseName: "",
// 			to:       "",
// 			from:     "",
// 			subject:  "",
// 			body:     "",
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		_, err := NewMessage(tc.to, tc.from, tc.subject, tc.body, msgSecretKey)
// 		if err == nil {
// 			t.Errorf("case: %q\nwas not able to create message due to: %q", tc.caseName, err)
// 		}
// 	}
// }
//
// // tests whether the signature is being computed correctly
// func TestVerifySignatures(t *testing.T) {
// 	tt := []struct {
// 		caseName  string
// 		to        string
// 		from      string
// 		subject   string
// 		body      string
// 		signature string
// 	}{
// 		{
// 			caseName:  "testing that it works correctly",
// 			to:        "Bob",
// 			from:      "Kevin",
// 			subject:   "Re: Tuesday",
// 			body:      "This is super important that we work on this tuesday.",
// 			signature: "6786dea088cb7648a016f5d3736e9ee3521d7f42a4b402e063da64d51ba4d4c0",
// 		},
// 		{
// 			caseName:  "testing that it works correctly",
// 			to:        "Kevin",
// 			from:      "Bob",
// 			subject:   "Re: Re: Tuesday",
// 			body:      "I am not sure that I can do tuesday though.",
// 			signature: "dd5c02b3e521b5a855e7fe85806f2ad2a66948e5569c31a32bff50843530113b",
// 		},
// 		{
// 			caseName:  "testing that it works correctly",
// 			to:        "Bob",
// 			from:      "Kevin",
// 			subject:   "Re: Re: Re: Tuesday",
// 			body:      "You can do tue, just let me know when on tue. I am free after lunch.",
// 			signature: "bcf00b7f415f54ff1391ba5b899d39b4405e3d9bb04293d62c907bf961592c73",
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		msg, err := NewMessage(tc.to, tc.from, tc.subject, tc.body, msgSecretKey)
// 		if err != nil {
// 			fmt.Printf("Should work, but didn't. Got error: %q", err)
// 			t.Fail()
// 		}
//
// 		err = msg.verifySignature(msgSecretKey)
// 		if err != nil {
// 			fmt.Printf("Should work, but didn't. Got error: %q", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestPrepMessageForSigning(t *testing.T) {
// 	tt := []struct {
// 		caseName   string
// 		to         string
// 		from       string
// 		subject    string
// 		body       string
// 		preppedMsg []byte
// 	}{
// 		{
// 			caseName:   "testing concat",
// 			to:         "Bob",
// 			from:       "Kevin",
// 			subject:    "Re: Tuesday",
// 			body:       "This is super important that we work on this tuesday.",
// 			preppedMsg: []byte("Bob|Kevin|Re: Tuesday|This is super important that we work on this tuesday."),
// 		},
// 		{
// 			caseName:   "testing concat",
// 			to:         "Kevin",
// 			from:       "Bob",
// 			subject:    "Re: Re: Tuesday",
// 			body:       "I am not sure that I can do tuesday though.",
// 			preppedMsg: []byte("Kevin|Bob|Re: Re: Tuesday|I am not sure that I can do tuesday though."),
// 		},
// 		{
// 			caseName:   "testing conact",
// 			to:         "Bob",
// 			from:       "Kevin",
// 			subject:    "Re: Re: Re: Tuesday",
// 			body:       "You can do tue, just let me know when on tue. I am free after lunch.",
// 			preppedMsg: []byte("Bob|Kevin|Re: Re: Re: Tuesday|You can do tue, just let me know when on tue. I am free after lunch."),
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		msg := &RawMessage{To: tc.to, From: tc.from, Subject: tc.subject, Body: tc.body}
// 		msgData := msg.messageDataForSigning()
//
// 		if !bytes.Equal(tc.preppedMsg, msgData) {
// 			t.Errorf("case: %q\ndid not prepare message correctly. got %q, want %q", tc.caseName, msgData, tc.preppedMsg)
// 		}
// 	}
// }
//
// func TestMessagePackageTime(t *testing.T) {
// 	tt := []struct {
// 		to      string
// 		from    string
// 		subject string
// 		body    string
// 	}{
// 		{
// 			to:      "Bob",
// 			from:    "Kevin",
// 			subject: "Re: Tuesday",
// 			body:    "This is super important that we work on this tuesday.",
// 		},
// 		{
// 			to:      "Kevin",
// 			from:    "Bob",
// 			subject: "Re: Re: Tuesday",
// 			body:    "I am not sure that I can do tuesday though.",
// 		},
// 		{
// 			to:      "Bob",
// 			from:    "Kevin",
// 			subject: "Re: Re: Re: Tuesday",
// 			body:    "You can do tue, just let me know when on tue. I am free after lunch.",
// 		},
// 	}
//
// 	for _, tc := range tt {
// 		msg, err := NewMessage(tc.to, tc.from, tc.subject, tc.body, msgSecretKey)
// 		if err != nil {
// 			t.Errorf("was not able to create message due to: %q", err)
// 		}
//
// 		msgPackagedGot := msg.Packaged
// 		oneMinuteAgoWantAfter := time.Now().UTC().Add(-1 * time.Minute)
//
// 		if msgPackagedGot.Before(oneMinuteAgoWantAfter) {
// 			t.Errorf("packaged time is before one minute ago. got=%q want=%q", msgPackagedGot.UTC().String(), oneMinuteAgoWantAfter.UTC().String())
// 		}
//
// 		oneMinuteInFutureWantBefore := time.Now().UTC().Add(time.Minute)
//
// 		if msgPackagedGot.After(oneMinuteInFutureWantBefore) {
// 			t.Errorf("packaged time is after one minute in future. got=%q want=%q", msgPackagedGot.UTC().String(), oneMinuteInFutureWantBefore.UTC().String())
// 		}
// 	}
// }
