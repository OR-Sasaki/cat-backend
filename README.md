# cat-backend

Go (Gin) + SQLite で構成された猫アプリのバックエンド API サーバー。

## 技術スタック

- Go 1.25.4
- Gin (HTTP フレームワーク)
- GORM + SQLite (データベース)
- JWT 認証
- Docker

## セットアップ & 起動

Docker と Docker Compose が必要です。

### 1. イメージをビルド

```bash
make build
```

### 2. コンテナを起動

```bash
make up
```

バックグラウンドでコンテナが起動し、`http://localhost:8080` でアクセスできます。

### 3. コンテナを停止

```bash
make down
```
