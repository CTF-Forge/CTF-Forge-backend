<script>
  import { user } from '../lib/store.js';

  let email = '';
  let password = '';
  let errorMessage = '';

  async function handleLogin() {
    errorMessage = '';
    try {
      const response = await fetch('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
      });

      if (response.ok) {
        const data = await response.json();
        user.set(data.user);
        // トークンをlocalStorageに保存
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
      } else {
        const errorData = await response.json();
        errorMessage = errorData.error || 'Login failed';
      }
    } catch (error) {
      errorMessage = 'An error occurred during login.';
    }
  }
</script>

<div class="login-container">
  <h2>Login</h2>
  <form on:submit|preventDefault={handleLogin}>
    <div class="form-group">
      <label for="email">Email</label>
      <input type="email" id="email" bind:value={email} required>
    </div>
    <div class="form-group">
      <label for="password">Password</label>
      <input type="password" id="password" bind:value={password} required>
    </div>
    <button type="submit">Login</button>
    {#if errorMessage}
      <p class="error">{errorMessage}</p>
    {/if}
  </form>
</div>

<style>
  .login-container {
    max-width: 400px;
    margin: 2rem auto;
    padding: 2rem;
    border: 1px solid #ccc;
    border-radius: 5px;
  }
  .form-group {
    margin-bottom: 1rem;
  }
  label {
    display: block;
    margin-bottom: 0.5rem;
  }
  input {
    width: 100%;
    padding: 0.5rem;
    font-size: 1rem;
  }
  button {
    width: 100%;
    padding: 0.75rem;
    font-size: 1rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  .error {
    color: red;
    margin-top: 1rem;
  }
</style>
