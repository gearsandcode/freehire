<script lang="ts">
  import { Sun, Moon } from '@lucide/svelte';
  import { themeStore } from '$lib/theme.svelte';

  const isDark = $derived(themeStore.isDark);
</script>

<button
  type="button"
  onclick={() => themeStore.toggle()}
  title={isDark ? 'Dark · click for light' : 'Light · click for dark'}
  aria-label="Toggle theme"
  class="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
>
  {#if isDark}
    <Moon class="size-4" />
  {:else}
    <Sun class="size-4" />
  {/if}
</button>
