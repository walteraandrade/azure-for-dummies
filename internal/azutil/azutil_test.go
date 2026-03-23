package azutil

import "testing"

func ptr[T any](v T) *T { return &v }

func TestDerefStr(t *testing.T) {
	tests := []struct {
		name string
		in   *string
		want string
	}{
		{"nil", nil, ""},
		{"value", ptr("hello"), "hello"},
		{"zero", ptr(""), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DerefStr(tt.in); got != tt.want {
				t.Errorf("DerefStr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDerefInt32(t *testing.T) {
	tests := []struct {
		name string
		in   *int32
		want int32
	}{
		{"nil", nil, 0},
		{"value", ptr(int32(42)), 42},
		{"zero", ptr(int32(0)), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DerefInt32(tt.in); got != tt.want {
				t.Errorf("DerefInt32() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDerefFloat64(t *testing.T) {
	tests := []struct {
		name string
		in   *float64
		want float64
	}{
		{"nil", nil, 0},
		{"value", ptr(3.14), 3.14},
		{"zero", ptr(0.0), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DerefFloat64(tt.in); got != tt.want {
				t.Errorf("DerefFloat64() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestDerefBool(t *testing.T) {
	tests := []struct {
		name string
		in   *bool
		want bool
	}{
		{"nil", nil, false},
		{"true", ptr(true), true},
		{"false", ptr(false), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DerefBool(tt.in); got != tt.want {
				t.Errorf("DerefBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRG(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{"valid", "/subscriptions/sub-id/resourceGroups/my-rg/providers/Microsoft.App/containerApps/my-app", "my-rg"},
		{"empty", "", ""},
		{"no-rg", "/subscriptions/sub-id/providers/Microsoft.App", ""},
		{"case-insensitive", "/subscriptions/sub-id/RESOURCEGROUPS/My-RG/providers/foo", "My-RG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractRG(tt.id); got != tt.want {
				t.Errorf("ExtractRG() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{"valid", "/subscriptions/sub-id/resourceGroups/rg/providers/Microsoft.App/containerApps/my-app", "my-app"},
		{"empty", "", ""},
		{"single", "only-segment", "only-segment"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractName(tt.id); got != tt.want {
				t.Errorf("ExtractName() = %q, want %q", got, tt.want)
			}
		})
	}
}
