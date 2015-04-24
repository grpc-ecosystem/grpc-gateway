package runtime

import (
	"errors"
	"strings"

	"github.com/golang/glog"
)

var (
	// ErrNotMatch indicates that the given HTTP request path does not match to the pattern.
	ErrNotMatch = errors.New("not match to the path pattern")
	// ErrInvalidPattern indicates that the given definition of Pattern is not valid.
	ErrInvalidPattern = errors.New("invalid pattern")
)

type opcode int

// These constants are the valid values of opcode.
const (
	// opNop does nothing
	opNop = opcode(iota)
	// opPush pushes a component to stack
	opPush
	// opLitPush pushes a component to stack if it matches to the literal
	opLitPush
	// opPushM concatenates the remaining components and pushes it to stack
	opPushM
	// opPopN pops a N items from stack, concatenates them and pushes it to stack
	opConcatN
	// opCapture pops an item and binds it to the variable
	opCapture
	// opEnd is the least postive invalid opcode.
	opEnd
)

type op struct {
	code    opcode
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
		op := op{code: opcode(ops[i]), operand: ops[i+1]}
		switch op.code {
		case opNop:
			continue
		case opPush, opPushM:
			stack++
		case opLitPush:
			if op.operand < 0 || len(pool) <= op.operand {
				glog.V(2).Infof("negative literal index: %d", op.operand)
				return Pattern{}, ErrInvalidPattern
			}
			stack++
		case opConcatN:
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
		case opCapture:
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
		case opNop:
			continue
		case opPush, opLitPush:
			if pos >= l {
				glog.V(1).Infof("insufficient # of segments")
				return nil, ErrNotMatch
			}
			c := components[pos]
			if op.code == opLitPush {
				if lit := p.pool[op.operand]; c != lit {
					glog.V(1).Infof("literal segment mismatch: got %q; want %q", c, lit)
					return nil, ErrNotMatch
				}
			}
			stack = append(stack, c)
			pos++
		case opPushM:
			stack = append(stack, strings.Join(components[pos:], "/"))
			pos = len(components)
		case opConcatN:
			n := op.operand
			l := len(stack) - n
			stack = append(stack[:l], strings.Join(stack[l:], "/"))
		case opCapture:
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
