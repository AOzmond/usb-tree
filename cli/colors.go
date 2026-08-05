package cli

import "charm.land/lipgloss/v2"

const (
	gray      = "#888888"
	white     = "#ffffff"
	black     = "#000000"
	hotPink   = "#ff028d"
	orange    = "#FF5c00"
	red       = "#FF0000"
	green     = "#00FF00"
	cyan      = "#003a3a"
	skyBlue   = "#00BFFF"
	gold      = "#FFD700"
	coralRed  = "#FF6B6B"
	paleGreen = "#98FB98"
	plum      = "#DDA0DD"
)

var (
	backgroundColor           = lipgloss.Color(cyan)
	inactiveNodeBorderColor   = lipgloss.Color(gray)
	lineHighlightColor        = lipgloss.Color(white)
	lineHighlightTextColor    = lipgloss.Color(black)
	tooltipTextColor          = lipgloss.Color(white)
	activeNodeBorderColor     = lipgloss.Color(hotPink)
	edgeHighlightChangeColor  = lipgloss.Color(orange)
	childChangeHighlightColor = lipgloss.Color(orange)
	removedStateColor         = lipgloss.Color(red)
	addedStateColor           = lipgloss.Color(green)
	linkTextColor             = lipgloss.Color(skyBlue)
	busTextColor              = lipgloss.Color(plum)
	deviceTextColor           = lipgloss.Color(gold)
	vidTextColor              = lipgloss.Color(coralRed)
	pidTextColor              = lipgloss.Color(paleGreen)
	nameTextColor             = lipgloss.Color(gray)
)
