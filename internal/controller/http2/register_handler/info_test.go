package registerHandler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/common"
)

func TestInfo(t *testing.T) {
	t.Run("успешное получение информации о сборке", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})

		buildInfo := common.BuildInfo{
			Version:   "1.2.3",
			BuildDate: "2025-10-02",
			Commit:    "abc123def456",
		}

		Info(fiberApp, buildInfo)

		req := httptest.NewRequest("GET", "/info", nil)
		resp, err := fiberApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "1.2.3", result["version"])
		assert.Equal(t, "2025-10-02", result["build_date"])
		assert.Equal(t, "abc123def456", result["commit"])
	})

	t.Run("пустые значения сборки", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})

		buildInfo := common.BuildInfo{
			Version:   "",
			BuildDate: "",
			Commit:    "",
		}

		Info(fiberApp, buildInfo)

		req := httptest.NewRequest("GET", "/info", nil)
		resp, err := fiberApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "", result["version"])
		assert.Equal(t, "", result["build_date"])
		assert.Equal(t, "", result["commit"])
	})
}
