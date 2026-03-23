package router

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type stub struct{ id string }

func (s stub) Init() tea.Cmd                       { return nil }
func (s stub) Update(tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s stub) View() string                        { return s.id }

func TestPush(t *testing.T) {
	r := New(stub{id: "root"})
	r, _ = r.Update(PushMsg{Screen: stub{id: "second"}})

	if got := r.View(); got != "second" {
		t.Errorf("View() = %q, want %q", got, "second")
	}
	if len(r.stack) != 2 {
		t.Errorf("stack len = %d, want 2", len(r.stack))
	}
}

func TestPop(t *testing.T) {
	r := New(stub{id: "root"})
	r, _ = r.Update(PushMsg{Screen: stub{id: "second"}})
	r, _ = r.Update(PopMsg{})

	if got := r.View(); got != "root" {
		t.Errorf("View() after pop = %q, want %q", got, "root")
	}
	if len(r.stack) != 1 {
		t.Errorf("stack len = %d, want 1", len(r.stack))
	}
}

func TestPopOnOne(t *testing.T) {
	r := New(stub{id: "root"})
	r, _ = r.Update(PopMsg{})

	if len(r.stack) != 1 {
		t.Errorf("stack len = %d, want 1 (no-op)", len(r.stack))
	}
}

func TestReplaceRoot(t *testing.T) {
	r := New(stub{id: "root"})
	r, _ = r.Update(PushMsg{Screen: stub{id: "second"}})
	r, _ = r.ReplaceRoot(stub{id: "new-root"})

	if got := r.View(); got != "new-root" {
		t.Errorf("View() = %q, want %q", got, "new-root")
	}
	if len(r.stack) != 1 {
		t.Errorf("stack len = %d, want 1", len(r.stack))
	}
}

func TestStackDepth(t *testing.T) {
	r := New(stub{id: "1"})
	r, _ = r.Update(PushMsg{Screen: stub{id: "2"}})
	r, _ = r.Update(PushMsg{Screen: stub{id: "3"}})
	r, _ = r.Update(PushMsg{Screen: stub{id: "4"}})

	if len(r.stack) != 4 {
		t.Errorf("stack len = %d, want 4", len(r.stack))
	}
	if got := r.View(); got != "4" {
		t.Errorf("View() = %q, want %q", got, "4")
	}
}
