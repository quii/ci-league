package league

import (
	"reflect"
	"testing"
)

func TestExtractCoAuthor(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		want   []string
	}{
		{
			name: "With co-author",
			commit: `Remove default setting of anchor by useActiveNavItem

			Co-authored-by: LisaMcCormack <lisamccormack85@gmail.com>`,
			want: []string{"lisamccormack85@gmail.com"},
		},
		{
			name:   "Without co-author",
			commit: `refactoring`,
			want:   []string{""},
		},
		{
			name: "With several co-authors",
			commit: `Remove default setting of anchor by useActiveNavItem

			Co-authored-by: Eevee <eevee@letsgo.com> Co-authored-by: Pikachu <pika@pikapi.com> Co-authored-by: Charmander <charmer@nintendo.com>`,
			want: []string{"eevee@letsgo.com", "pika@pikapi.com", "charmer@nintendo.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCoAuthor(tt.commit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractCoAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}
