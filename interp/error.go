package interp

import (
	"fmt"
	"go/token"
	"strings"
)

type ErrorPosition interface {
	Pos() token.Position
	End() token.Position
	Reason() string
}

// A cfgError represents an error during CFG build stage.
type cfgError struct {
	pos    token.Position
	end    token.Position
	reason string
}

func (c *cfgError) Pos() token.Position {
	return c.pos
}

func (c *cfgError) End() token.Position {
	return c.end
}

func (c *cfgError) Reason() string {
	return c.reason
}

func (c *cfgError) Error() string {
	posString := c.pos.String()
	if c.pos.Filename == DefaultSourceName {
		posString = strings.TrimPrefix(posString, DefaultSourceName+":")
	}
	return fmt.Sprintf("%s %s", posString, c.reason)
}

func (n *node) cfgErrorf(format string, a ...interface{}) *cfgError {
	pos := n.interp.fset.Position(n.pos)
	end := n.interp.fset.Position(n.end)
	return &cfgError{pos, end, fmt.Sprintf(format, a...)}
}

// A runError represents an error during runtime.
type runError struct {
	pos    token.Position
	end    token.Position
	reason string
}

func (c *runError) Pos() token.Position {
	return c.pos
}

func (c *runError) End() token.Position {
	return c.end
}

func (c *runError) Reason() string {
	return c.reason
}

func (c *runError) Error() string {
	posString := c.pos.String()
	if c.pos.Filename == DefaultSourceName {
		posString = strings.TrimPrefix(posString, DefaultSourceName+":")
	}
	return fmt.Sprintf("%s %s", posString, c.reason)
}

func (n *node) runErrorf(format string, a ...interface{}) *runError {
	pos := n.interp.fset.Position(n.pos)
	end := n.interp.fset.Position(n.end)
	return &runError{pos, end, fmt.Sprintf(format, a...)}
}
