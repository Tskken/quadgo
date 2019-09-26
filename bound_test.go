package quadgo

import (
	"reflect"
	"testing"
)

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

func TestBound_IsEqual(t *testing.T) {
	type fields struct {
		lhs Bound
	}
	type args struct {
		rhs Bound
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "value given is equal",
			fields: fields{
				lhs: Bound{
					Min: Point{
						X: 0,
						Y: 0,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
				},
			},
			args: args{
				rhs: Bound{
					Min: Point{
						X: 0,
						Y: 0,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
				},
			},
			want: true,
		},
		{
			name: "value given is not equal",
			fields: fields{
				Bound{
					Min: Point{
						X: 0,
						Y: 0,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
				},
			},
			args: args{
				Bound{
					Min: Point{
						X: 10,
						Y: 10,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.lhs.IsEqual(tt.args.rhs); got != tt.want {
				t.Errorf("Point.IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBound_IsIntersectBound(t *testing.T) {
	type fields struct {
		Min    Point
		Max    Point
		Center Point
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
			},
			args: args{
				bounds: Bound{
					Min:    Point{5, 5},
					Max:    Point{15, 15},
					Center: Point{float64(15 - (10 / 2)), float64(15 - (10 / 2))},
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
			},
			args: args{
				bounds: Bound{
					Min:    Point{55, 55},
					Max:    Point{105, 105},
					Center: Point{float64(105 - 25), float64(105 - 25)},
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
			}
			if got := b.IsIntersectPoint(tt.args.point); got != tt.want {
				t.Errorf("Bound.IsIntersectPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestBound_String(t *testing.T) {
// 	type fields struct {
// 		Min    Point
// 		Max    Point
// 		Center Point
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			name: "bound to string",
// 			fields: fields{
// 				Min:    Point{0, 0},
// 				Max:    Point{50, 50},
// 				Center: Point{25, 25},
// 			},
// 			want: "Min: X: 0, Y: 0, Max: X: 50, Y: 50, Center: X: 25, Y: 25\n",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := Bound{
// 				Min:    tt.fields.Min,
// 				Max:    tt.fields.Max,
// 				Center: tt.fields.Center,
// 			}
// 			if got := b.String(); got != tt.want {
// 				t.Errorf("Bound.String() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
