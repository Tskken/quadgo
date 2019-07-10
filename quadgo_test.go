package quadgo

import (
	"reflect"
	"testing"
)

func TestSetBounds(t *testing.T) {
	type args struct {
		width  float64
		height float64
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "basic test",
			args: args{
				width:  1920,
				height: 1080,
			},
			want: &Options{
				Width:       1920,
				Height:      1080,
				MaxEntities: defaultOption.MaxEntities,
				MaxDepth:    defaultOption.MaxDepth,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetBounds(tt.args.width, tt.args.height)
			op := defaultOption
			got(op)
			if !reflect.DeepEqual(op, tt.want) {
				t.Errorf("Options = %v, want %v", op, tt.want)
			}
		})
	}
}

func TestSetMaxEntities(t *testing.T) {
	type args struct {
		maxEntities int
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "basic test",
			args: args{
				maxEntities: 5,
			},
			want: &Options{
				Width:       defaultOption.Width,
				Height:      defaultOption.Height,
				MaxEntities: 5,
				MaxDepth:    defaultOption.MaxDepth,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetMaxEntities(5)
			op := defaultOption
			got(op)
			if !reflect.DeepEqual(op, tt.want) {
				t.Errorf("Options = %v, want %v", op, tt.want)
			}
		})
	}
}

