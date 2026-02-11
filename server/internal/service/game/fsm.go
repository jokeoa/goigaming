package game

import (
	"github.com/jokeoa/goigaming/internal/core/domain"
)

var validTransitions = map[domain.GameStage][]domain.GameStage{
	domain.StageWaiting:  {domain.StagePreflop},
	domain.StagePreflop:  {domain.StageFlop, domain.StageShowdown, domain.StageComplete},
	domain.StageFlop:     {domain.StageTurn, domain.StageShowdown, domain.StageComplete},
	domain.StageTurn:     {domain.StageRiver, domain.StageShowdown, domain.StageComplete},
	domain.StageRiver:    {domain.StageShowdown, domain.StageComplete},
	domain.StageShowdown: {domain.StageComplete},
}

type GameFSM struct {
	stage domain.GameStage
}

func NewGameFSM() GameFSM {
	return GameFSM{stage: domain.StageWaiting}
}

func NewGameFSMFromStage(stage domain.GameStage) GameFSM {
	return GameFSM{stage: stage}
}

func (f GameFSM) Transition(to domain.GameStage) (GameFSM, error) {
	allowed, ok := validTransitions[f.stage]
	if !ok {
		return f, domain.ErrInvalidTransition
	}

	for _, s := range allowed {
		if s == to {
			return GameFSM{stage: to}, nil
		}
	}

	return f, domain.ErrInvalidTransition
}

func (f GameFSM) Stage() domain.GameStage {
	return f.stage
}

func (f GameFSM) IsTerminal() bool {
	return f.stage == domain.StageComplete
}

func (f GameFSM) NextStage() domain.GameStage {
	switch f.stage {
	case domain.StageWaiting:
		return domain.StagePreflop
	case domain.StagePreflop:
		return domain.StageFlop
	case domain.StageFlop:
		return domain.StageTurn
	case domain.StageTurn:
		return domain.StageRiver
	case domain.StageRiver:
		return domain.StageShowdown
	case domain.StageShowdown:
		return domain.StageComplete
	default:
		return domain.StageComplete
	}
}
