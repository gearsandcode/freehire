<script lang="ts">
  import { ExternalLink } from '@lucide/svelte';
  import { getJob } from '$lib/api';
  import type { Job } from '$lib/types';
  import { Badge, Button } from '$lib/ui';
  import { formatDate } from '$lib/utils';
  import States from './States.svelte';

  let { id }: { id: string } = $props();

  let job = $state.raw<Job | null>(null);
  let status = $state<'loading' | 'error' | 'ready'>('loading');

  // Reload whenever the route id changes.
  $effect(() => {
    const current = id;
    status = 'loading';
    job = null;
    getJob(current)
      .then((j) => {
        if (current !== id) return;
        job = j;
        status = 'ready';
      })
      .catch(() => {
        if (current !== id) return;
        status = 'error';
      });
  });
</script>

{#if status === 'loading'}
  <States state="loading" rows={3} />
{:else if status === 'error' || !job}
  <States state="error" message="Job not found." />
{:else}
  {@const posted = formatDate(job.posted_at)}
  <article class="flex flex-col gap-4">
    <div>
      <h1 class="text-xl font-semibold tracking-tight">{job.title}</h1>
      <p class="mt-1 text-sm text-muted-foreground">
        {#if job.company_slug}
          <a href={`/companies/${job.company_slug}`} class="hover:text-foreground hover:underline">
            {job.company || 'Unknown company'}
          </a>
        {:else}
          {job.company || 'Unknown company'}
        {/if}
        {#if job.location}· {job.location}{/if}
      </p>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      {#if job.remote}<Badge variant="secondary">Remote</Badge>{/if}
      <Badge variant="outline">{job.source}</Badge>
      {#if posted}<span class="text-xs text-muted-foreground">Posted {posted}</span>{/if}
    </div>

    {#if job.description}
      <p class="whitespace-pre-wrap text-sm leading-relaxed">{job.description}</p>
    {/if}

    <div>
      <Button variant="primary" href={job.url} target="_blank" rel="noopener noreferrer">
        Apply <ExternalLink class="size-4" />
      </Button>
    </div>
  </article>
{/if}
