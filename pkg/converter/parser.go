package converter

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/veresnikov/statemachines/pkg/logger"
	"github.com/veresnikov/statemachines/pkg/machine"
)

type parser struct {
	log logger.Logger
}

func (p *parser) ParseMealy(path string) (map[string]*machine.MealyState, []*machine.MealyState, error) {
	p.log.Info(fmt.Sprintf("start parsing %v", path))
	data, err := p.readFile(path)
	if err != nil {
		return nil, nil, err
	}
	p.log.Info("complete parsing")

	idxStates, states := p.getMealyStates(data)
	p.fillMealyTransitions(idxStates, states, data)
	return idxStates, states, nil
}

func (p *parser) readFile(path string) ([][]string, error) {
	inputFile, err := os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return nil, p.log.Error(err)
	}
	defer func() {
		_ = p.log.Error(inputFile.Close())
	}()
	csvreader := csv.NewReader(inputFile)
	csvreader.Comma = ';'
	data, err := csvreader.ReadAll()
	if err != nil {
		return nil, p.log.Error(err)
	}
	return data, nil
}

func (p *parser) getMealyStates(data [][]string) (map[string]*machine.MealyState, []*machine.MealyState) {
	p.log.Info("parsing mealy states...")
	idxStates := make(map[string]*machine.MealyState)
	states := make([]*machine.MealyState, 0)
	for _, value := range data[0] {
		if value == "" {
			continue
		}
		state := &machine.MealyState{
			Name:        value,
			Transitions: make(map[string]machine.MealyTransition),
		}
		idxStates[value] = state
		states = append(states, state)
	}
	p.log.Info(fmt.Sprintf("complete parsing: parsed %v states", len(states)))
	return idxStates, states
}

func (p *parser) fillMealyTransitions(idxState map[string]*machine.MealyState, states []*machine.MealyState, data [][]string) {
	p.log.Info("fill mealy states transitions...")
	for i := 1; i < len(data); i++ {
		input := data[i][0]
		for n := 1; n < len(data[i]); n++ {
			state := states[n-1]
			v := strings.Split(data[i][n], "/")
			state.Transitions[input] = machine.MealyTransition{
				Signal: v[1],
				State:  idxState[v[0]],
			}
		}
	}
	p.log.Info("complete fill")
}
