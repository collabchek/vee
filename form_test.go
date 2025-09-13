package vee

import (
	"strings"
	"testing"
)

func TestFormRendering(t *testing.T) {
	type Demo struct {
		Username string
		Password string `vee:"type:password"`
	}
	noActionMethod, err := Render(&Demo{}, FormActionScriptOption())
	if err != nil {
		t.Errorf("Expected nil error, got %q", err)
	}
	methodIdx := strings.Index(noActionMethod, "method=\"POST\"")
	if methodIdx != -1 {
		t.Errorf("Expected no method attribute, got '%s'", noActionMethod[methodIdx:])
	}
	actionIdx := strings.Index(noActionMethod, "action=\"")
	if actionIdx != -1 {
		t.Errorf("Expected no action attribute, got '%s'", noActionMethod[actionIdx:])
	}

	actionMethod, err := Render(&Demo{}, FormActionOption("/register-user"))
	if err != nil {
		t.Errorf("Expected nil error, got %q", err)
	}
	methodIdx = strings.Index(actionMethod, "method=\"POST\"")
	if methodIdx == -1 {
		t.Errorf("Expected method=\"POST\", found nothing")
	}
	actionIdx = strings.Index(actionMethod, "action=\"/register-user\"")
	if actionIdx == -1 {
		t.Errorf("Expected action=\"/register-user\", found nothing")
	}
}
