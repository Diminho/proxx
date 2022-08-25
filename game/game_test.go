package game

import (
	"github.com/golang/mock/gomock"
	"github.com/proxx/game/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestGame_Start(t *testing.T) {
	type fields struct {
		playground func(ctrl *gomock.Controller) Playground
		state      State
		userInput  string
	}
	tests := []struct {
		name          string
		fields        fields
		setupInput    func() *os.File
		teardownInput func(*os.File)
		expectedErr   error
		wantErr       bool
	}{
		{
			name: "success_lose",
			fields: fields{
				playground: func(ctrl *gomock.Controller) Playground {
					p := mocks.NewMockPlayground(ctrl)
					p.EXPECT().Print().Return().Times(2)
					p.EXPECT().Click([]int{3, 4}).Return(nil)

					return p
				},
				userInput: "4 5\n",
				state:     lose,
			},
			setupInput: func() *os.File {
				// simulating user input
				in, err := os.CreateTemp("", "in")
				if err != nil {
					t.Fatal(err)
				}

				_, err = io.WriteString(in, "4 5\n")
				if err != nil {
					t.Fatal(err)
				}

				_, err = in.Seek(0, io.SeekStart)
				if err != nil {
					t.Fatal(err)
				}

				return in
			},
			teardownInput: func(in *os.File) {
				in.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			g := &Game{
				playground: tt.fields.playground(ctrl),
				state:      tt.fields.state,
			}

			// simulating user input
			in := tt.setupInput()
			defer tt.teardownInput(in)

			err := g.Start(in)
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
