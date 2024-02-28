package service

import (
	"context"
	"errors"
	"kvart-info/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_GetSummaryInfo(t *testing.T) {

	tests := []struct {
		name   string
		uc     *Usecase
		want   []model.SummaryInfo
		gotErr error
	}{
		{name: "пустая структура", want: []model.SummaryInfo{}, gotErr: nil},
		{name: "есть ошибка 1", want: nil, gotErr: errors.New("Get PrepareNamedContext")},
		{name: "есть ошибка 2", want: nil, gotErr: errors.New("Get SelectContext")},
		{name: "заполненная структура",
			want:   []model.SummaryInfo{{FinID: 1, TipName: "УК", CountBuild: 10, CountLic: 90}},
			gotErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mockRepository{}
			mockRepo.On("Get", mock.Anything).Return(tt.want, tt.gotErr)
			uc := New(&mockRepo)
			ctx := context.Background()

			got, err := uc.GetSummaryInfo(ctx)
			if tt.gotErr != nil {
				require.Error(t, err)
				require.Nil(t, got)
				require.EqualError(t, tt.gotErr, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
