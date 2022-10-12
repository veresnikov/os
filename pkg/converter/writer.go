package converter

import (
	"encoding/csv"
	"fmt"
	"github.com/veresnikov/statemachines/pkg/logger"
	"github.com/veresnikov/statemachines/pkg/machine"
	"os"
)

type writer struct {
	log logger.Logger
}

func (w *writer) WriteMooreStatemachine(output string, mooreMachine map[string]*machine.MooreState) error {
	outputFile, err := w.createFile(output)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.log.Error(outputFile.Close())
	}()
	csvwriter := csv.NewWriter(outputFile)
	csvwriter.Comma = ';'
	data := w.convertMooreMachineToCsv(mooreMachine)
	err = csvwriter.WriteAll(data)
	if err != nil {
		return w.log.Error(err)
	}
	return nil
}

func (w *writer) createFile(path string) (*os.File, error) {
	w.log.Info(fmt.Sprintf("create file %v", path))
	output, err := os.Create(path)
	if err != nil {
		return nil, w.log.Error(err)
	}
	return output, nil
}

func (w *writer) convertMooreMachineToCsv(idxMooreMachine map[string]*machine.MooreState) [][]string {
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
	for _, transition := range finalTransitions {
		data = append(data, transition)
	}
	return data
}
