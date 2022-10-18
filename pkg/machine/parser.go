package machine

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/veresnikov/statemachines/pkg/logger"
)

type Parser struct {
	Log logger.Logger
}

func (p *Parser) ParseMealy(path string) (map[string]*MealyState, []*MealyState, error) {
	p.Log.Info(fmt.Sprintf("start parsing %v", path))
	data, err := p.readFile(path)
	if err != nil {
		return nil, nil, err
	}
	p.Log.Info("complete parsing")

	idxStates, states := p.getMealyStates(data)
	p.fillMealyTransitions(idxStates, states, data)
	return idxStates, states, nil
}

func (p *Parser) readFile(path string) ([][]string, error) {
	inputFile, err := os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return nil, p.Log.Error(err)
	}
	defer func() {
		_ = p.Log.Error(inputFile.Close())
	}()
	csvreader := csv.NewReader(inputFile)
	csvreader.Comma = ';'
	data, err := csvreader.ReadAll()
	if err != nil {
		return nil, p.Log.Error(err)
	}
	return data, nil
}

func (p *Parser) getMealyStates(data [][]string) (map[string]*MealyState, []*MealyState) {
	p.Log.Info("parsing mealy states...")
	idxStates := make(map[string]*MealyState)
	states := make([]*MealyState, 0)
	for _, value := range data[0] {
		if value == "" {
			continue
		}
		state := &MealyState{
			Name:        value,
			Transitions: make(map[string]MealyTransition),
		}
		idxStates[value] = state
		states = append(states, state)
	}
	p.Log.Info(fmt.Sprintf("complete parsing: parsed %v states", len(states)))
	return idxStates, states
}

func (p *Parser) fillMealyTransitions(idxState map[string]*MealyState, states []*MealyState, data [][]string) {
	p.Log.Info("fill mealy states transitions...")
	for i := 1; i < len(data); i++ {
		input := data[i][0]
		for n := 1; n < len(data[i]); n++ {
			state := states[n-1]
			v := strings.Split(data[i][n], "/")
			state.Transitions[input] = MealyTransition{
				Signal: v[1],
				State:  idxState[v[0]],
			}
		}
	}
	p.Log.Info("complete fill")
}

func (p *Parser) ParseMoore(path string) (map[string]*MooreState, []*MooreState, error) {
	p.Log.Info(fmt.Sprintf("start parsing %v", path))
	data, err := p.readFile(path)
	if err != nil {
		return nil, nil, err
	}
	p.Log.Info("complete parsing")

	idxStates, states := p.getMooreStates(data)
	p.fillMooreTransitions(idxStates, states, data)
	return idxStates, states, nil
}

func (p *Parser) getMooreStates(data [][]string) (map[string]*MooreState, []*MooreState) {
	p.Log.Info("parsing moore states...")
	idxStates := make(map[string]*MooreState)
	states := make([]*MooreState, 0)
	for i := 0; i < len(data[0]); i++ {
		signal := data[0][i]
		stateName := data[1][i]
		if signal == "" {
			continue
		}
		state := &MooreState{
			Name:        stateName,
			Signal:      signal,
			Transitions: make(map[string]MooreTransition),
		}
		idxStates[stateName] = state
		states = append(states, state)
	}
	p.Log.Info(fmt.Sprintf("complete parsing: parsed %v states", len(states)))
	return idxStates, states
}

func (p *Parser) fillMooreTransitions(idxState map[string]*MooreState, states []*MooreState, data [][]string) {
	p.Log.Info("fill moore states transitions...")
	for i := 2; i < len(data); i++ {
		input := data[i][0]
		for n := 1; n < len(data[i]); n++ {
			state := states[n-1]
			state.Transitions[input] = MooreTransition{
				State: idxState[data[i][n]],
			}
		}
	}
	p.Log.Info("complete fill")
}
