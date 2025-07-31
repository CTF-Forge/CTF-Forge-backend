# CTFLab

誰もがCTFの問題を作成し、公開できるアプリ。
webとdesktopで公開する。

技術スタック

| コンポーネント             | 推奨技術                          | 主要な役割・選定理由                                                                 |
|----------------------------|-----------------------------------|----------------------------------------------------------------------------------------|
| バックエンド言語           | Go                                | ユーザー要件。高いパフォーマンスと並行処理性能。                                     |
| バックエンドフレームワーク | Gin                               | 高性能なHTTPルーター、豊富なミドルウェア、広範なコミュニティ。                        |
| デスクトップフレームワーク | Wails                             | Web技術を用いたUI開発、Goとの直接連携、軽量なバイナリ生成。                           |
| フロントエンドフレームワーク | Svelte                            | コンパイラベースによる高速なランタイム性能と軽量なバンドルサイズ。                     |
| データベース               | PostgreSQL                        | 堅牢性、豊富な機能セット（JSONB, 全文検索）、高い信頼性。                            |
| ORM/データアクセス         | GORM                              | 開発効率の向上、マイグレーション機能、豊富な関連付け機能。                             |
| データベースマイグレーション | golang-migrate / Atlas            | バージョン管理された安全なスキーマ変更を実現。                                        |
| API認証                    | JWT (Access + Refresh Tokens)     | ステートレス認証、セキュアなセッション管理。                                          |
| リクエスト検証            | go-playground/validator           | 構造体タグベースの宣言的な入力値検証。                                                |
| 設定管理                  | YAML (Viper)                      | 環境ごとの設定を柔軟に管理。                                                           |
| テスト（バックエンド）     | httptest                          | Go標準ライブラリによるHTTPハンドラの単体テスト。                                      |
| テスト（フロントエンド）   | Vitest + Svelte Testing Library   | Viteネイティブの高速なコンポーネントテスト。                                          |
| デプロイメント             | Docker                            | 環境の再現性とポータビリティを確保。                                                  |

ディレクトリ構成
```text
ctflab/
├── cmd/
│   ├── web/                  # Webサーバ用エントリポイント (Gin)
│   │   └── main.go
│   └── desktop/              # Wailsアプリ用エントリポイント
│       └── main.go
├── internal/
│   ├── config/               # 設定読み込み（.env, yaml 等）
│   │   └── config.go
│   ├── domain/               # ドメインモデル (構造体・インターフェース)
│   │   ├── user.go
│   │   └── challenge.go
│   ├── repository/           # DB操作の実装 (GORM)
│   │   ├── user_repo.go
│   │   └── challenge_repo.go
│   ├── service/              # ビジネスロジック
│   │   ├── auth_service.go
│   │   └── challenge_service.go
│   ├── handler/              # Web/デスクトップ共通ハンドラ
│   │   ├── user_handler.go
│   │   └── challenge_handler.go
│   ├── transport/
│   │   ├── http/             # Ginのルーターとミドルウェア
│   │   │   ├── router.go
│   │   │   ├── middleware.go
│   │   │   └── subdomain.go
│   │   └── desktop/          # Wailsのバインディング
│   │       └── bindings.go
├── frontend/                 # SvelteKit（Web/デスクトップ共通UI）
│   ├── src/
│   └── vite.config.ts
├── migrations/               # DBマイグレーションSQL
├── scripts/                  # 初期化・CI/CD用スクリプト
│   └── init_db.sh
├── test/                     # 統合・ユニットテスト
│   └── challenge_test.go
├── docs/                     # API仕様書、設計資料など
│   └── openapi.yaml
├── go.mod
└── README.md
```

開発順番
```
Step 1. 機能一覧と画面構成
Step 2. ER図 / GORMモデル作成
Step 3. GinのAPI設計（router + handler）
Step 4. UIモック（SvelteKitで仮組み）
Step 5. API連携
Step 6. Wails統合
```

# Step 1. 機能一覧と画面構成

## MVP
- 認証機能
    - Github / Microsoft アカウントでログイン (OAuth2)
- ユーザー管理
    - 自分のプロフィール閲覧(表示名、問題数など)
    - ログアウト
- 問題作成・管理
    - 問題新規作成(title,description,カテゴリ,flag,添付ファイル)
    - 問題の保存・編集・削除(自分の問題のみ)
    - 問題にDockerベースの環境を添付(例:pwn用コンテナ)
- 問題公開・共有
    - 演題一覧ページ(自作&他人の問題が見える)
    - 問題詳細ページ(description,添付,flag提出用フォーム)
    - 公開・非公開フラグ
- Flag提出・判定
    - 解答フォーム(flag入力→自動判定)
    - 解答履歴(自分がどれを解いたか)
- 管理者・運用視点の最低限機能
    - サンドボックス環境の自動デプロイ(例: `problem.user.ctflab.dev`)
    - 問題ごとに独立したDocker環境を立てる
    - 運用負荷を減らす問題クリーンアップ(いい定時間後自動停止)