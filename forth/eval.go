//go:build !solution

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNotEnoughArgs = errors.New("not enough arguments in stack")
)

type Stack struct {
	arr []int
}

func NewStack() *Stack {
	return &Stack{arr: make([]int, 0)}
}

func (s *Stack) len() int {
	return len(s.arr)
}

func (s *Stack) push(n int) {
	s.arr = append(s.arr, n)
}

func (s *Stack) pop() int {
	if len(s.arr) == 0 {
		panic(errors.New("stack is empty"))
	}
	n := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return n
}

func (s *Stack) peek() int {
	if len(s.arr) == 0 {
		panic(errors.New("stack is empty"))
	}
	return s.arr[len(s.arr)-1]
}

type Evaluator struct {
	stack         *Stack
	allowCommand  map[string]bool
	customCommand map[string][]string
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		NewStack(),
		map[string]bool{
			"+":    true,
			"-":    true,
			"*":    true,
			"/":    true,
			"dup":  true,
			"over": true,
			"drop": true,
			"swap": true,
		},
		make(map[string][]string)}
}

func (e *Evaluator) plus() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	n := e.stack.pop() + e.stack.pop()
	e.stack.push(n)
	return nil
}

func (e *Evaluator) minus() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	n := -e.stack.pop() - e.stack.pop()
	e.stack.push(n)
	return nil
}

func (e *Evaluator) multiply() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	n := e.stack.pop() * e.stack.pop()
	e.stack.push(n)
	return nil
}

func (e *Evaluator) divide() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	n1 := e.stack.pop()
	if n1 == 0 {
		return ErrNotEnoughArgs
	}
	n2 := e.stack.pop()
	n := n2 / n1
	e.stack.push(n)
	return nil
}

func (e *Evaluator) dup() error {
	if e.stack.len() < 1 {
		return ErrNotEnoughArgs
	}
	e.stack.push(e.stack.peek())
	return nil
}

func (e *Evaluator) over() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	num1 := e.stack.pop()
	num2 := e.stack.peek()
	e.stack.push(num1)
	e.stack.push(num2)
	return nil
}

func (e *Evaluator) drop() error {
	if e.stack.len() < 1 {
		return ErrNotEnoughArgs
	}
	e.stack.pop()
	return nil
}

func (e *Evaluator) swap() error {
	if e.stack.len() < 2 {
		return ErrNotEnoughArgs
	}
	num1 := e.stack.pop()
	num2 := e.stack.pop()
	e.stack.push(num1)
	e.stack.push(num2)
	return nil
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {

	reserveStack := NewStack()
	for _, element := range e.stack.arr {
		reserveStack.push(element)
	}
	reserveCommand := make(map[string][]string)
	for k, v := range e.customCommand {
		reserveCommand[k] = v
	}
	row = strings.ToLower(row)
	args := strings.Split(row, " ")
	define := make([]string, 0)
	i := 0
	for ; i < len(args); i++ {
		if args[i] != ":" && !e.allowCommand[args[i]] {
			if _, ok := e.customCommand[args[i]]; !ok {
				n, err := strconv.Atoi(args[i])
				if err != nil {
					e.stack = reserveStack
					e.customCommand = reserveCommand
					return nil, errors.New(fmt.Sprint(args[i], ": invalid argument"))
				}
				e.stack.push(n)
			}
		}

		if args[i] == ":" {
			i++
			for i < len(args) && args[i] != ";" {
				define = append(define, args[i])
				i++
			}
			if i == len(args) {
				e.stack = reserveStack
				e.customCommand = reserveCommand
				return nil, errors.New(fmt.Sprint(args[i], ": invalid argument"))
			}
			if _, ok := strconv.Atoi(define[0]); ok == nil {
				e.stack = reserveStack
				e.customCommand = reserveCommand
				return nil, errors.New(fmt.Sprint(define[0], ": invalid argument"))
			}
			temp := make([]string, 0)
			for _, arg := range define[1:] {
				if _, ok := e.customCommand[arg]; ok {
					temp = append(temp, e.customCommand[arg]...)
				} else if e.allowCommand[arg] {
					temp = append(temp, arg)
				} else {
					if _, err := strconv.Atoi(arg); err != nil {
						e.stack = reserveStack
						e.customCommand = reserveCommand
						return nil, errors.New(fmt.Sprint(arg, ": invalid argument"))
					} else {
						temp = append(temp, arg)
					}
				}
			}
			e.customCommand[define[0]] = make([]string, len(temp))
			copy(e.customCommand[define[0]], temp)
		}

		if _, ok := e.customCommand[args[i]]; ok || e.allowCommand[args[i]] {
			temp := make([]string, 0)
			if _, ok := e.customCommand[args[i]]; ok {
				temp = e.customCommand[args[i]]
			} else {
				temp = append(temp, args[i])
			}
			var errArgument error
			for _, arg := range temp {
				errArgument = nil
				switch arg {
				case "*":
					errArgument = e.multiply()
				case "/":
					errArgument = e.divide()
				case "+":
					errArgument = e.plus()
				case "-":
					errArgument = e.minus()
				case "dup":
					errArgument = e.dup()
				case "over":
					errArgument = e.over()
				case "drop":
					errArgument = e.drop()
				case "swap":
					errArgument = e.swap()
				default:
					n, _ := strconv.Atoi(arg)
					e.stack.push(n)
				}
				if errArgument != nil {
					e.stack = reserveStack
					e.customCommand = reserveCommand
					return nil, errArgument
				}
			}

		}
	}
	return e.stack.arr, nil
}
