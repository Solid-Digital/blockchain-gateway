package domain

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/unchainio/pkg/xsync"
	"math/rand"
	"time"
)

type Request struct {
	Tag    string
	Output Output
	Error  error
}

type Output map[string]interface{}

func NewRequest(output Output) *Request {
	return &Request{
		Tag:    NewTag(),
		Output: output,
	}
}


func NewRequestError(err error) *Request {
	return &Request{
		Tag:   NewTag(),
		Error: err,
	}
}


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
