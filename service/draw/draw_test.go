package draw

import "testing"

func TestExtractImageURL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name:  "string output",
			input: "https://example.com/image.png",
			want:  "https://example.com/image.png",
		},
		{
			name:  "array output",
			input: []interface{}{"https://example.com/image.png"},
			want:  "https://example.com/image.png",
		},
		{
			name:    "empty array",
			input:   []interface{}{},
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   map[string]interface{}{"url": "https://example.com/image.png"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractImageURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
