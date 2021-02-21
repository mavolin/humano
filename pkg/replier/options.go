package replier

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type (
	// Options holds the configuration of the replier.
	Options struct {
		// DefaultDelayFunc is the DelayFunc used if no replacement was
		// provided.
		//
		// Default: RandomizedDelay(500*time.Millesecond, 0.1)
		DefaultDelayFunc DelayFunc
		// Splitter is the Splitter used to split messages.
		//
		// Default: NoSplit
		Splitter Splitter
		// NoTyping specifies that no typing event shall be sent.
		//
		// Default: false
		NoTyping bool
	}

	// DelayFunc is the function type used to calculate how long a specific
	// message shall be delayed.
	DelayFunc func(message string) time.Duration
	// Splitter is the function used to split messages into smaller ones.
	Splitter func(message string) []string
)

func (o *Options) fillDefaults() {
	if o.DefaultDelayFunc == nil {
		o.DefaultDelayFunc = RandomizedDelay(500*time.Millisecond, 0.1)
	}

	if o.Splitter == nil {
		o.Splitter = NoSplit
	}
}

// NoSplit is a Splitter that never splits a message.
func NoSplit(message string) []string { return []string{message} }

var _ Splitter = NoSplit

// FieldsFuncSplitter creates a new Splitter that splits, as if it was handed
// to strings.FieldsFunc.
func FieldsFuncSplitter(f func(rune) bool) Splitter {
	return func(message string) []string {
		return strings.FieldsFunc(message, f)
	}
}

// StaticDelay is a DelayFunc that always returns the same delay.
func StaticDelay(d time.Duration) DelayFunc {
	return func(string) time.Duration { return d }
}

// RandomizedDelay returns a delay that consists of a fixed character-based
// delay and an additional randomized delay constructed randomFactor.
// randomFactor is a number between 0 and 1.0.
//
// It will be used to generate a random multiplier within the range of
// randomFactor, which is then used to create (or substract) an extra delay.
// This means the returned delay will be (len(message) * runeDelay +
// len(message) * runeDelay * randomFactor).
func RandomizedDelay(runeDelay time.Duration, randomFactor float64) DelayFunc {
	return func(message string) time.Duration {
		d := float64(int(runeDelay) * len(message))
		d += d * ((2 * rand.Float64() * randomFactor) - randomFactor)

		return time.Duration(d)
	}
}
