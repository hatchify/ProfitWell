package profitwell

import (
	"os"
	"testing"

	"github.com/Hatch1fy/errors"
)

var testAuthToken = os.Getenv("PROFITWELL_AUTH_TOKEN")

func TestNew(t *testing.T) {
	if _, err := New(testAuthToken); err != nil {
		t.Fatal(err)
	}
}

func TestProfitWell_SetUserAction(t *testing.T) {
	var (
		p   *ProfitWell
		err error
	)

	if p, err = New(testAuthToken); err != nil {
		t.Fatal(err)
	}

	tcs := []testCase{
		newTestCase("test@hatchify.co", nil),
		newTestCase("", errors.Error("Customer email or user ID not provided!")),
	}

	for _, tc := range tcs {
		if err = p.SetUserAction(tc.userEmail); !tc.isOK(err) {
			t.Fatalf("invalid error, expected %v and received %v", tc.err, err)
		}
	}
}

func newTestCase(userEmail string, err error) (t testCase) {
	t.userEmail = userEmail
	t.err = err
	return
}

type testCase struct {
	userEmail string
	err       error
}

func (t *testCase) isOK(err error) (ok bool) {
	if t.err == nil && err == nil {
		return true
	}

	if t.err == nil && err != nil {
		return
	}

	if t.err != nil && err == nil {
		return
	}

	return t.err.Error() == err.Error()
}
