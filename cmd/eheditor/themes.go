package main

import (
	"github.com/go-curses/cdk/lib/paint"
)

var (
	ViewerTheme paint.Theme
	WindowTheme paint.Theme
)

func init() {
	theme := paint.GetDefaultColorTheme()

	ViewerTheme = theme.Clone()

	borders, _ := paint.GetDefaultBorderRunes(paint.RoundedBorder)
	arrows, _ := paint.GetArrows(paint.WideArrow)

	style := paint.GetDefaultColorStyle()
	style = style.Background(paint.ColorNavy)
	styleLight := style.Foreground(paint.ColorWhite)
	styleDark := style.Foreground(paint.ColorDarkGray)

	ViewerTheme.Content.Normal = styleDark.Dim(true)
	ViewerTheme.Content.Selected = styleLight.Dim(false)
	ViewerTheme.Content.Active = styleLight.Dim(false)
	ViewerTheme.Content.Prelight = styleLight.Dim(false)
	ViewerTheme.Content.Insensitive = styleDark.Dim(true)
	ViewerTheme.Content.FillRune = paint.DefaultFillRune
	ViewerTheme.Content.BorderRunes = borders
	ViewerTheme.Content.ArrowRunes = arrows
	ViewerTheme.Content.Overlay = false

	ViewerTheme.Border.Normal = styleDark.Dim(true)
	ViewerTheme.Border.Selected = styleLight.Dim(false)
	ViewerTheme.Border.Active = styleLight.Dim(false)
	ViewerTheme.Border.Prelight = styleLight.Dim(false)
	ViewerTheme.Border.Insensitive = styleDark.Dim(true)
	ViewerTheme.Border.FillRune = paint.DefaultFillRune
	ViewerTheme.Border.BorderRunes = borders
	ViewerTheme.Border.ArrowRunes = arrows
	ViewerTheme.Border.Overlay = false

	paint.RegisterTheme(paint.ColorTheme, ViewerTheme)

	WindowTheme = theme.Clone()
	WindowTheme.Content.Normal = styleLight.Dim(false)
	WindowTheme.Content.Selected = WindowTheme.Content.Normal
	WindowTheme.Content.Active = WindowTheme.Content.Normal
	WindowTheme.Content.Prelight = WindowTheme.Content.Normal
	WindowTheme.Content.Insensitive = WindowTheme.Content.Normal
	WindowTheme.Content.BorderRunes = borders
	WindowTheme.Content.ArrowRunes = arrows
	WindowTheme.Content.Overlay = false
	WindowTheme.Border.Normal = WindowTheme.Content.Normal
	WindowTheme.Border.Active = WindowTheme.Content.Normal
	WindowTheme.Border.Prelight = WindowTheme.Content.Normal
	WindowTheme.Border.Insensitive = WindowTheme.Content.Normal
	// WindowTheme.Border.BorderRunes = borders
	WindowTheme.Border.ArrowRunes = arrows
	WindowTheme.Border.Overlay = false
	paint.RegisterTheme(paint.DisplayTheme, WindowTheme)
}