package main

import (
	"github.com/maxzender/jsonexplorer/jsonfmt"
	"github.com/maxzender/jsonexplorer/treemodel"
	"github.com/nsf/termbox-go"
)

type colorWriter struct {
	Lines    []treemodel.Line
	colorMap map[jsonfmt.TokenType]termbox.Attribute
	line     int
	bgColor  termbox.Attribute
}

func NewColorWriter(colorMap map[jsonfmt.TokenType]termbox.Attribute, bgColor termbox.Attribute) *colorWriter {
	writer := &colorWriter{
		colorMap: colorMap,
		bgColor:  bgColor,
	}

	writer.Lines = append(writer.Lines, treemodel.Line{})

	return writer
}

func (w *colorWriter) Write(s string, t jsonfmt.TokenType) {
	for _, c := range s {
		w.Lines[w.line] = append(w.Lines[w.line], treemodel.Char{c, w.colorMap[t]})
	}
}

func (w *colorWriter) Newline() {
	w.Lines = append(w.Lines, treemodel.Line{})
	w.line++
}
