# JWT Token Package

CTFForgeのJWT認証機能を提供するパッケージです。

## 機能

- アクセストークンとリフレッシュトークンの生成・検証
- 詳細なエラーハンドリング
- Ginミドルウェア
- トークンの期限切れ判定

## 使用方法

### 1. JWTマネージャーの初期化

```go
import (
    "time"
    "github.com/Saku0512/CTFForge/ctfforge/config"
    "github.com/Saku0512/CTFForge/ctfforge/pkg/token"
)

// 設定からJWTマネージャーを作成
jwtManager := token.NewJWTManager(
    config.GetJWTAccessSecret(),
    config.GetJWTRefreshSecret(),
    config.GetJWTIssuer(),
    config.GetJWTAccessExpireDuration(),
    config.GetJWTRefreshExpireDuration(),
)
```

### 2. トークンペアの生成

```go
// ユーザー情報からトークンペアを生成
tokenPair, err := jwtManager.GenerateTokenPair(userID, username)
if err != nil {
    // エラーハンドリング
}

// レスポンス例
c.JSON(http.StatusOK, gin.H{
    "access_token":  tokenPair.AccessToken,
    "refresh_token": tokenPair.RefreshToken,
    "expires_in":    tokenPair.ExpiresIn,
})
```

### 3. トークンの検証

```go
// アクセストークンの検証
claims, err := jwtManager.VerifyAccessToken(tokenString)
if err != nil {
    switch err {
    case token.ErrExpiredToken:
        // トークン期限切れ
    case token.ErrInvalidToken:
        // 無効なトークン
    default:
        // その他のエラー
    }
}

// ユーザー情報の取得
userID := claims.UserID
username := claims.Username
```

### 4. リフレッシュトークンでの更新

```go
// リフレッシュトークンを使用して新しいトークンペアを生成
newTokenPair, err := jwtManager.RefreshTokenPair(refreshToken)
if err != nil {
    // エラーハンドリング
}
```

### 5. Ginミドルウェアの使用

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/Saku0512/CTFForge/ctfforge/pkg/token"
)

// 認証必須のミドルウェア
router.Use(token.AuthMiddleware(jwtManager))

// オプショナル認証のミドルウェア
router.Use(token.OptionalAuthMiddleware(jwtManager))

// ハンドラー内でユーザー情報を取得
func ProtectedHandler(c *gin.Context) {
    userID, exists := token.GetUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
        return
    }
    
    username, _ := token.GetUsername(c)
    // ユーザー情報を使用した処理
}
```

## 環境変数

以下の環境変数を設定してください：

```env
# JWT設定
JWT_ACCESS_SECRET=your_access_secret_key
JWT_REFRESH_SECRET=your_refresh_secret_key
JWT_ISSUER=ctfforge
JWT_ACCESS_EXPIRE_HOURS=1
JWT_REFRESH_EXPIRE_HOURS=168

# 後方互換性のため、以下の設定も使用可能
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24
```

## エラー型

```go
var (
    ErrInvalidToken     = errors.New("invalid token")
    ErrExpiredToken     = errors.New("token has expired")
    ErrInvalidSignature = errors.New("invalid token signature")
    ErrInvalidClaims    = errors.New("invalid token claims")
)
```

## セキュリティ考慮事項

1. **シークレットキー**: 強力なシークレットキーを使用し、環境変数で管理
2. **トークン期限**: アクセストークンは短期限（1時間）、リフレッシュトークンは長期限（7日）
3. **HTTPS**: 本番環境では必ずHTTPSを使用
4. **トークン保存**: クライアント側ではセキュアな方法でトークンを保存 