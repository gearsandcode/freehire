<script lang="ts">
  import { router } from '$lib/router.svelte';
  import { cn } from '$lib/utils';
  import ThemeToggle from './ThemeToggle.svelte';

  const name = $derived(router.route.name);

  const links = [
    { href: '/', label: 'Jobs', match: ['jobs', 'job'] },
    { href: '/companies', label: 'Companies', match: ['companies', 'company'] },
  ];
</script>

<header class="border-b border-border">
  <div class="mx-auto flex h-14 max-w-3xl items-center gap-6 px-4">
    <a href="/" class="text-sm font-semibold tracking-tight">hire</a>

    <nav class="flex items-center gap-4 text-sm">
      {#each links as link (link.href)}
        <a
          href={link.href}
          class={cn(
            'transition-colors hover:text-foreground',
            link.match.includes(name) ? 'text-foreground' : 'text-muted-foreground',
          )}
        >
          {link.label}
        </a>
      {/each}
    </nav>

    <div class="ml-auto">
      <ThemeToggle />
    </div>
  </div>
</header>
