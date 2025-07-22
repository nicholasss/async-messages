package msg

import (
	"bytes"
	"errors"
	"testing"
)

var rawMsgSecretKey = []byte("GgfY0UssupyYBlFy92/ENsq5/Qy8dq3bh3Mp8hZcPMDEdSnxMgi5E1TPzJuHVHzRs60aq6r7gKyLGwbauaUn1Q==")

func TestToPackagedMessage(t *testing.T) {
	tt := []struct {
		rawMsg    RawMessage
		wantError *MissingFieldError
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
			wantError: nil,
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
			wantError: nil,
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
			wantError: nil,
		},
		{
			rawMsg: RawMessage{
				ToName:     "",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantError: &MissingFieldError{Field: "ToName"},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantError: &MissingFieldError{Field: "ToVessel"},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantError: &MissingFieldError{Field: "FromName"},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantError: &MissingFieldError{Field: "FromVessel"},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantError: &MissingFieldError{Field: "Subject"},
		},
		{
			rawMsg: RawMessage{
				ToName:     "Bob",
				ToVessel:   "Snow",
				FromName:   "Kevin",
				FromVessel: "Liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "",
			},
			wantError: &MissingFieldError{Field: "Body"},
		},
	}

	for _, tc := range tt {
		_, gotErr := tc.rawMsg.ToPackagedMessage(rawMsgSecretKey)

		if tc.wantError != nil {
			//
			// we are expecting a specific error
			var gotAsError *MissingFieldError
			if !errors.As(gotErr, &gotAsError) {
				t.Errorf("Expected a MissingFieldError, but got %v (type %T)", gotErr, gotErr)
			} else {
				// compare the fields to ensure they match as well
				if gotAsError.Field != tc.wantError.Field {
					t.Errorf("MissingFieldError field mismatch: got=%q, want=%q", gotAsError.Field, tc.wantError.Field)
				}
			}
		} else {
			//
			// we do not expect an error
			if gotErr != nil {
				t.Errorf("Did not expect error: got=%q", gotErr)
			}
		}
	}
}

func TestRawMessageDataForSigning(t *testing.T) {
	tt := []struct {
		rawMsg          RawMessage
		wantMessageData []byte
	}{
		{
			rawMsg: RawMessage{
				ToName:     "bob",
				ToVessel:   "snow",
				FromName:   "kevin",
				FromVessel: "liberty",
				Subject:    "Tuesday",
				Body:       "I am planning on proceeding on tuesday since there is a break in the weather",
			},
			wantMessageData: []byte("bob@snow|kevin@liberty|Tuesday|I am planning on proceeding on tuesday since there is a break in the weather"),
		},
		{
			rawMsg: RawMessage{
				ToName:     "kevin",
				ToVessel:   "liberty",
				FromName:   "bob",
				FromVessel: "snow",
				Subject:    "Re: Tuesday",
				Body:       "I will need to wait longer because of needed repair work. Hope to catch up.",
			},
			wantMessageData: []byte("kevin@liberty|bob@snow|Re: Tuesday|I will need to wait longer because of needed repair work. Hope to catch up."),
		},
		{
			rawMsg: RawMessage{
				ToName:     "bob",
				ToVessel:   "snow",
				FromName:   "kevin",
				FromVessel: "liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantMessageData: []byte("bob@snow|kevin@liberty|Re: Re: Tuesday|Good idea. Take your time. I will send out a message when we get into Cambridge Bay."),
		},
	}

	for _, tc := range tt {
		gotMessageData := tc.rawMsg.messageDataForSigning()

		if !bytes.Equal(gotMessageData, tc.wantMessageData) {
			t.Errorf("message data mismatch. got=%s want=%s", gotMessageData, tc.wantMessageData)
		}
	}
}

func TestCreateSignature(t *testing.T) {
	tt := []struct {
		rawMsg        RawMessage
		wantSignature string
	}{
		{
			rawMsg: RawMessage{
				ToName:     "bob",
				ToVessel:   "snow",
				FromName:   "kevin",
				FromVessel: "liberty",
				Subject:    "Tuesday",
				Body:       "I am planning on proceeding on tuesday since there is a break in the weather",
			},
			wantSignature: "2d18d51dee10b1c5b5e673652b9744ebc8f0900c114cb34262b1429738a116f9",
		},
		{
			rawMsg: RawMessage{
				ToName:     "kevin",
				ToVessel:   "liberty",
				FromName:   "bob",
				FromVessel: "snow",
				Subject:    "Re: Tuesday",
				Body:       "I will need to wait longer because of needed repair work. Hope to catch up.",
			},
			wantSignature: "7bf816982b5d8c52822a366e68d9fd2662c9192e4a86473d4dfdfbd29cfd64d5",
		},
		{
			rawMsg: RawMessage{
				ToName:     "bob",
				ToVessel:   "snow",
				FromName:   "kevin",
				FromVessel: "liberty",
				Subject:    "Re: Re: Tuesday",
				Body:       "Good idea. Take your time. I will send out a message when we get into Cambridge Bay.",
			},
			wantSignature: "1a936de118b68b1cea9469ec5cd22b7bd65bef71027eb35b88c931ad66eaa870",
		},
	}

	for _, tc := range tt {
		gotSignature, err := tc.rawMsg.createSignature(rawMsgSecretKey)
		if err != nil {
			t.Errorf("create signature failed unexpectedly due to: %q", err)
		}

		if gotSignature != tc.wantSignature {
			t.Errorf("signature mismatch. got=%s want=%s", gotSignature, tc.wantSignature)
		}
	}
}
