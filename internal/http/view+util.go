package http

import (
	"github.com/labstack/echo/v4"
	"github.com/mssola/useragent"
)

const (
	MobileBreakpoint  Breakpoint = 560
	TabletBreakpoint  Breakpoint = 860
	DesktopBreakpoint Breakpoint = 1440
)

type Breakpoint int
type UserAgent useragent.UserAgent
type ContextView struct {
	echo.Context
}

func NewContextView(
	ctx echo.Context,
) ContextView {
	return ContextView{
		Context: ctx,
	}
}

func (v ContextView) Breakpoint() Breakpoint {
	return extractUserAgent(v.Context).Breakpoint()
}

func (v ContextView) HasAdminRights() bool {
	return hasAdminRights(v.Context)
}

func (ua UserAgent) Breakpoint() Breakpoint {
	if u := useragent.UserAgent(ua); u.Mobile() {
		return MobileBreakpoint
	}

	return DesktopBreakpoint
}

func (bp Breakpoint) Mobile() bool {
	return bp == MobileBreakpoint
}

func (bp Breakpoint) Tablet() bool {
	return bp == TabletBreakpoint
}

func (bp Breakpoint) Desktop() bool {
	return bp == DesktopBreakpoint
}

func (bp Breakpoint) ChartSize() (int, int) {
	switch bp {
	case MobileBreakpoint:
		return 300, 400
	case TabletBreakpoint:
		return 500, 400
	default:
		return 1200, 400
	}
}
