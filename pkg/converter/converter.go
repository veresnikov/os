package converter

import (
	"context"
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
		parser: machine.Parser{Log: log},
		writer: machine.Writer{Log: log},
	}
}

type converter struct {
	log    logger.Logger
	parser machine.Parser
	writer machine.Writer
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
			for input := range mealyState.Transitions {
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
	idxMooreStates, mooreStates, err := c.parser.ParseMoore(input)
	if err != nil {
		return err
	}
	idxMealyStates, _ := c.generateMealyStatesFormMooreStates(mooreStates)
	c.fillMealyTransitions(idxMooreStates, idxMealyStates)
	return c.writer.WriteMealyStatemachine(output, idxMealyStates)
}

func (c *converter) generateMealyStatesFormMooreStates(
	mooreStates []*machine.MooreState,
) (map[string]*machine.MealyState, []*machine.MealyState) {
	c.log.Info("generating mealy states...")
	idxMealyStates := make(map[string]*machine.MealyState)
	mealyStates := make([]*machine.MealyState, 0)
	for _, mooreState := range mooreStates {
		mealyState := &machine.MealyState{
			Name:        mooreState.Name,
			Transitions: make(map[string]machine.MealyTransition),
		}
		idxMealyStates[mooreState.Name] = mealyState
		mealyStates = append(mealyStates, mealyState)
	}
	c.log.Info(fmt.Sprintf("complete generating: generated %v states", len(mealyStates)))
	return idxMealyStates, mealyStates
}

func (c *converter) fillMealyTransitions(
	idxMooreStates map[string]*machine.MooreState,
	idxMealyStates map[string]*machine.MealyState,
) {
	for _, mealyState := range idxMealyStates {
		mooreState := idxMooreStates[mealyState.Name]
		for input, transition := range mooreState.Transitions {
			nextMooreState := transition.State
			mealyState.Transitions[input] = machine.MealyTransition{
				Signal: nextMooreState.Signal,
				State:  idxMealyStates[nextMooreState.Name],
			}
		}
	}
}

func isTransitionExist(filledTransitions map[string][]string, src, dst string) bool {
	for _, t := range filledTransitions[src] {
		if dst == t {
			return true
		}
	}
	return false
}
