package decoder

import (
	"testing"

	"github.com/google/uuid"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		id   uuid.UUID
		want string
	}{
		{uuid.MustParse("c9c8b97f-8bf3-4a71-a0e2-42131a30bb90"), "YzljOGI5N2YtOGJmMy00YTcxLWEwZTItNDIxMzFhMzBiYjkw"},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			got, err := Decode(tt.id)
			if err != nil {
				t.Errorf("error got while decoding: %v", err)
			}

			// fmt.Println("--got", got)

			// id, err := Encode(got)
			// if err != nil {
			// 	t.Errorf("error got while encoding: %v", err)
			// }
			// fmt.Println("--uid", id)

			if got != tt.want {
				t.Errorf("expected %s, want %s", got, tt.want)
			}
		})
	}
}
