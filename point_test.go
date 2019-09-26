package quadgo

import "testing"

func TestNewPoint(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		{
			name: "basic new point",
			args: args{
				x: 50,
				y: 50,
			},
			want: Point{
				X: 50,
				Y: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPoint(tt.args.x, tt.args.y)
			if got.X != tt.want.X || got.Y != tt.want.Y {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_IsEqual(t *testing.T) {
	type fields struct {
		lhs Point
	}
	type args struct {
		rhs Point
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
				lhs: Point{
					X: 50,
					Y: 50,
				},
			},
			args: args{
				rhs: Point{
					X: 50,
					Y: 50,
				},
			},
			want: true,
		},
		{
			name: "value given is not equal",
			fields: fields{
				lhs: Point{
					X: 50,
					Y: 50,
				},
			},
			args: args{
				rhs: Point{
					X: 60,
					Y: 60,
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

// func TestPoint_String(t *testing.T) {
// 	type fields struct {
// 		X float64
// 		Y float64
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			name:   "point to string",
// 			fields: fields{5, 5},
// 			want:   "X: 5, Y: 5",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := Point{
// 				X: tt.fields.X,
// 				Y: tt.fields.Y,
// 			}
// 			if got := p.String(); got != tt.want {
// 				t.Errorf("Point.String() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
