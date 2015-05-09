package runtime

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gengo/grpc-gateway/internal"
	"github.com/golang/glog"
)

var (
	// ErrNotMatch indicates that the given HTTP request path does not match to the pattern.
	ErrNotMatch = errors.New("not match to the path pattern")
	// ErrInvalidPattern indicates that the given definition of Pattern is not valid.
	ErrInvalidPattern = errors.New("invalid pattern")
)

type op struct {
	code    internal.OpCode
	operand int
}

// Pattern is a template pattern of http request paths defined in third_party/googleapis/google/api/http.proto.
type Pattern struct {
	// ops is a list of operations
	ops []op
	// pool is a constant pool indexed by the operands or vars.
	pool []string
	// vars is a list of variables names to be bound by this pattern
	vars []string
	// stacksize is the max depth of the stack
	stacksize int
	// verb is the VERB part of the path pattern. It is empty if the pattern does not have VERB part.
	verb string
}

// NewPattern returns a new Pattern from the given definition values.
// "ops" is a sequence of op codes. "pool" is a constant pool.
// "verb" is the verb part of the pattern. It is empty if the pattern does not have the part.
// "version" must be 1 for now.
// It returns an error if the given definition is invalid.
func NewPattern(version int, ops []int, pool []string, verb string) (Pattern, error) {
	if version != 1 {
		glog.V(2).Infof("unsupported version: %d", version)
		return Pattern{}, ErrInvalidPattern
	}

	l := len(ops)
	if l%2 != 0 {
		glog.V(2).Infof("odd number of ops codes: %d", l)
		return Pattern{}, ErrInvalidPattern
	}

	var typedOps []op
	var stack, maxstack int
	var vars []string
	for i := 0; i < l; i += 2 {
		op := op{code: internal.OpCode(ops[i]), operand: ops[i+1]}
		switch op.code {
		case internal.OpNop:
			continue
		case internal.OpPush, internal.OpPushM:
			stack++
		case internal.OpLitPush:
			if op.operand < 0 || len(pool) <= op.operand {
				glog.V(2).Infof("negative literal index: %d", op.operand)
				return Pattern{}, ErrInvalidPattern
			}
			stack++
		case internal.OpConcatN:
			if op.operand <= 0 {
				glog.V(2).Infof("negative concat size: %d", op.operand)
				return Pattern{}, ErrInvalidPattern
			}
			stack -= op.operand
			if stack < 0 {
				glog.V(2).Info("stack underflow")
				return Pattern{}, ErrInvalidPattern
			}
			stack++
		case internal.OpCapture:
			if op.operand < 0 || len(pool) <= op.operand {
				glog.V(2).Infof("variable name index out of bound: %d", op.operand)
				return Pattern{}, ErrInvalidPattern
			}
			v := pool[op.operand]
			op.operand = len(vars)
			vars = append(vars, v)
			stack--
			if stack < 0 {
				glog.V(2).Info("stack underflow")
				return Pattern{}, ErrInvalidPattern
			}
		default:
			glog.V(2).Infof("invalid opcode: %d", op.code)
			return Pattern{}, ErrInvalidPattern
		}

		if maxstack < stack {
			maxstack = stack
		}
		typedOps = append(typedOps, op)
	}
	glog.V(3).Info("pattern successfully built")
	return Pattern{
		ops:       typedOps,
		pool:      pool,
		vars:      vars,
		stacksize: maxstack,
		verb:      verb,
	}, nil
}

// MustPattern is a helper function which makes it easier to call NewPattern in variable initialization.
func MustPattern(p Pattern, err error) Pattern {
	if err != nil {
		glog.Fatalf("Pattern initialization failed: %v", err)
	}
	return p
}

// Match examines components if it matches to the Pattern.
// If it matches, the function returns a mapping from field paths to their captured values.
// If otherwise, the function returns an error.
func (p Pattern) Match(components []string, verb string) (map[string]string, error) {
	glog.V(2).Infof("matching (%q, %q) to %v", components, verb, p)

	if p.verb != verb {
		return nil, ErrNotMatch
	}

	var pos int
	stack := make([]string, 0, p.stacksize)
	captured := make([]string, len(p.vars))
	l := len(components)
	for _, op := range p.ops {
		switch op.code {
		case internal.OpNop:
			continue
		case internal.OpPush, internal.OpLitPush:
			if pos >= l {
				glog.V(1).Infof("insufficient # of segments")
				return nil, ErrNotMatch
			}
			c := components[pos]
			if op.code == internal.OpLitPush {
				if lit := p.pool[op.operand]; c != lit {
					glog.V(1).Infof("literal segment mismatch: got %q; want %q", c, lit)
					return nil, ErrNotMatch
				}
			}
			stack = append(stack, c)
			pos++
		case internal.OpPushM:
			stack = append(stack, strings.Join(components[pos:], "/"))
			pos = len(components)
		case internal.OpConcatN:
			n := op.operand
			l := len(stack) - n
			stack = append(stack[:l], strings.Join(stack[l:], "/"))
		case internal.OpCapture:
			n := len(stack) - 1
			captured[op.operand] = stack[n]
			stack = stack[:n]
		}
	}
	if pos < l {
		glog.V(1).Infof("remaining segments: %q", components[pos:])
		return nil, ErrNotMatch
	}
	bindings := make(map[string]string)
	for i, val := range captured {
		bindings[p.vars[i]] = val
	}
	return bindings, nil
}

// Verb returns the verb part of the Pattern.
func (p Pattern) Verb() string { return p.verb }

func (p Pattern) String() string {
	var stack []string
	for _, op := range p.ops {
		switch op.code {
		case internal.OpNop:
			continue
		case internal.OpPush:
			stack = append(stack, "*")
		case internal.OpLitPush:
			stack = append(stack, p.pool[op.operand])
		case internal.OpPushM:
			stack = append(stack, "**")
		case internal.OpConcatN:
			n := op.operand
			l := len(stack) - n
			stack = append(stack[:l], strings.Join(stack[l:], "/"))
		case internal.OpCapture:
			n := len(stack) - 1
			stack[n] = fmt.Sprintf("{%s=%s}", p.vars[op.operand], stack[n])
		}
	}
	segs := strings.Join(stack, "/")
	if p.verb != "" {
		return fmt.Sprintf("/%s:%s", segs, p.verb)
	}
	return "/" + segs
}
