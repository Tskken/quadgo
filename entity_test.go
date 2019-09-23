package quadgo

import (
	"reflect"
	"testing"
)

func TestNewEntity(t *testing.T) {
	type args struct {
		minX float64
		minY float64
		maxX float64
		maxY float64
		objs []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Entity
	}{
		{
			name: "basic new entity",
			args: args{
				0, 0, 25, 25, nil,
			},
			want: &Entity{
				Bound: Bound{
					Min: Point{
						0, 0,
					},
					Max: Point{
						25, 25,
					},
					Center: Point{
						float64(25) / 2, float64(25) / 2,
					},
				},
				Objects: nil,
			},
		},
		{
			name: "object function on new entity",
			args: args{
				0, 0, 25, 25, []interface{}{
					struct {
						name string
					}{"test struct object"},
				},
			},
			want: &Entity{
				Bound: Bound{
					Min: Point{
						0, 0,
					},
					Max: Point{
						25, 25,
					},
					Center: Point{
						float64(25) / 2, float64(25) / 2,
					},
					Width:  25,
					Height: 25,
				},
				Objects: []interface{}{
					struct {
						name string
					}{"test struct object"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEntity(tt.args.minX, tt.args.minY, tt.args.maxX, tt.args.maxY, tt.args.objs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_String(t *testing.T) {
	type fields struct {
		Bound   Bound
		Objects []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "entity to string",
			fields: fields{
				Bound: Bound{
					Min:    Point{0, 0},
					Max:    Point{50, 50},
					Center: Point{25, 25},
					Width:  50,
					Height: 50,
				},
				Objects: nil,
			},
			want: "Bounds: Min: X: 0, Y: 0, Max: X: 50, Y: 50, Center: X: 25, Y: 25\n Width: 50, Height 50\n\n Objects: []\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entity{
				Bound:   tt.fields.Bound,
				Objects: tt.fields.Objects,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("Entity.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBound(t *testing.T) {
	type args struct {
		minX float64
		minY float64
		maxX float64
		maxY float64
	}
	tests := []struct {
		name string
		args args
		want Bound
	}{
		{
			name: "new bounds",
			args: args{0, 0, 50, 50},
			want: Bound{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBound(tt.args.minX, tt.args.minY, tt.args.maxX, tt.args.maxY); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBound_IsIntersectBound(t *testing.T) {
	type fields struct {
		Min    Point
		Max    Point
		Center Point
		Width  float64
		Height float64
	}
	type args struct {
		bounds Bound
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "is intersected test",
			fields: fields{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
			args: args{
				bounds: Bound{
					Min:    Point{5, 5},
					Max:    Point{15, 15},
					Center: Point{float64(15 - (10 / 2)), float64(15 - (10 / 2))},
					Width:  10,
					Height: 10,
				},
			},
			want: true,
		},
		{
			name: "is not interacted test",
			fields: fields{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
			args: args{
				bounds: Bound{
					Min:    Point{55, 55},
					Max:    Point{105, 105},
					Center: Point{float64(105 - 25), float64(105 - 25)},
					Width:  25,
					Height: 25,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bound{
				Min:    tt.fields.Min,
				Max:    tt.fields.Max,
				Center: tt.fields.Center,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := b.IsIntersectBound(tt.args.bounds); got != tt.want {
				t.Errorf("Bound.IsIntersectBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBound_IsIntersectPoint(t *testing.T) {
	type fields struct {
		Min    Point
		Max    Point
		Center Point
		Width  float64
		Height float64
	}
	type args struct {
		point Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "is intersected test",
			fields: fields{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
			args: args{
				point: Point{5, 5},
			},
			want: true,
		},
		{
			name: "is not interacted test",
			fields: fields{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
			args: args{
				point: Point{55, 55},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bound{
				Min:    tt.fields.Min,
				Max:    tt.fields.Max,
				Center: tt.fields.Center,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := b.IsIntersectPoint(tt.args.point); got != tt.want {
				t.Errorf("Bound.IsIntersectPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBound_String(t *testing.T) {
	type fields struct {
		Min    Point
		Max    Point
		Center Point
		Width  float64
		Height float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "bound to string",
			fields: fields{
				Min:    Point{0, 0},
				Max:    Point{50, 50},
				Center: Point{25, 25},
				Width:  50,
				Height: 50,
			},
			want: "Min: X: 0, Y: 0, Max: X: 50, Y: 50, Center: X: 25, Y: 25\n Width: 50, Height 50\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bound{
				Min:    tt.fields.Min,
				Max:    tt.fields.Max,
				Center: tt.fields.Center,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := b.String(); got != tt.want {
				t.Errorf("Bound.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	type fields struct {
		X float64
		Y float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "point to string",
			fields: fields{5, 5},
			want:   "X: 5, Y: 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Point{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("Point.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
