package quadgo

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestEntities_FindAndRemove(t *testing.T) {
	type fields struct {
		entities Entities
	}
	type args struct {
		entity *Entity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Entities
		wantErr error
	}{
		{
			name: "basic remove from list of 3",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 60,
								Y: 60,
							},
							Max: Point{
								X: 800,
								Y: 800,
							},
						},
						Action: nil,
					},
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 20,
								Y: 20,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				entity: &Entity{
					Bound: Bound{
						Min: Point{
							X: 60,
							Y: 60,
						},
						Max: Point{
							X: 800,
							Y: 800,
						},
					},
					Action: nil,
				},
			},
			want: Entities{
				&Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
				&Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 20,
							Y: 20,
						},
					},
					Action: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "remove from only 1 item list",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				entity: &Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			want:    Entities{},
			wantErr: nil,
		},
		{
			name: "remove last item in list",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 60,
								Y: 60,
							},
							Max: Point{
								X: 100,
								Y: 100,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				entity: &Entity{
					Bound: Bound{
						Min: Point{
							X: 60,
							Y: 60,
						},
						Max: Point{
							X: 100,
							Y: 100,
						},
					},
					Action: nil,
				},
			},
			want: Entities{
				&Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "could not find item in list error",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 60,
								Y: 60,
							},
							Max: Point{
								X: 100,
								Y: 100,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				entity: &Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 60,
						},
						Max: Point{
							X: 100,
							Y: 100,
						},
					},
					Action: nil,
				},
			},
			want:    nil,
			wantErr: errors.New("could not find entity in tree to remove"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent, err := tt.fields.entities.FindAndRemove(tt.args.entity)
			if !reflect.DeepEqual(ent, tt.want) || !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Entities.FindAndRemove() = %v, %v, wanted %v, %v", ent, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestEntities_Contains(t *testing.T) {
	type fields struct {
		entities Entities
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
			name: "basic find true",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				&Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			want: true,
		},
		{
			name: " not found",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				&Entity{
					Bound: Bound{
						Min: Point{
							X: 10,
							Y: 10,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.entities.Contains(tt.args.entity); got != tt.want {
				t.Errorf("Entities.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntities_isIntersect(t *testing.T) {
	type fields struct {
		entities Entities
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
			name: "intersect bound true",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
							Center: Point{
								X: 25,
								Y: 25,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				bound: Bound{
					Min: Point{
						X: 10,
						Y: 10,
					},
					Max: Point{
						X: 40,
						Y: 40,
					},
					Center: Point{
						X: 25,
						Y: 25,
					},
				},
			},
			want: true,
		},
		{
			name: "intersect bound false",
			fields: fields{
				entities: Entities{
					&Entity{
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
							Center: Point{
								X: 25,
								Y: 25,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				bound: Bound{
					Min: Point{
						X: 100,
						Y: 100,
					},
					Max: Point{
						X: 150,
						Y: 150,
					},
					Center: Point{
						X: 125,
						Y: 125,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.entities.isIntersect(tt.args.bound); got != tt.want {
				t.Errorf("Entities.isIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntities_intersects(t *testing.T) {
	type fields struct {
		entities Entities
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
			name: "intersect bound true",
			fields: fields{
				entities: Entities{
					&Entity{
						ID: 1,
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
							Center: Point{
								X: 25,
								Y: 25,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				bound: Bound{
					Min: Point{
						X: 10,
						Y: 10,
					},
					Max: Point{
						X: 40,
						Y: 40,
					},
					Center: Point{
						X: 25,
						Y: 25,
					},
				},
			},
			want: Entities{
				&Entity{
					ID: 1,
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
						Center: Point{
							X: 25,
							Y: 25,
						},
					},
					Action: nil,
				},
			},
		},
		{
			name: "intersect bound false",
			fields: fields{
				entities: Entities{
					&Entity{
						ID: 1,
						Bound: Bound{
							Min: Point{
								X: 0,
								Y: 0,
							},
							Max: Point{
								X: 50,
								Y: 50,
							},
							Center: Point{
								X: 25,
								Y: 25,
							},
						},
						Action: nil,
					},
				},
			},
			args: args{
				bound: Bound{
					Min: Point{
						X: 100,
						Y: 100,
					},
					Max: Point{
						X: 150,
						Y: 150,
					},
					Center: Point{
						X: 125,
						Y: 125,
					},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.entities.intersects(tt.args.bound)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Entities.intersectBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntity(t *testing.T) {
	type args struct {
		minX float64
		minY float64
		maxX float64
		maxY float64
	}
	tests := []struct {
		name string
		args args
		want *Entity
	}{
		{
			name: "basic new entity",
			args: args{
				0, 0, 25, 25,
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
				Action: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEntity(tt.args.minX, tt.args.minY, tt.args.maxX, tt.args.maxY); !got.Bound.IsEqual(tt.want.Bound) {
				t.Errorf("NewEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntityWithAction(t *testing.T) {
	type args struct {
		minX   float64
		minY   float64
		maxX   float64
		maxY   float64
		action Action
	}
	tests := []struct {
		name string
		args args
		want *Entity
	}{
		{
			name: "basic new entity",
			args: args{
				minX:   0,
				minY:   0,
				maxX:   50,
				maxY:   50,
				action: func() { fmt.Println("this is a test func on action") },
			},
			want: &Entity{
				Bound: Bound{
					Min: Point{
						X: 0,
						Y: 0,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
					Center: Point{
						X: 25,
						Y: 25,
					},
				},
				Action: func() { fmt.Println("this is a test func on action") },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEntityWithAction(tt.args.minX, tt.args.minY, tt.args.maxX, tt.args.maxY, tt.args.action); !got.Bound.IsEqual(tt.want.Bound) || got.Action == nil {
				t.Errorf("NewEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_SetAction(t *testing.T) {
	type fields struct {
		entity *Entity
	}
	type args struct {
		action Action
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Entity
	}{
		{
			name: "basic set action function",
			fields: fields{
				entity: &Entity{
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			args: args{
				action: func() { fmt.Println("test set action action function.") },
			},
			want: &Entity{
				Bound: Bound{
					Min: Point{
						X: 0,
						Y: 0,
					},
					Max: Point{
						X: 50,
						Y: 50,
					},
				},
				Action: func() { fmt.Println("test set action action function.") },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.entity.SetAction(tt.args.action); tt.fields.entity.Action == nil {
				t.Errorf("Entity.SetAction() = %v, want %v", tt.fields.entity, tt.want)
			}
		})
	}
}

func TestEntity_IsEqual(t *testing.T) {
	type fields struct {
		entity *Entity
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
			name: "is equal true",
			fields: fields{
				entity: &Entity{
					ID: 1,
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			args: args{
				entity: &Entity{
					ID: 1,
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			want: true,
		},
		{
			name: "is equal false",
			fields: fields{
				entity: &Entity{
					ID: 2,
					Bound: Bound{
						Min: Point{
							X: 0,
							Y: 0,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			args: args{
				entity: &Entity{
					ID: 1,
					Bound: Bound{
						Min: Point{
							X: 10,
							Y: 10,
						},
						Max: Point{
							X: 50,
							Y: 50,
						},
					},
					Action: nil,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.entity.IsEqual(tt.args.entity); got != tt.want {
				t.Errorf("Entity.IsEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestEntity_String(t *testing.T) {
// 	type fields struct {
// 		Bound  Bound
// 		action func()
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			name: "entity to string",
// 			fields: fields{
// 				Bound: Bound{
// 					Min:    Point{0, 0},
// 					Max:    Point{50, 50},
// 					Center: Point{25, 25},
// 				},
// 				action: nil,
// 			},
// 			want: "Bounds: Min: X: 0, Y: 0, Max: X: 50, Y: 50, Center: X: 25, Y: 25\n Action: <nil>",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &Entity{
// 				Bound:  tt.fields.Bound,
// 				Action: tt.fields.action,
// 			}
// 			if got := e.String(); got != tt.want {
// 				t.Errorf("Entity.String() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
