package main

import (
    "strconv"
)


type ColorBrush struct {
    enabled bool
}


func colorize(enabled bool, s string, id int, bold bool) string {
    if enabled {
        var boldInt = "0"
        if bold {
            boldInt = "1"
        }

        var prefix = "\033[0;" + boldInt + ";3" + strconv.Itoa(id) + "m"
        var suffix = "\033[0;0m"
        return prefix + s + suffix
    }
    return s
}


func (b *ColorBrush) black(s string) string {
    return colorize(b.enabled, s, 0, false)
}

func (b *ColorBrush) red(s string) string {
    return colorize(b.enabled, s, 1, false)
}

func (b *ColorBrush) green(s string) string {
    return colorize(b.enabled, s, 2, false)
}

func (b *ColorBrush) yellow(s string) string {
    return colorize(b.enabled, s, 3, false)
}

func (b *ColorBrush) blue(s string) string {
    return colorize(b.enabled, s, 4, false)
}

func (b *ColorBrush) magenta(s string) string {
    return colorize(b.enabled, s, 5, false)
}

func (b *ColorBrush) cyan(s string) string {
    return colorize(b.enabled, s, 6, false)
}

func (b *ColorBrush) white(s string) string {
    return colorize(b.enabled, s, 7, false)
}


func (b *ColorBrush) bblack(s string) string {
    return colorize(b.enabled, s, 0, true)
}

func (b *ColorBrush) bred(s string) string {
    return colorize(b.enabled, s, 1, true)
}

func (b *ColorBrush) bgreen(s string) string {
    return colorize(b.enabled, s, 2, true)
}

func (b *ColorBrush) byellow(s string) string {
    return colorize(b.enabled, s, 3, true)
}

func (b *ColorBrush) bblue(s string) string {
    return colorize(b.enabled, s, 4, true)
}

func (b *ColorBrush) bmagenta(s string) string {
    return colorize(b.enabled, s, 5, true)
}

func (b *ColorBrush) bcyan(s string) string {
    return colorize(b.enabled, s, 6, true)
}

func (b *ColorBrush) bwhite(s string) string {
    return colorize(b.enabled, s, 7, true)
}
