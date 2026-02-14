# CLAUDE.md

## Swagger コメント規約

APIハンドラ関数（`func XxxHandler(c *gin.Context)` など）を定義する際は、関数の直前にswag形式のコメントを必ず記述すること。

```go
// @Summary ユーザー登録
// @Description 新規ユーザーを登録し、IDとパスワードを返す
// @Tags users
// @Accept json
// @Produce json
// @Param request body UserRegisterRequest true "ユーザー登録リクエスト"
// @Success 200 {object} UserRegisterResponse
// @Router /api/users/register [post]
func UserRegister(c *gin.Context) {
```

必須アノテーション:
- `@Summary` - APIの簡潔な説明
- `@Description` - APIの詳細な説明
- `@Tags` - グループ分け用タグ
- `@Accept` / `@Produce` - リクエスト/レスポンスのContent-Type
- `@Param` - リクエストパラメータ（body, query, path など）
- `@Success` - レスポンスのステータスコードと型(@Failureの場合は記載しない)
- `@Router` - エンドポイントのパスとHTTPメソッド(基本的に/apiから始まる)

認証が必要なエンドポイントには `@Security` も追加すること:
```go
// @Security BearerAuth
```