package msg

import (
	"bytes"
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
				Signature: "e9d99039e92f5f0adb162e2f9fa44e3f61356df0",
				// not real signatures
			},
			wantString: "To: Bob@Snow\nFrom: Kevin@Liberty\nSubject: Strong currents ahead\nBody: There are strong currents ahead, ~5nm NE of our position.\nSignature: e9d99039e92f5f0adb162e2f9fa44e3f61356df0\n",
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
				Signature: "e9d99039e92f5f0adb162e2f9fa44e3f61356df0",
				// not real signatures
			},
			wantString: "To: Kevin@Liberty\nFrom: Bob@Snow\nSubject: Roger\nBody: We will slow to below 4kt and watch our drifting. Thanks.\nSignature: e9d99039e92f5f0adb162e2f9fa44e3f61356df0\n",
		},
	}

	for _, tc := range tt {
		gotString := tc.msg.String()
		if gotString != tc.wantString {
			t.Errorf("stringified version of packaged message not equal. got=%q want=%q", gotString, tc.wantString)
		}
	}
}

func TestPackagedMessageDataForSigning(t *testing.T) {
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
		pkgMsg, err := tc.rawMsg.ToPackagedMessage(pkgMsgSecretKey)
		if err != nil {
			t.Errorf("failed to package message due to: %q", err)
		}

		gotMessageData := pkgMsg.messageDataForSigning()

		if !bytes.Equal(gotMessageData, tc.wantMessageData) {
			t.Errorf("message data mismatch. got=%s want=%s", gotMessageData, tc.wantMessageData)
		}
	}
}
