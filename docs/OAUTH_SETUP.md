# OAuth設定ガイド

## GitHub OAuth設定

### 1. GitHub OAuth App作成

1. GitHubにログイン
2. Settings > Developer settings > OAuth Apps
3. "New OAuth App"をクリック
4. 以下の情報を入力：
   - **Application name**: CTFForge
   - **Homepage URL**: `http://localhost:8080`
   - **Authorization callback URL**: `http://localhost:8080/auth/github/callback`
5. "Register application"をクリック
6. **Client ID**と**Client Secret**をコピー

### 2. 環境変数の設定

`.env`ファイルに以下を追加：

```env
GITHUB_KEY=your_github_client_id
GITHUB_SECRET=your_github_client_secret
GITHUB_CALLBACK=http://localhost:8080/auth/github/callback
```

## Google OAuth設定

### 1. Google Cloud Console設定

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. プロジェクトを作成または選択
3. "APIs & Services" > "Credentials"
4. "Create Credentials" > "OAuth client ID"
5. Application type: "Web application"
6. 以下の情報を入力：
   - **Name**: CTFForge
   - **Authorized JavaScript origins**: `http://localhost:8080`
   - **Authorized redirect URIs**: `http://localhost:8080/auth/google/callback`
7. "Create"をクリック
8. **Client ID**と**Client Secret**をコピー

### 2. 環境変数の設定

`.env`ファイルに以下を追加：

```env
GOOGLE_KEY=your_google_client_id
GOOGLE_SECRET=your_google_client_secret
GOOGLE_CALLBACK=http://localhost:8080/auth/google/callback
```

## テスト手順

### 1. 環境変数の確認

```bash
# .envファイルが存在するか確認
ls -la ctfforge/.env

# 環境変数が読み込まれているか確認
cd ctfforge
go run main.go
```

### 2. サーバー起動

```bash
cd ctfforge
go run main.go
```

### 3. OAuth認証テスト

#### GitHub認証テスト
1. ブラウザで `http://localhost:8080/auth/github` にアクセス
2. GitHubの認証画面が表示される
3. 認証を許可
4. コールバックでJWTトークンが返される

#### Google認証テスト
1. ブラウザで `http://localhost:8080/auth/google` にアクセス
2. Googleの認証画面が表示される
3. 認証を許可
4. コールバックでJWTトークンが返される

### 4. 期待されるレスポンス

認証成功時のレスポンス例：

```json
{
  "message": "oauth authentication successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "username": "your_github_username",
    "email": "your_email@example.com",
    "provider": "github"
  }
}
```

## トラブルシューティング

### よくある問題

1. **"invalid provider"エラー**
   - URLが正しいか確認（`/auth/github` または `/auth/google`）

2. **"oauth authentication failed"エラー**
   - 環境変数が正しく設定されているか確認
   - OAuth Appの設定が正しいか確認

3. **"failed to process oauth user"エラー**
   - データベースが正常に動作しているか確認
   - マイグレーションが実行されているか確認

### デバッグ方法

1. **ログの確認**
   ```bash
   go run main.go
   ```

2. **環境変数の確認**
   ```bash
   cd ctfforge
   go run -c "fmt.Println(os.Getenv(\"GITHUB_KEY\"))" .
   ```

3. **データベースの確認**
   ```sql
   SELECT * FROM oauth_accounts;
   SELECT * FROM users;
   ```

## セキュリティ注意事項

1. **Client Secretの管理**
   - `.env`ファイルをGitにコミットしない
   - 本番環境では環境変数で管理

2. **HTTPSの使用**
   - 本番環境では必ずHTTPSを使用
   - コールバックURLもHTTPSに変更

3. **スコープの設定**
   - 必要最小限のスコープのみを要求
   - ユーザーのプライバシーを尊重 