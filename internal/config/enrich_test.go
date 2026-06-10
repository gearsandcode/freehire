package config

import (
	"strings"
	"testing"
)

func TestLoadEnrich_missingRequiredFailsFast(t *testing.T) {
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_API_KEY", "")
	t.Setenv("LLM_MODEL", "")

	_, err := LoadEnrich()
	if err == nil {
		t.Fatal("expected error when LLM_* are unset, got nil")
	}
	for _, want := range []string{"LLM_BASE_URL", "LLM_API_KEY", "LLM_MODEL"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should name missing %s", err.Error(), want)
		}
	}
}

func TestLoadEnrich_namesOnlyTheMissingOne(t *testing.T) {
	t.Setenv("LLM_BASE_URL", "http://gateway:4000/v1")
	t.Setenv("LLM_API_KEY", "sk-test")
	t.Setenv("LLM_MODEL", "")

	_, err := LoadEnrich()
	if err == nil {
		t.Fatal("expected error when LLM_MODEL is unset")
	}
	if !strings.Contains(err.Error(), "LLM_MODEL") {
		t.Errorf("error %q should name LLM_MODEL", err.Error())
	}
	if strings.Contains(err.Error(), "LLM_BASE_URL") || strings.Contains(err.Error(), "LLM_API_KEY") {
		t.Errorf("error %q should not name the set vars", err.Error())
	}
}

func TestLoadEnrich_defaultsAndOverrides(t *testing.T) {
	t.Setenv("LLM_BASE_URL", "http://gateway:4000/v1")
	t.Setenv("LLM_API_KEY", "sk-test")
	t.Setenv("LLM_MODEL", "qwen2.5-72b")

	got, err := LoadEnrich()
	if err != nil {
		t.Fatalf("LoadEnrich: %v", err)
	}
	if got.LLMBaseURL != "http://gateway:4000/v1" || got.LLMModel != "qwen2.5-72b" {
		t.Errorf("unexpected config: %+v", got)
	}
	// Tunables fall back to conservative defaults.
	if got.BatchSize != 50 || got.LeaseSeconds != 300 || got.MaxAttempts != 3 {
		t.Errorf("defaults wrong: batch=%d lease=%d max=%d", got.BatchSize, got.LeaseSeconds, got.MaxAttempts)
	}

	t.Setenv("ENRICH_BATCH_SIZE", "10")
	got, err = LoadEnrich()
	if err != nil {
		t.Fatalf("LoadEnrich: %v", err)
	}
	if got.BatchSize != 10 {
		t.Errorf("batch override = %d, want 10", got.BatchSize)
	}
}
