package expbackoff

import (
	"math"
	"math/rand"
	"time"
)

//Backoff sleeps a random amount so we can.
//retry failed requests using a randomized exponential backoff:
// Example:
// Backoff(try, 0, 16, time.Second)
// try=1: wait a random period between [0..1] seconds and retry; if that fails,
// try=2: wait a random period between [0..2] seconds and retry; if that fails,
// try=3: wait a random period between [0..4] seconds and retry; if that fails,
// try=4: wait a random period between [0..8] seconds and retry, and so on,
// try=N: with an upper bounds to the wait period being 16 seconds.
//https://play.golang.org/p/O-PjlWl-zBS
func Backoff(try, min, max int, duration time.Duration) {
	end := EndingBounds(try, min, max, duration)
	d := RandomDuration(min, end, duration)
	time.Sleep(d)
}

func RandomDuration(start, end int, duration time.Duration) time.Duration {
	r := rand.Int31n(int32(end))
	d := time.Duration(r+int32(start)) * duration
	return d
}

func EndingBounds(try, min, max int, duration time.Duration) int {
	nf := math.Pow(2, float64(try))
	nf = math.Max(float64(min), nf)
	nf = math.Min(nf, float64(max))
	return int(nf)
}
