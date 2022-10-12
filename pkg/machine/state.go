package machine

type MealyState struct {
	Name        string
	Transitions map[string]MealyTransition
}

type MealyTransition struct {
	Signal string
	State  *MealyState
}

type MooreState struct {
	Name        string
	Signal      string
	Transitions map[string]MooreTransition
}

type MooreTransition struct {
	State *MooreState
}
