# ビルドステージ
FROM golang:1.25.4-alpine AS builder

# CGOビルドに必要なツールをインストール
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/cat-backend ./main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/cat-backend .

# dbディレクトリを作成
RUN mkdir -p /app/db

EXPOSE 8080

CMD ["./cat-backend"]

