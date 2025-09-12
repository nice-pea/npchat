package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_JWTC_Create(t *testing.T) {
	t.Run("uid и sid могут содержать любые строковые данные", func(t *testing.T) {
		jwtC := NewJWTCreator("secret", 2*time.Minute)

		tests := []struct {
			uid string
			sid string
		}{
			{
				uid: "",
				sid: "",
			},
			{
				uid: "флыыфжвлфывджл в 1лш21лд4л124 лл12д4лд214лл241",
				sid: strings.Repeat("a", 51),
			},
			{
				uid: "    ",
				sid: " 1 2 3 ",
			},
			{
				uid: "fakj35jkl1jkj1k2jlk412",
				sid: "aflklkrflk;1k;124kkl;124l24",
			},
		}
		for _, tt := range tests {
			token, err := jwtC.Create(tt.uid, tt.sid)
			assert.NoError(t, err)
			assert.NotZero(t, token)

		}
	})
	t.Run("jwt токены созданные с одинковыми даными, созданные в разные промежутки времени - неравны", func(t *testing.T) {
		jwtC := NewJWTCreator("secret", 2*time.Minute)
		var (
			uid = "123"
			sid = "123"
		)
		t1, _ := jwtC.Create(uid, sid)
		time.Sleep(1 * time.Second)
		t2, _ := jwtC.Create(uid, sid)
		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с одинковыми даными и в одно время - равны", func(t *testing.T) {
		jwtC := NewJWTCreator("secret", 2*time.Minute)
		var (
			uid = "123"
			sid = "123"
		)
		t1, _ := jwtC.Create(uid, sid)
		t2, _ := jwtC.Create(uid, sid)
		assert.Equal(t, t1, t2)
	})
	t.Run("jwt токены созданные с одинковыми даными но с разными ttl - неравны", func(t *testing.T) {
		var (
			uid = "123"
			sid = "123"
		)
		jwtC := NewJWTCreator("secret", 2*time.Minute)
		t1, _ := jwtC.Create(uid, sid)

		jwtC = NewJWTCreator("secret", 3*time.Minute)
		t2, _ := jwtC.Create(uid, sid)
		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с разными secret - неравны", func(t *testing.T) {
		var (
			uid = "123"
			sid = "123"
		)
		jwtC := NewJWTCreator("secret1", 2*time.Minute)
		t1, _ := jwtC.Create(uid, sid)

		jwtC = NewJWTCreator("secret2", 2*time.Minute)
		t2, _ := jwtC.Create(uid, sid)
		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с zero secret - невалидны", func(t *testing.T) {
		var (
			uid = "123"
			sid = "123"
		)
		jwtC := NewJWTCreator("", 2*time.Minute)
		_, err := jwtC.Create(uid, sid)
		assert.Error(t, err)
	})
}
