package converter

import (
	"context"
	stderr "errors"
	"fmt"

	"github.com/veresnikov/statemachines/pkg/logger"
	"github.com/veresnikov/statemachines/pkg/machine"
)

type Converter interface {
	MealyToMoore(ctx context.Context, input string, output string) error
	MooreToMealy(ctx context.Context, input string, output string) error
}

func NewConverter(log logger.Logger) Converter {
	return &converter{
		log:    log,
		parser: parser{log: log},
		writer: writer{log: log},
	}
}

type converter struct {
	log    logger.Logger
	parser parser
	writer writer
}

func (c *converter) MealyToMoore(_ context.Context, input, output string) error {
	idxMealyStates, mealyStates, err := c.parser.ParseMealy(input)
	if err != nil {
		return err
	}

	idxMooreStates, _ := c.generateMooreStatesFormMealyStates(mealyStates)
	c.fillMooreTransitions(idxMealyStates, idxMooreStates)
	return c.writer.WriteMooreStatemachine(output, idxMooreStates)
}

func (c *converter) generateMooreStatesFormMealyStates(mealyStates []*machine.MealyState) (map[string]*machine.MooreState, []*machine.MooreState) {
	c.log.Info("generating moore states...")
	idxMooreStates := make(map[string]*machine.MooreState)
	mooreStates := make([]*machine.MooreState, 0)

	for _, mealyState := range mealyStates {
		for _, transition := range mealyState.Transitions {
			name := transition.State.Name + "/" + transition.Signal
			if _, ok := idxMooreStates[name]; !ok {
				mooreState := &machine.MooreState{
					Name:        name,
					Signal:      transition.Signal,
					Transitions: make(map[string]machine.MooreTransition),
				}
				idxMooreStates[name] = mooreState
				mooreStates = append(mooreStates, mooreState)
			}
		}
	}
	c.log.Info(fmt.Sprintf("complete generating: generated %v states", len(mooreStates)))
	return idxMooreStates, mooreStates
}

func (c *converter) fillMooreTransitions(
	idxMealyStates map[string]*machine.MealyState,
	idxMooreStates map[string]*machine.MooreState,
) {
	c.log.Info("fill moore transitions...")
	countTransitions := 0
	filledTransitions := make(map[string][]string)
	for _, mealyState := range idxMealyStates {
		for _, transition := range mealyState.Transitions {
			for input, _ := range mealyState.Transitions {
				currentStateName := transition.State.Name + "/" + transition.Signal
				currentMooreState := idxMooreStates[currentStateName]

				nextMealyState := transition.State
				nextStateName := nextMealyState.Transitions[input].State.Name + "/" + nextMealyState.Transitions[input].Signal
				nextMooreState := idxMooreStates[nextStateName]

				if isTransitionExist(filledTransitions, currentStateName, nextStateName) {
					continue
				}

				currentMooreState.Transitions[input] = machine.MooreTransition{
					State: nextMooreState,
				}
				countTransitions++
				filledTransitions[currentStateName] = append(filledTransitions[currentStateName], nextStateName)
			}
		}
	}
	c.log.Info(fmt.Sprintf("complete fill: filled %v transitions", countTransitions))
}

func (c *converter) MooreToMealy(_ context.Context, input, output string) error {
	return c.log.Error(stderr.New("not implemented"))
}

func isTransitionExist(filledTransitions map[string][]string, src, dst string) bool {
	for _, t := range filledTransitions[src] {
		if dst == t {
			return true
		}
	}
	return false
}
