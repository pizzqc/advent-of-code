package main

import (
	"container/list"
	"reflect"
	"testing"
)

func TestPosition_GetNeighbors(t *testing.T) {
	tests := []struct {
		name string
		p    Position
		want []Position
	}{
		{
			name: "Test with 0,0 return node at N, E, S, W",
			p:    Position{X: 0, Y: 0},
			want: []Position{
				{X: -1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: -1},
			},
		},
		{
			name: "Test with 5,6 return node at N, E, S, W",
			p:    Position{X: 5, Y: 6},
			want: []Position{
				{X: 4, Y: 6},
				{X: 5, Y: 7},
				{X: 6, Y: 6},
				{X: 5, Y: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetNeighbors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.GetNeighbors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElevationMap_FindNeighbors(t *testing.T) {
	type args struct {
		row int
		col int
	}
	tests := []struct {
		name string
		e    *ElevationMap
		args args
		want []Position
	}{
		{
			name: "All valid - all smaller values",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 1, 1, 1},
					{1, 1, 4, 1, 1},
					{1, 1, 1, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 1,
				col: 2,
			},
			want: []Position{
				{X: 0, Y: 2},
				{X: 1, Y: 3},
				{X: 2, Y: 2},
				{X: 1, Y: 1},
			},
		},
		{
			name: "All valid - all equals",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 4, 1, 1},
					{1, 4, 4, 4, 1},
					{1, 1, 4, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 1,
				col: 2,
			},
			want: []Position{
				{X: 0, Y: 2},
				{X: 1, Y: 3},
				{X: 2, Y: 2},
				{X: 1, Y: 1},
			},
		},
		{
			name: "All valid - all plus 1",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 5, 1, 1},
					{1, 5, 4, 5, 1},
					{1, 1, 5, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 1,
				col: 2,
			},
			want: []Position{
				{X: 0, Y: 2},
				{X: 1, Y: 3},
				{X: 2, Y: 2},
				{X: 1, Y: 1},
			},
		},
		{
			name: "All invalid destination returns",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 5, 1, 1},
					{1, 5, 3, 5, 1},
					{1, 1, 5, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 1,
				col: 2,
			},
			want: []Position{},
		},
		{
			name: "Test out-of-bound with TopLeft",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 6, 1, 1},
					{1, 5, 4, 6, 1},
					{1, 1, 6, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 0,
				col: 0,
			},
			want: []Position{
				{X: 0, Y: 1},
				{X: 1, Y: 0},
			},
		},
		{
			name: "Test out-of-bound with TopRight",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 6, 1, 1},
					{1, 5, 4, 6, 1},
					{1, 1, 6, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 0,
				col: 4,
			},
			want: []Position{
				{X: 1, Y: 4},
				{X: 0, Y: 3},
			},
		},
		{
			name: "Test out-of-bound with LowerLeft",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 6, 1, 1},
					{1, 5, 4, 6, 1},
					{1, 1, 6, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 3,
				col: 0,
			},
			want: []Position{
				{X: 2, Y: 0},
				{X: 3, Y: 1},
			},
		},
		{
			name: "Test out-of-bound with BottomRight",
			e: &ElevationMap{
				Elevation: [][]int{
					{1, 1, 6, 1, 1},
					{1, 5, 4, 6, 1},
					{1, 1, 6, 1, 1},
					{1, 1, 1, 1, 1},
				},
			},
			args: args{
				row: 3,
				col: 4,
			},
			want: []Position{
				{X: 2, Y: 4},
				{X: 3, Y: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.FindNeighbors(tt.args.row, tt.args.col); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElevationMap.FindNeighbors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTraveler(t *testing.T) {
	type args struct {
		emap  ElevationMap
		start Position
	}
	tests := []struct {
		name string
		args args
		want *Traveler
	}{
		{
			name: "Init all with -1",
			args: args{
				emap: ElevationMap{
					Elevation: [][]int{
						{1, 2, 3, 4, 5},
						{1, 2, 3, 4, 5},
					}},
				start: Position{X: 1, Y: 2},
			},
			want: &Traveler{
				Elevation: ElevationMap{
					Elevation: [][]int{
						{1, 2, 3, 4, 5},
						{1, 2, 3, 4, 5},
					},
				},
				DistanceTracker: [][]int{
					{-1, -1, -1, -1, -1},
					{-1, -1, -1, -1, -1},
				},
				Start: Position{X: 1, Y: 2},
				Queue: &list.List{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTraveler(tt.args.emap, tt.args.start)
			if !reflect.DeepEqual(got.DistanceTracker, tt.want.DistanceTracker) &&
				!reflect.DeepEqual(got.Elevation, tt.want.Elevation) &&
				!reflect.DeepEqual(got.Start, tt.want.Start) {
				t.Errorf("NewTraveler() = %v, want %v", got, tt.want)
			}
			if got.Queue.Len() != 1 {
				t.Errorf("NewTraveler() Queue Size = %v, want %v", got.Queue.Len(), 1)
			}
		})
	}
}
