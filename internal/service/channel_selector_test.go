package service

import (
	"testing"
	"uni-search-hub/internal/model"
)

func TestRoundRobinSelector_Select(t *testing.T) {
	s := &RoundRobinSelector{}
	channels := []*model.Channel{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	got1, err := s.Select("tavily", channels)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got2, _ := s.Select("tavily", channels)
	got3, _ := s.Select("tavily", channels)
	got4, _ := s.Select("tavily", channels)

	if got1.Id != 1 || got2.Id != 2 || got3.Id != 3 || got4.Id != 1 {
		t.Fatalf("round robin mismatch: %d %d %d %d", got1.Id, got2.Id, got3.Id, got4.Id)
	}
}

func TestRoundRobinSelector_Select_Empty(t *testing.T) {
	s := &RoundRobinSelector{}
	_, err := s.Select("tavily", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

