package runtime

import (
	"math"
	"testing"
	"time"
)

func TestTimeoutDecode(t *testing.T) {
	for _, spec := range []struct {
		timeout string
		want    time.Duration
		wantErr bool
	}{
		{timeout: "17H", want: 17 * time.Hour},
		{timeout: "19M", want: 19 * time.Minute},
		{timeout: "23S", want: 23 * time.Second},
		{timeout: "1009m", want: 1009 * time.Millisecond},
		{timeout: "1000003u", want: 1000003 * time.Microsecond},
		{timeout: "100000007n", want: 100000007 * time.Nanosecond},
		// Values that would overflow time.Duration (an int64 nanosecond count)
		// must be clamped to the maximum representable duration rather than
		// wrapping to a bogus, often negative, deadline.
		{timeout: "99999999H", want: time.Duration(math.MaxInt64)},
		{timeout: "3000000H", want: time.Duration(math.MaxInt64)},
		{timeout: "999999999M", want: time.Duration(math.MaxInt64)},
		{timeout: "10000000000000000000n", want: time.Duration(math.MaxInt64)},
		{timeout: "-5H", wantErr: true},
		{timeout: "H", wantErr: true},
	} {
		got, err := timeoutDecode(spec.timeout)
		if spec.wantErr {
			if err == nil {
				t.Errorf("timeoutDecode(%q) = %v, nil; want error", spec.timeout, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("timeoutDecode(%q) failed with %v; want success", spec.timeout, err)
			continue
		}
		if got != spec.want {
			t.Errorf("timeoutDecode(%q) = %v; want %v", spec.timeout, got, spec.want)
		}
		if got < 0 {
			t.Errorf("timeoutDecode(%q) = %v; want non-negative duration", spec.timeout, got)
		}
	}
}
