# CTFForge API ドキュメント

## 認証API

### Email認証

#### ユーザー登録
```http
POST /auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**レスポンス**
```json
{
  "message": "user registered successfully"
}
```

#### ログイン
```http
POST /auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

**レスポンス**
```json
{
  "message": "login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

#### トークン更新
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンス**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

#### ログアウト
```http
POST /auth/logout
```

**レスポンス**
```json
{
  "message": "logout successful"
}
```

### OAuth認証

#### OAuth認証開始
```http
GET /auth/{provider}
```

`{provider}` は `github` または `google`

#### OAuth認証コールバック
```http
GET /auth/{provider}/callback
```

**レスポンス**
```json
{
  "message": "oauth authentication successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "username": "testuser",
    "email": "test@example.com",
    "provider": "github"
  }
}
```

#### OAuthトークン更新
```http
POST /auth/oauth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### OAuthログアウト
```http
POST /auth/oauth/logout
```

## 保護されたAPI

### ユーザー情報取得
```http
GET /api/me
Authorization: Bearer {access_token}
```

**レスポンス**
```json
{
  "user_id": 1,
  "username": "testuser"
}
```

## 公開API

### 問題一覧取得
```http
GET /api/public/challenges
```

**認証なしの場合**
```json
{
  "message": "public challenges"
}
```

**認証ありの場合**
```json
{
  "message": "public challenges",
  "user_id": 1
}
```

## エラーレスポンス

### 400 Bad Request
```json
{
  "error": "invalid request"
}
```

### 401 Unauthorized
```json
{
  "error": "authentication failed"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

## 認証フロー

### Email認証フロー
1. ユーザー登録 (`POST /auth/register`)
2. ログイン (`POST /auth/login`)
3. アクセストークンを使用してAPIにアクセス
4. トークン期限切れ時は更新 (`POST /auth/refresh`)

### OAuth認証フロー
1. OAuth認証開始 (`GET /auth/{provider}`)
2. プロバイダー側で認証
3. コールバック処理 (`GET /auth/{provider}/callback`)
4. アクセストークンを使用してAPIにアクセス
5. トークン期限切れ時は更新 (`POST /auth/oauth/refresh`)

## 環境変数設定

必要な環境変数を設定してください：

```env
# JWT設定
JWT_ACCESS_SECRET=your_access_secret_key
JWT_REFRESH_SECRET=your_refresh_secret_key
JWT_ISSUER=ctfforge
JWT_ACCESS_EXPIRE_HOURS=1
JWT_REFRESH_EXPIRE_HOURS=168

# セッション設定
SESSION_SECRET=your_session_secret

# OAuth設定
GITHUB_KEY=your_github_client_id
GITHUB_SECRET=your_github_client_secret
GITHUB_CALLBACK=http://localhost:8080/auth/github/callback

GOOGLE_KEY=your_google_client_id
GOOGLE_SECRET=your_google_client_secret
GOOGLE_CALLBACK=http://localhost:8080/auth/google/callback
``` 