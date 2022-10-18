package executor

import (
	stderr "errors"

	"github.com/veresnikov/statemachines/pkg/logger"
	"github.com/veresnikov/statemachines/pkg/machine"
)

type Executor interface {
	Run(state interface{}, input []string) ([]string, error)
}

func NewExecutor(
	log logger.Logger,
	useWarnings bool,
) Executor {
	return &executor{
		log:         log,
		useWarnings: useWarnings,
	}
}

type executor struct {
	log         logger.Logger
	useWarnings bool
}

func (e executor) Run(state interface{}, input []string) ([]string, error) {
	move, err := getMoveStrategy(state)
	if err != nil {
		return nil, e.log.Error(err)
	}
	result := make([]string, 0)
	currentState := state
	for _, v := range input {
		output, nextState, moveErr := move(v, currentState)
		if moveErr != nil {
			if e.useWarnings {
				e.log.Warn(moveErr)
				continue
			}
			return nil, e.log.Error(moveErr)
		}
		result = append(result, output)
		currentState = nextState
	}
	return result, nil
}

type moveFunc func(input string, state interface{}) (output string, nextState interface{}, err error)

func getMoveStrategy(state interface{}) (moveFunc, error) {
	switch state.(type) {
	case *machine.MealyState:
		return func(input string, state interface{}) (output string, nextState interface{}, err error) {
			s := state.(*machine.MealyState)
			transition, ok := s.Transitions[input]
			if !ok {
				return "", nil, stderr.New("unexpected input symbol")
			}
			return transition.Signal, transition.State, nil
		}, nil
	case *machine.MooreState:
		return func(input string, state interface{}) (output string, nextState interface{}, err error) {
			s := state.(*machine.MooreState)
			transition, ok := s.Transitions[input]
			if !ok {
				return "", nil, stderr.New("unexpected input symbol")
			}
			return transition.State.Signal, transition.State, nil
		}, nil
	default:
		return nil, stderr.New("undefined state type")
	}
}
