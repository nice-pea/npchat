package jwt_create

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_JWTC_Issue(t *testing.T) {
	t.Run("uid и sid могут содержать любые строковые данные", func(t *testing.T) {
		jwtC := Issuer{"secret"}

		tests := []struct {
			name   string
			claims map[string]any
		}{
			{
				"claims может иметь пустые параметры",
				map[string]any{
					"UserID":    "",
					"SessionID": "",
				},
			},
			{
				"claims может иметь любые параметры",
				map[string]any{
					"UserID": "asfjkjasfjasfl;klfsaklfklsasakfl;kflsa",
					"itakoe": "123456789",
				},
			},
			{
				"claims может быть пустым",
				map[string]any{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				token, err := jwtC.Issue(tt.claims)
				assert.NoError(t, err)
				assert.NotZero(t, token)
			})
		}
	})

	t.Run("jwt токены созданные с одинковыми даными но с разными ttl - неравны", func(t *testing.T) {
		var (
			claims1 = map[string]any{
				"UserID":    "123",
				"SessionID": "123",
				"exp":       time.Now().Add(2 * time.Second).Unix(),
			}
			claims2 = map[string]any{
				"UserID":    "123",
				"SessionID": "123",
				"exp":       time.Now().Add(1 * time.Second).Unix(),
			}
		)

		jwtC := Issuer{"secret"}
		t1, _ := jwtC.Issue(claims1)
		t2, _ := jwtC.Issue(claims2)
		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с разными secret - неравны", func(t *testing.T) {
		var claims = map[string]any{
			"UserID":    "123",
			"SessionID": "123",
		}
		jwtC := Issuer{"secret1"}
		t1, _ := jwtC.Issue(claims)

		jwtC = Issuer{"secret2"}
		t2, _ := jwtC.Issue(claims)
		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с zero secret - невалидны", func(t *testing.T) {
		var claims = map[string]any{
			"UserID":    "123",
			"SessionID": "123",
		}
		jwtC := Issuer{}
		token, err := jwtC.Issue(claims)
		assert.Error(t, err)
		assert.Zero(t, token)
	})

	t.Run("jwt токены созданные с разными даными - неравны", func(t *testing.T) {
		var (
			claims1 = map[string]any{
				"UserID":    "123",
				"SessionID": "123",
			}
			claims2 = map[string]any{
				"UserID":    "123",
				"SessionID": "1234",
			}
		)

		jwtC := Issuer{"secret"}
		t1, _ := jwtC.Issue(claims1)
		t2, _ := jwtC.Issue(claims2)
		assert.NotEqual(t, t1, t2)
	})
}
