package provider

import (
	"context"
	"testing"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

type fakeProvider struct{ name string }

func (f fakeProvider) Name() string { return f.name }

func (f fakeProvider) Search(_ context.Context, _ *model.Channel, _ *dto.UnifiedSearchRequest) (*dto.UnifiedSearchResponse, error) {
	return &dto.UnifiedSearchResponse{Query: "ok"}, nil
}

func TestGetProvider_DefaultTavilyRegistered(t *testing.T) {
	_, ok := GetProvider("tavily")
	if !ok {
		t.Fatalf("expected default tavily provider registered")
	}
}

func TestRegisterProvider_DuplicatePanics(t *testing.T) {
	name := "fake-dup"
	RegisterProvider(fakeProvider{name: name})
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on duplicate register")
		}
	}()
	RegisterProvider(fakeProvider{name: name})
}

