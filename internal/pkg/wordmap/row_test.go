package wordmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRow_FromString(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *Row
		wantErr bool
	}{
		{
			name: "ok",
			str:  "word\t1234",
			want: &Row{
				Word:    "word",
				Counter: 1234,
			},
			wantErr: false,
		},
		{
			name:    "empty str",
			str:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty word",
			str:     "\t1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "overflow",
			str:     "word\t9999999999",
			want:    nil,
			wantErr: true,
		},
		{
			name: "tab in word",
			str:  "my\tword\t1234",
			want: &Row{
				Word:    "my\tword",
				Counter: 1234,
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row := &Row{}
			err := row.FromString(test.str)

			if test.wantErr {
				assert.NotNil(t, err)
				return
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.want, row)
		})
	}
}
