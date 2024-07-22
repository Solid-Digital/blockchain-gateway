package adapter

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"

	"github.com/unchainio/pkg/xsync"
)

var globalCounter xsync.Counter

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewTag(optFuncs ...MessageOptsFunc) string {
	opts := defaultOpts

	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}

	if opts.tag == "" {
		opts.tag = fmt.Sprintf("%d", globalCounter.Add(1))
	}

	return opts.tag
}

type MessageOpts struct {
	tag string
}

var defaultOpts = MessageOpts{}

type MessageOptsFunc func(opt *MessageOpts)

func randomTag() string {
	return fmt.Sprintf("%d", rand.Uint64())
}

func WithTag(tag string) MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.tag = tag
	}
}
func WithUUID() MessageOptsFunc {
	return func(opts *MessageOpts) {
		tag, _ := uuid.NewV4()

		opts.tag = tag.String()
	}
}

func WithRandomTag() MessageOptsFunc {
	return func(opts *MessageOpts) {
		opts.tag = randomTag()
	}
}
