package machine

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/veresnikov/statemachines/pkg/logger"
)

type Writer struct {
	Log logger.Logger
}

func (w *Writer) WriteMooreStatemachine(output string, mooreMachine map[string]*MooreState) error {
	outputFile, err := w.createFile(output)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Log.Error(outputFile.Close())
	}()
	csvwriter := csv.NewWriter(outputFile)
	csvwriter.Comma = ';'
	data := w.convertMooreMachineToCsv(mooreMachine)
	err = csvwriter.WriteAll(data)
	if err != nil {
		return w.Log.Error(err)
	}
	return nil
}

func (w *Writer) createFile(path string) (*os.File, error) {
	w.Log.Info(fmt.Sprintf("create file %v", path))
	output, err := os.Create(path)
	if err != nil {
		return nil, w.Log.Error(err)
	}
	return output, nil
}

func (w *Writer) convertMooreMachineToCsv(idxMooreMachine map[string]*MooreState) [][]string {
	data := make([][]string, 0)

	signals := []string{""}
	states := []string{""}
	inputs := make([]string, 0)
	for _, state := range idxMooreMachine {
		for input := range state.Transitions {
			inputs = append(inputs, input)
		}
		break
	}

	transitions := make([][]string, 0)
	for _, state := range idxMooreMachine {
		signals = append(signals, state.Signal)
		states = append(states, state.Name)
		currentTransitions := make([]string, 0)
		for _, transition := range state.Transitions {
			currentTransitions = append(currentTransitions, transition.State.Name)
		}
		transitions = append(transitions, currentTransitions)
	}
	finalTransitions := make([][]string, len(inputs))
	for i := 0; i < len(finalTransitions); i++ {
		finalTransitions[i] = append(finalTransitions[i], inputs[i])
		for _, transition := range transitions {
			finalTransitions[i] = append(finalTransitions[i], transition[i])
		}
	}

	data = append(data, signals, states)
	data = append(data, finalTransitions...)
	return data
}

func (w *Writer) WriteMealyStatemachine(output string, mealyMachine map[string]*MealyState) error {
	outputFile, err := w.createFile(output)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Log.Error(outputFile.Close())
	}()
	csvwriter := csv.NewWriter(outputFile)
	csvwriter.Comma = ';'
	data := w.convertMealyMachineToCsv(mealyMachine)
	err = csvwriter.WriteAll(data)
	if err != nil {
		return w.Log.Error(err)
	}
	return nil
}

func (w *Writer) convertMealyMachineToCsv(idxMealyMachine map[string]*MealyState) [][]string {
	data := make([][]string, 0)

	states := []string{""}
	inputs := make([]string, 0)
	for _, state := range idxMealyMachine {
		for input := range state.Transitions {
			inputs = append(inputs, input)
		}
		break
	}

	transitions := make([][]string, 0)
	for _, state := range idxMealyMachine {
		states = append(states, state.Name)
		currentTransitions := make([]string, 0)
		for _, transition := range state.Transitions {
			currentTransitions = append(currentTransitions, transition.State.Name+"/"+transition.Signal)
		}
		transitions = append(transitions, currentTransitions)
	}
	finalTransitions := make([][]string, len(inputs))
	for i := 0; i < len(finalTransitions); i++ {
		finalTransitions[i] = append(finalTransitions[i], inputs[i])
		for _, transition := range transitions {
			finalTransitions[i] = append(finalTransitions[i], transition[i])
		}
	}

	data = append(data, states)
	data = append(data, finalTransitions...)
	return data
}
