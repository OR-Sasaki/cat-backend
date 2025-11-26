package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
	"github.com/OR-Sasaki/cat-backend/models"
	"github.com/OR-Sasaki/cat-backend/test"
)

func TestGetAllOutfits(t *testing.T) {
	test.TestApi(t, getAllOutfitsTestDatas(), "/api/outfits", OutfitsRouter, true, http.MethodGet)
}

func TestGetAllOutfitsUnauthorized(t *testing.T) {
	testDatas := []test.TestData[[]OutfitResponse]{
		{
			Name:           "異常系: 認証なし",
			RequestBody:    nil,
			ExpectedStatus: http.StatusUnauthorized,
		},
	}

	test.TestApi(t, testDatas, "/api/outfits", OutfitsRouter, false, http.MethodGet)
}

func getAllOutfitsTestDatas() []test.TestData[[]OutfitResponse] {
	testDatas := []test.TestData[[]OutfitResponse]{}

	// 正常系: 認証ありで全outfitsを取得
	{
		testDatas = append(testDatas, test.TestData[[]OutfitResponse]{
			Before: func(t *testing.T, _ *models.User) map[string]any {
				series := models.Series{
					Name: "testseries",
				}
				gorm.G[models.Series](config.DB).Create(t.Context(), &series)

				outfit1 := models.Outfit{
					Name:      "outfit1",
					Type:      models.OutfitTypeDefault,
					SeriesID:  series.ID,
					AssetPath: "/path/to/outfit1",
				}
				gorm.G[models.Outfit](config.DB).Create(t.Context(), &outfit1)

				outfit2 := models.Outfit{
					Name:      "outfit2",
					Type:      models.OutfitTypeDefault,
					SeriesID:  series.ID,
					AssetPath: "/path/to/outfit2",
				}
				gorm.G[models.Outfit](config.DB).Create(t.Context(), &outfit2)

				return nil
			},
			Name:           "正常系: 認証ありで全outfitsを取得",
			RequestBody:    nil,
			ExpectedStatus: http.StatusOK,
			ValidateResponse: func(t *testing.T, w *httptest.ResponseRecorder, responseBody []OutfitResponse) {
				assert.Equal(t, 2, len(responseBody), "outfitsの数は2である必要があります")
				assert.Equal(t, "outfit1", responseBody[0].Name)
				assert.Equal(t, "default", responseBody[0].Type)
				assert.Equal(t, "/path/to/outfit1", responseBody[0].AssetPath)
			},
		})
	}

	return testDatas
}
