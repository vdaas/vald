//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// This file is a workaround for placing panels from top left.
// This file can be removed when https://github.com/grafana/grafana-foundation-sdk/issues/673 is resolved.

package main

import "github.com/grafana/grafana-foundation-sdk/go/dashboard"

const maxGridWidth uint32 = 24

type GridMap struct {
	data      map[uint32]map[uint32]bool
	minHeight uint32
}

func NewGridMap() *GridMap {
	return &GridMap{data: make(map[uint32]map[uint32]bool)}
}

func (g *GridMap) IsFree(x, y, w, h uint32) bool {
	for dy := range h {
		for dx := range w {
			if g.data[y+dy] != nil && g.data[y+dy][x+dx] {
				return false
			}
		}
	}
	return true
}

func (g *GridMap) Reserve(x, y, w, h uint32) {
	for dy := range h {
		row := y + dy
		if g.data[row] == nil {
			g.data[row] = make(map[uint32]bool)
		}
		for dx := range w {
			g.data[row][x+dx] = true
		}
	}
}

func (g *GridMap) FindNextAvailablePosition(w, h uint32) (uint32, uint32) {
	for y := g.minHeight; ; y++ {
		isRowFull := true
		for x := uint32(0); x <= maxGridWidth-w; x++ {
			if !g.data[y][x] {
				isRowFull = false
			}
			if g.IsFree(x, y, w, h) {
				g.Reserve(x, y, w, h)
				return x, y
			}
		}
		if isRowFull {
			g.minHeight = y + 1
		}
	}
}

func ArrangePanels(builder *dashboard.Dashboard) {
	grid := NewGridMap()

	for _, panel := range builder.Panels {

		var gridPos *dashboard.GridPos
		if panel.Panel != nil {
			gridPos = panel.Panel.GridPos
		} else if panel.RowPanel != nil {
			gridPos = panel.RowPanel.GridPos
		}

		w, h := gridPos.W, gridPos.H
		gridPos.X, gridPos.Y = grid.FindNextAvailablePosition(w, h)
	}
}
