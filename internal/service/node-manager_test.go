package service_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func TestNodeManagerServiceCodeReturnsLastCodeCreated(t *testing.T) {
	s := service.NewNodeManagerService()

	// create new code
	c1 := s.Code()
	c2 := s.Code()

	if c1 != c2 {
		t.Errorf("expected %v, but got %v", c1, c2)
	}
}

func TestNodeManagerServiceValid(t *testing.T) {
	s := service.NewNodeManagerService()

	testCases := []struct {
		desc   string
		input  string
		output bool
	}{
		{
			desc:   "returns false if no code has been created yet",
			input:  "my.code",
			output: false,
		},
		{
			desc:   "returns true if code matches current code",
			input:  func() string { return s.Code().Code }(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if output := s.Valid(tC.input); output != tC.output {
				t.Errorf("expected %v but got %v", tC.output, output)
			}
		})
	}
}

func TestNodeManagerServiceIsAuthenticated(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)

	testCases := []struct {
		desc   string
		input  models.Node
		output bool
	}{
		{
			desc:   "returns false if node is not contained in the network",
			input:  models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}),
			output: false,
		},
		{
			desc:   "returns true if node is contained in the network",
			input:  n,
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if output := s.IsAuthenticated(tC.input); output != tC.output {
				t.Errorf("expected %v but got %v", tC.output, output)
			}
		})
	}
}

func TestNodeManagerServiceAuthenticate(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)
	code := s.Code()

	testCases := []struct {
		desc  string
		node  models.Node
		code  string
		error bool
	}{
		{
			desc:  "returns err if node is already contained in the network",
			node:  n,
			code:  models.UUID(),
			error: true,
		},
		{
			desc:  "returns err if code is invalid",
			code:  models.UUID(),
			node:  models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}),
			error: true,
		},

		{
			desc:  "does not return err if code is valid",
			node:  models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}),
			code:  code.Code,
			error: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := s.Authenticate(tC.node, tC.code)
			if error := output != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestNodeManagerServiceJoin(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)

	testCases := []struct {
		desc  string
		input models.Node
		error bool
	}{
		{
			desc:  "returns err if node is not authenticated",
			input: models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}),
			error: true,
		},
		{
			desc:  "does not return err if node is authenticated",
			input: n,
			error: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := s.Join(tC.input)
			if error := output != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestNodeManagerServiceUpdate(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)

	testCases := []struct {
		desc  string
		input models.Node
		error bool
	}{
		{
			desc:  "returns err if node is not authenticated",
			input: models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}),
			error: true,
		},
		{
			desc:  "does not return err if node is authenticated",
			input: n,
			error: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := s.Update(tC.input)
			if error := output != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestNodeManagerServiceNode(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)

	testCases := []struct {
		desc  string
		input string
		node  models.Node
		ok    bool
	}{
		{
			desc:  "returns empty node, false if network does not contain a node with input id",
			input: models.UUID(),
			node:  models.Node{},
			ok:    false,
		},
		{
			desc:  "returns node, true if network contains node with input id",
			input: n.ID,
			node:  n,
			ok:    true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			node, ok := s.Node(tC.input)
			if ok != tC.ok {
				t.Errorf("expected (%v, %v) but got (%v, %v)", tC.node, tC.ok, node, ok)
			}
		})
	}
}

func TestNodeManagerServiceStreamHandlesConcurrentAccesses(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)

	defer func() {
		if r := recover(); r != nil {
			t.Error("expected not to panic, but panicked")
		}
	}()
	for i := 0; i < 20; i++ {
		go func() {
			ch := s.Stream()
			s.Release(ch)
		}()
	}
}

func TestNodeManagerServiceAuthenticateHandlesConcurrentAccesses(t *testing.T) {
	n := models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{})
	s := service.NewNodeManagerService(n)
	c := s.Code()

	defer func() {
		if r := recover(); r != nil {
			t.Error("expected not to panic, but panicked")
		}
	}()
	for i := 0; i < 20; i++ {
		go func() {
			err := s.Authenticate(models.NewNodeWithoutStats(models.UUID(), models.MachineInfo{}), c.Code)
			if err != nil {
				t.Errorf("expected not to error, but got %v", err)
			}

		}()
	}
}