func TestSetMaxDepth(t *testing.T) {
	type args struct {
		maxDepth int
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "basic test",
			args: args{
				maxDepth: 5,
			},
			want: &Options{
				Width:       defaultOption.Width,
				Height:      defaultOption.Height,
				MaxEntities: defaultOption.MaxEntities,
				MaxDepth:    5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetMaxDepth(5)
			op := defaultOption
			got(op)
			if !reflect.DeepEqual(op, tt.want) {
				t.Errorf("Options = %v, want %v", op, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ops []Option
	}
	tests := []struct {
		name string
		args args
		want *QuadGo
	}{
		{
			name: "default new",
			want: &QuadGo{
				node: &node{
					parent:   nil,
					bounds:   NewBound(0, 0, defaultOption.Width, defaultOption.Height),
					entities: make(Entities, 0, defaultOption.MaxEntities),
					children: make(nodes, 0, 4),
					depth:    0,
				},
				maxDepth: defaultOption.MaxDepth,
			},
		},
		{
			name: "SetBounds on New",
			args: args{
				ops: []Option{
					SetBounds(1920, 1080),
				},
			},
			want: &QuadGo{
				node: &node{
					parent:   nil,
					bounds:   NewBound(0, 0, 1920, 1080),
					entities: make(Entities, 0, defaultOption.MaxEntities),
					children: make(nodes, 0, 4),
					depth:    0,
				},
				maxDepth: defaultOption.MaxDepth,
			},
		},
		{
			name: "Set all options test",
			args: args{
				ops: []Option{
					SetBounds(1920, 1080),
					SetMaxEntities(5),
					SetMaxDepth(1),
				},
			},
			want: &QuadGo{
				node: &node{
					parent:   nil,
					bounds:   NewBound(0, 0, 1920, 1080),
					entities: make(Entities, 0, 5),
					children: make(nodes, 0, 4),
					depth:    0,
				},
				maxDepth: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.ops...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_Insert(t *testing.T) {
	type fields struct {
		options []Option
	}
	type args struct {
		minX float64
		minY float64
		maxX float64
		maxY float64
		objs []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    []args
		wantErr bool
	}{
		{
			name: "basic insert",
			args: []args{
				{0, 0, 50, 50, nil},
			},
			wantErr: false,
		},
		{
			name: "insert to split",
			fields: fields{
				[]Option{
					SetMaxEntities(1),
				},
			},
			args: []args{
				{0, 0, 50, 50, nil},
				{500, 500, 700, 700, nil},
				{700, 800, 900, 1000, nil},
			},
			wantErr: false,
		},
		{
			name: "insert all 4 quadrants",
			fields: fields{
				options: []Option{
					SetMaxEntities(1),
					SetBounds(100, 100),
				},
			},
			args: []args{
				{0, 0, 50, 50, nil},
				{51, 51, 100, 100, nil},
				{0, 51, 50, 100, nil},
				{51, 0, 100, 50, nil},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, a := range tt.args {
				if err := q.Insert(a.minX, a.minY, a.maxX, a.maxY, a.objs...); (err != nil) != tt.wantErr {
					t.Errorf("QuadGo.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func TestQuadGo_InsertEntity(t *testing.T) {
	type fields struct {
		options []Option
	}
	type args struct {
		entities Entities
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic insert",
			args: args{
				entities: Entities{
					NewEntity(0,0,50,50),
				},
			},
			wantErr: false,
		},
		{
			name:"no entities given error",
			wantErr:true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			if err := q.InsertEntity(tt.args.entities...); (err != nil) != tt.wantErr {
				t.Errorf("QuadGo.InsertEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuadGo_Remove(t *testing.T) {
	type fields struct {
		options  []Option
		entities Entities
	}
	type args struct {
		entities Entities
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic remove entity",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			wantErr: false,
		},
		{
			name: "no entity found in tree",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				entities: Entities{
					NewEntity(0, 0, 25, 25),
				},
			},
			wantErr: true,
		},
		{
			name: "remove from deeper node",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(55, 55, 105, 105),
					NewEntity(400, 400, 600, 600),
				},
				options: []Option{
					SetMaxEntities(1),
				},
			},
			args: args{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			wantErr: false,
		},
		{
			name: "collapse test",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(100, 100, 150, 150),
					NewEntity(200, 200, 250, 250),
				},
				options: []Option{
					SetMaxEntities(2),
				},
			},
			args: args{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(100, 100, 150, 150),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.Remove() error = %v", err)
				}
			}
			for _, e := range tt.args.entities {
				if err := q.Remove(e); (err != nil) != tt.wantErr {
					t.Errorf("QuadGo.Remove() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestQuadGo_RetrieveFromPoint(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		point Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Entities
	}{
		{
			name: "basic retrieve from point",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				point: Point{5, 5},
			},
			want: Entities{
				NewEntity(0, 0, 50, 50),
			},
		},
		{
			name: "find none from retrieve from point",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(5, 5, 55, 55),
				},
				options: []Option{
					SetBounds(100, 100),
					SetMaxEntities(1),
				},
			},
			args: args{
				point: Point{100, 100},
			},
			want: Entities{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.RetrieveFromPoint() for %v", e)
				}
			}
			if got := q.RetrieveFromPoint(tt.args.point); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuadGo.RetrieveFromPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_RetrieveFromBound(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		bound Bound
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Entities
	}{
		{
			name: "basic retrieve from bound",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				bound: NewBound(5, 5, 10, 10),
			},
			want: Entities{
				NewEntity(0, 0, 50, 50),
			},
		},
		{
			name: "find none from retrieve from bound",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(5, 5, 55, 55),
				},
				options: []Option{
					SetBounds(100, 100),
					SetMaxEntities(1),
				},
			},
			args: args{
				bound: NewBound(80, 80, 100, 100),
			},
			want: Entities{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.RetrieveFromBound() for %v", e)
				}
			}
			if got := q.RetrieveFromBound(tt.args.bound); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuadGo.RetrieveFromBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_IsEntity(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		entity *Entity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "basic IsEntity",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(100, 100, 150, 150),
				},
			},
			args: args{
				NewEntity(0, 0, 50, 50),
			},
			want: true,
		},
		{
			name: "not found for IsEntity",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
					NewEntity(100, 100, 150, 150),
				},
			},
			args: args{
				NewEntity(5, 5, 50, 50),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.IsEntity() for %v", e)
				}
			}
			if got := q.IsEntity(tt.args.entity); got != tt.want {
				t.Errorf("QuadGo.IsEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_IsIntersectPoint(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
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
			name: "basic IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				point: Point{5, 5},
			},
			want: true,
		},
		{
			name: "is not IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				point: Point{55, 55},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.IsIntersectPoint() for %v", e)
				}
			}
			if got := q.IsIntersectPoint(tt.args.point); got != tt.want {
				t.Errorf("QuadGo.IsIntersectPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_IsIntersectBound(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		bound Bound
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "basic IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				bound: NewBound(5, 5, 10, 10),
			},
			want: true,
		},
		{
			name: "is not IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				bound: NewBound(55, 55, 60, 60),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.IsIntersectPoint() for %v", e)
				}
			}
			if got := q.IsIntersectBound(tt.args.bound); got != tt.want {
				t.Errorf("QuadGo.IsIntersectBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuadGo_IntersectsPoint(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		point Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Entities
	}{
		{
			name: "basic IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				point: Point{5, 5},
			},
			want: Entities{
				NewEntity(0, 0, 50, 50),
			},
		},
		{
			name: "is not IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				point: Point{55, 55},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.IsIntersectPoint() for %v", e)
				}
			}
			if gotIntersects := q.IntersectsPoint(tt.args.point); !reflect.DeepEqual(gotIntersects, tt.want) {
				t.Errorf("QuadGo.IntersectsPoint() = %v, want %v", gotIntersects, tt.want)
			}
		})
	}
}

func TestQuadGo_IntersectsBound(t *testing.T) {
	type fields struct {
		entities Entities
		options  []Option
	}
	type args struct {
		bound Bound
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Entities
	}{
		{
			name: "basic IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				bound: NewBound(5, 5, 10, 10),
			},
			want: Entities{
				NewEntity(0, 0, 50, 50),
			},
		},
		{
			name: "is not IsIntersectPoint",
			fields: fields{
				entities: Entities{
					NewEntity(0, 0, 50, 50),
				},
			},
			args: args{
				bound: NewBound(55, 55, 60, 60),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.fields.options...)
			for _, e := range tt.fields.entities {
				err := q.InsertEntity(e)
				if err != nil {
					t.Errorf("error in QuadGo.InsertEntity() found in QuadGo.IsIntersectPoint() for %v", e)
				}
			}
			if gotIntersects := q.IntersectsBound(tt.args.bound); !reflect.DeepEqual(gotIntersects, tt.want) {
				t.Errorf("QuadGo.IntersectsBound() = %v, want %v", gotIntersects, tt.want)
			}
		})
	}
}
