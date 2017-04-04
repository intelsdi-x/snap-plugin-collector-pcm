package parser

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"

	"math"

	"github.com/pkg/errors"
)

// When true, parser behaves like old one, it's returning last cached line
// instead waiting for next one.
var CompatibilityMode = true

type (
	Key struct {
		Component, MetricName string
	}

	Values map[Key]float64

	ValuesOrError struct {
		Values Values
		Error  error
	}

	Parser interface {
		GetKeys(ctx context.Context) ([]Key, error)
		GetValues(ctx context.Context) (Values, error)
	}
)

func RunParser(reader io.Reader) Parser {
	p := &pcmParser{
		source:          reader,
		keysReady:       make(chan struct{}),
		keysInfoMutex:   new(sync.RWMutex),
		streamInfoReady: make(chan struct{}),
		streamInfoMutex: new(sync.RWMutex),
	}
	go p.run()
	return p
}

func RunStreamedParser(reader io.Reader, chanLen int) (Parser, <-chan ValuesOrError) {
	p := &pcmParser{
		source:        reader,
		sink:          make(chan ValuesOrError, chanLen),
		keysReady:     make(chan struct{}),
		keysInfoMutex: new(sync.RWMutex),
	}
	go p.run()

	return p, p.sink
}

func (p *pcmParser) GetKeys(ctx context.Context) ([]Key, error) {
	var sig chan struct{}
	withRLock(p.keysInfoMutex, func() {
		sig = p.keysReady
	})
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-sig:
			var keys []Key
			withRLock(p.keysInfoMutex, func() {
				keys = p.keys
			})
			return keys, nil
		}
	}
}

func (p *pcmParser) GetValues(ctx context.Context) (Values, error) {
	if p.sink != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case vals := <-p.sink:
			return vals.Values, vals.Error
		}
	}

	var sig chan struct{}
	withRLock(p.streamInfoMutex, func() {
		sig = p.streamInfoReady
	})

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-sig:
			var sVal ValuesOrError
			withRLock(p.streamInfoMutex, func() {
				sVal = p.stream
			})
			return sVal.Values, sVal.Error
		}
	}
}

type pcmParser struct {
	source io.Reader
	sink   chan ValuesOrError

	keys          []Key
	keysReady     chan struct{}
	keysInfoMutex *sync.RWMutex

	stream          ValuesOrError
	streamInfoReady chan struct{}
	streamInfoMutex *sync.RWMutex
}

func (p *pcmParser) run() {
	scanner := bufio.NewScanner(p.source)
	var first, second, current []string
	streamInfoClosed := false
	for scanner.Scan() {
		if first == nil {
			first = splitLine(scanner.Text())
			fillHeader(first)
			first = first[2:]
			continue
		}
		if second == nil {
			second = splitLine(scanner.Text())
			second = second[2:]
			if len(first) != len(second) {
				panic("data length changed")
			}

			withLock(p.keysInfoMutex, func() {
				p.keys = make([]Key, len(first))
				for i, topHeader := range first {
					p.keys[i] = Key{Component: topHeader, MetricName: second[i]}
				}
				close(p.keysReady)
			})
			continue
		}

		current = splitLine(scanner.Text())
		current = current[2:]
		if len(first) != len(current) {
			panic("data length changed")
		}

		vals := ValuesOrError{Values: Values{}}
		for i, field := range current {
			v, err := strconv.ParseFloat(field, 64)
			k := Key{Component: first[i], MetricName: second[i]}
			if err == nil {
				vals.Values[k] = v
			} else {
				if strings.ToLower(field) == "n/a" {
					//TODO: make sure this is desired, maybe entry should be just missing
					vals.Values[k] = math.NaN()
				} else {
					vals = ValuesOrError{Error: errors.Wrapf(err, "parsing %v = %v failed", k, field)}
					streamInfoClosed = true
					break
				}
			}
		}
		if p.sink == nil {
			withLock(p.streamInfoMutex, func() {
				orgSig := p.streamInfoReady
				if !CompatibilityMode {
					p.streamInfoReady = make(chan struct{})
				}
				p.stream = vals
				if !CompatibilityMode || !streamInfoClosed {
					close(orgSig)
					streamInfoClosed = true
				}
			})
		} else {
			p.sink <- vals
		}
	}
	if p.sink == nil {
		if !CompatibilityMode {
			withLock(p.streamInfoMutex, func() {
				orgSig := p.streamInfoReady
				p.stream = ValuesOrError{Error: errors.New("stream not running")}
				close(orgSig)
				streamInfoClosed = true

			})
		}
	} else {
		close(p.sink)
	}
}

func fillHeader(headerRef []string) {
	for i := 1; i < len(headerRef); i++ {
		if len(headerRef[i]) == 0 {
			headerRef[i] = headerRef[i-1]
		}
	}
}

func splitLine(line string) []string {
	line = strings.TrimSuffix(line, ";")
	split := strings.Split(line, ";")
	for i, field := range split {
		split[i] = strings.TrimSpace(field)
	}

	return split
}

func withRLock(rwMutex *sync.RWMutex, f func()) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	f()
}

func withLock(mutex *sync.RWMutex, f func()) {
	mutex.Lock()
	defer mutex.Unlock()
	f()
}