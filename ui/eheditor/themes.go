// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eheditor

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