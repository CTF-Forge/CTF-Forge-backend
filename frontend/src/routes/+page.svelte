<script>
  import Login from './lib/components/Login.svelte';
  import { user } from './lib/store.js';
</script>

{#if $user}
  <h1>Welcome, {$user.username}!</h1>
  <p>Your user ID is: {$user.id}</p>
{:else}
  <Login />
{/if}