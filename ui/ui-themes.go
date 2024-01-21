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

package ui

import (
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/ctk"
)

var (
	WindowTheme paint.Theme

	ActiveButtonTheme  paint.Theme
	DefaultButtonTheme paint.Theme

	SidebarFrameTheme  paint.Theme
	SidebarHeaderTheme paint.Theme
	SidebarButtonTheme paint.Theme
	SidebarActiveTheme paint.Theme
)

var ButtonActiveTheme paint.ThemeName = "toggle-button-active"

func init() {
	theme := paint.GetDefaultColorTheme()

	borders, _ := paint.GetDefaultBorderRunes(paint.RoundedBorder)
	arrows, _ := paint.GetArrows(paint.WideArrow)

	style := paint.GetDefaultColorStyle()
	style = style.Background(paint.ColorNavy)
	styleLight := style.Foreground(paint.ColorWhite)

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

	DefaultButtonTheme, _ = paint.GetTheme(ctk.ButtonColorTheme)

	styleNormal := style.Foreground(paint.ColorWhite).Background(paint.ColorDarkGreen)
	styleActive := style.Foreground(paint.ColorWhite).Background(paint.ColorForestGreen)
	styleInsensitive := style.Foreground(paint.ColorDarkSlateGray).Background(paint.ColorRosyBrown)
	ActiveButtonTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(false).Bold(true),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	}
	paint.RegisterTheme(ButtonActiveTheme, ActiveButtonTheme)

	styleNormal = style.Foreground(paint.ColorWhite).Background(paint.ColorNavy)
	styleActive = style.Foreground(paint.ColorWhite).Background(paint.ColorNavy)
	styleInsensitive = style.Foreground(paint.ColorDarkSlateGray).Background(paint.ColorNavy)
	SidebarButtonTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(true),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	}
	paint.RegisterTheme("sidebar-button-theme", SidebarButtonTheme)

	SidebarHeaderTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(false).Bold(true),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultNilRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(false).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultNilRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	}
	paint.RegisterTheme("sidebar-header-theme", SidebarHeaderTheme)

	styleNormal = style.Foreground(paint.ColorWhite).Background(paint.ColorDarkGreen)
	styleActive = style.Foreground(paint.ColorWhite).Background(paint.ColorDarkGreen)
	styleInsensitive = style.Foreground(paint.ColorDarkSlateGray).Background(paint.ColorDarkGreen)
	SidebarActiveTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(true),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	}
	paint.RegisterTheme("sidebar-active-theme", SidebarActiveTheme)

	SidebarFrameTheme = theme.Clone()

	styleNormal = style.Foreground(paint.ColorWhite).Background(paint.ColorNavy)
	styleActive = style.Foreground(paint.ColorWhite).Background(paint.ColorNavy)
	styleInsensitive = style.Foreground(paint.ColorDarkSlateGray).Background(paint.ColorNavy)
	SidebarFrameTheme.Content = paint.ThemeAspect{
		Normal:      styleNormal.Dim(true).Bold(false),
		Selected:    styleActive.Dim(true).Bold(false),
		Active:      styleActive.Dim(true).Bold(false).Reverse(true),
		Prelight:    styleActive.Dim(true),
		Insensitive: styleInsensitive.Dim(true),
		FillRune:    paint.DefaultFillRune,
		ArrowRunes:  arrows,
		Overlay:     false,
	}
	SidebarFrameTheme.Border = paint.ThemeAspect{
		Normal:      styleNormal.Dim(true).Bold(false),
		Selected:    styleActive.Dim(true).Bold(false),
		Active:      styleActive.Dim(true).Bold(false).Reverse(true),
		Prelight:    styleActive.Dim(true),
		Insensitive: styleInsensitive.Dim(true),
		FillRune:    paint.DefaultFillRune,
		BorderRunes: paint.BorderRuneSet{
			TopRight:    paint.DefaultNilRune, // paint.RuneURCornerRounded,
			Top:         paint.RuneHLine,
			TopLeft:     paint.DefaultNilRune, // paint.RuneULCornerRounded,
			Left:        paint.DefaultNilRune,
			BottomLeft:  paint.DefaultNilRune, // paint.RuneLLCornerRounded,
			Bottom:      paint.DefaultNilRune, // paint.RuneHLine,
			BottomRight: paint.DefaultNilRune, // paint.RuneLRCornerRounded,
			Right:       paint.DefaultNilRune,
		},
		ArrowRunes: arrows,
		Overlay:    false,
	}

	paint.RegisterTheme("sidebar-frame-theme", SidebarFrameTheme)

	entryColorTheme, _ := paint.GetTheme(ctk.EntryColorTheme)
	entryColorTheme.Border.Normal = entryColorTheme.Border.Normal.Foreground(paint.ColorYellow)
	entryColorTheme.Content.Normal = entryColorTheme.Content.Normal.Foreground(paint.ColorYellow)
	paint.RegisterTheme(ctk.EntryColorTheme, entryColorTheme)
}
