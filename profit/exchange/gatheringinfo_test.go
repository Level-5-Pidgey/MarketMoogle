package exchange

import (
	"math"
	"testing"
)

func withinTolerance(a, b, e float64) bool {
	if a == b {
		return true
	}

	d := math.Abs(a - b)

	if b == 0 {
		return d < e
	}

	return (d / math.Abs(b)) < e
}

func TestGatheringInfo_GetCost(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  int
	}{
		{
			name:  "Level cap Endwalker gathering costs 1650 gil",
			level: 90,
			want:  1650,
		},
		{
			name:  "Gathering above mid-Stormblood still costs 1275 gil",
			level: 75,
			want:  1275,
		},
		{
			name:  "Early game ARR gathering costs 375 gil",
			level: 20,
			want:  150,
		},
		{
			name:  "Level cap ARR gathering still costs 600 gil",
			level: 50,
			want:  600,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gatheringInfo := GatheringInfo{
					Level: tt.level,
				}
				if got := gatheringInfo.GetCost(); got != tt.want {
					t.Errorf("GetCost() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGatheringInfo_GetCostPerItem(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		isCollectible bool
		isHidden      bool
		want          int
	}{
		{
			name:  "Level cap Endwalker gathering costs approx 61 gil per item",
			level: 90,
			want:  61,
		},
		{
			name:          "Level cap Endwalker collectible gathering costs approx 550 gil per item",
			level:         90,
			isCollectible: true,
			want:          550,
		},
		{
			name:  "Mid Stormblood gathering costs 41 gil per item",
			level: 75,
			want:  41,
		},
		{
			name:          "Mid Stormblood collectible gathering costs 425 gil per item",
			level:         75,
			isCollectible: true,
			want:          425,
		},
		{
			name:     "Level cap hidden Endwalker gathering costs 68 gil per item",
			level:    90,
			isHidden: true,
			want:     68,
		},
		{
			name:          "Level cap ARR gathering costs per 18 gil per item",
			level:         50,
			isCollectible: false,
			isHidden:      false,
			want:          18,
		},
		{
			name:          "Level cap ARR hidden gathering costs 20 gil per item",
			level:         50,
			isCollectible: false,
			isHidden:      true,
			want:          20,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gatheringInfo := GatheringInfo{
					Level:         tt.level,
					IsCollectible: tt.isCollectible,
					IsHidden:      tt.isHidden,
				}
				if got := gatheringInfo.GetCostPerItem(); got != tt.want {
					t.Errorf("GetCostPerItem() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGatheringInfo_GetEffortFactor(t *testing.T) {
	type fields struct {
		Level         int
		IsCollectible bool
		IsHidden      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Endwalker collectibles are at an effort factor of 1.1",
			fields: fields{
				Level:         90,
				IsCollectible: true,
				IsHidden:      false,
			},
			want: 1.1,
		},
		{
			name: "Endwalker gatherables are at an effort factor of 1.1",
			fields: fields{
				Level:         90,
				IsCollectible: false,
				IsHidden:      false,
			},
			want: 1.1,
		},
		{
			name: "Shadowbringers gatherables are at an effort factor of 1.0",
			fields: fields{
				Level:         80,
				IsCollectible: false,
				IsHidden:      false,
			},
			want: 1.0,
		},
		{
			name: "Shadowbringers hidden gatherables are at an effort factor of 1.1",
			fields: fields{
				Level:         80,
				IsCollectible: false,
				IsHidden:      true,
			},
			want: 1.1,
		},
		{
			name: "Stormblood hidden gatherables are at an effort factor of 1.05",
			fields: fields{
				Level:         70,
				IsCollectible: false,
				IsHidden:      true,
			},
			want: 1.05,
		},
		{
			name: "Stormblood gatherables are at an effort factor of 0.95",
			fields: fields{
				Level:         70,
				IsCollectible: false,
				IsHidden:      false,
			},
			want: 0.95,
		},
		{
			name: "Heavensward gatherables are at an effort factor of 0.925",
			fields: fields{
				Level:         60,
				IsCollectible: false,
				IsHidden:      false,
			},
			want: 0.925,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gatheringInfo := GatheringInfo{
					Level:         tt.fields.Level,
					IsHidden:      tt.fields.IsHidden,
					IsCollectible: tt.fields.IsCollectible,
				}
				if got := gatheringInfo.GetEffortFactor(); !withinTolerance(got, tt.want, 1e-3) {
					t.Errorf("GetEffortFactor() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGatheringInfo_GetQuantity(t *testing.T) {
	type fields struct {
		Level         int
		IsCollectible bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Endwalker collectibles drop 3 per node",
			fields: fields{
				Level:         90,
				IsCollectible: true,
			},
			want: 3,
		},
		{
			name: "Endwalker gatherables drop approx. 27 per node",
			fields: fields{
				Level:         90,
				IsCollectible: false,
			},
			want: 27,
		},
		{
			name: "Shadowbringers gatherables drop approx. 30 per node",
			fields: fields{
				Level:         80,
				IsCollectible: false,
			},
			want: 30,
		},
		{
			name: "Shadowbringers hidden gatherables drop approx. 30 per node",
			fields: fields{
				Level:         80,
				IsCollectible: false,
			},
			want: 30,
		},
		{
			name: "Stormblood hidden gatherables drop approx. 31 per node",
			fields: fields{
				Level:         70,
				IsCollectible: false,
			},
			want: 31,
		},
		{
			name: "Stormblood gatherables drop approx. 31 per node",
			fields: fields{
				Level:         70,
				IsCollectible: false,
			},
			want: 31,
		},
		{
			name: "Heavensward gatherables drop approx. 32 per node",
			fields: fields{
				Level:         60,
				IsCollectible: false,
			},
			want: 32,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gatheringInfo := GatheringInfo{
					Level:         tt.fields.Level,
					IsCollectible: tt.fields.IsCollectible,
				}
				if got := gatheringInfo.GetQuantity(); got != tt.want {
					t.Errorf("GetQuantity() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
