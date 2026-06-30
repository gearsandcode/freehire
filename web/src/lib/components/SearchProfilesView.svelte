<script lang="ts">
  import { ApiError, facetCounts } from '$lib/api';
  import { isAuthenticated } from '$lib/auth.svelte';
  import { openAuthDialog } from '$lib/auth-dialog.svelte';
  import { CATEGORY_OPTIONS, categoryLabel, type FacetOption } from '$lib/facets';
  import { searchProfiles } from '$lib/searchProfiles.svelte';
  import type { SearchProfile } from '$lib/types';
  import { Button, Input } from '$lib/ui';
  import SearchSelect from './facets/SearchSelect.svelte';
  import States from './States.svelte';

  let status = $state<'loading' | 'error' | 'ready'>('loading');
  const profiles = $derived(searchProfiles.items);

  // The universe of skills (canonical tokens with job counts) for the picker, fetched
  // from the facet-distribution endpoint — the same source the filter panel's dynamic
  // skills facet uses. Empty params = the whole catalogue's skill distribution.
  let skillDist = $state.raw<FacetOption[]>([]);

  // Form doubles as create (editingId null) and edit (editingId set). specialization is
  // single-valued; skills is a set. A profile needs a name, a specialization, and at
  // least one skill — the same invariants the server enforces.
  let editingId = $state<number | null>(null);
  let name = $state('');
  let specialization = $state('');
  let skills = $state.raw<string[]>([]);
  let formError = $state<string | null>(null);
  let busy = $state(false);

  const canSubmit = $derived(name.trim() !== '' && specialization !== '' && skills.length > 0);

  // The picker options: the distribution plus any selected skills not in it (so an
  // edited profile's skills stay visible and removable even if they have no live count).
  const skillOptions = $derived.by((): FacetOption[] => {
    const known = new Set(skillDist.map((o) => o.value));
    const extra = skills.filter((s) => !known.has(s)).map((s) => ({ value: s, label: s }));
    return [...skillDist, ...extra];
  });

  async function loadProfiles() {
    status = 'loading';
    try {
      await searchProfiles.ensureLoaded();
      status = 'ready';
    } catch {
      status = 'error';
    }
  }

  async function loadSkills() {
    try {
      const counts = await facetCounts(new URLSearchParams());
      const dist = counts.facets?.skills ?? {};
      skillDist = Object.entries(dist)
        .map(([value, count]) => ({ value, label: value, count }))
        .toSorted((a, b) => b.count - a.count || a.label.localeCompare(b.label));
    } catch {
      // best-effort: an empty distribution still lets the user type-and-pick nothing,
      // but selected skills (on edit) remain visible via skillOptions.
    }
  }

  // Load once the session is confirmed (the boot-time /me resolution may still be in
  // flight when the page is opened directly), mirroring ApiKeysView. Reset the per-user
  // cache on sign-out so a different user does not see the previous one's profiles.
  $effect(() => {
    if (isAuthenticated()) {
      void loadProfiles();
      void loadSkills();
    } else {
      searchProfiles.reset();
    }
  });

  function resetForm() {
    editingId = null;
    name = '';
    specialization = '';
    skills = [];
    formError = null;
  }

  function startEdit(p: SearchProfile) {
    editingId = p.id;
    name = p.name;
    specialization = p.specialization;
    skills = [...p.skills];
    formError = null;
  }

  // Single-select: clicking the chosen specialization clears it, any other replaces it.
  function toggleSpecialization(value: string) {
    specialization = specialization === value ? '' : value;
  }

  function toggleSkill(value: string) {
    skills = skills.includes(value) ? skills.filter((s) => s !== value) : [...skills, value];
  }

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!canSubmit || busy) return;
    busy = true;
    formError = null;
    try {
      const patch = { name: name.trim(), specialization, skills };
      if (editingId === null) {
        await searchProfiles.create(patch.name, specialization, skills);
      } else {
        await searchProfiles.update(editingId, patch);
      }
      resetForm();
    } catch (err) {
      formError = err instanceof ApiError ? err.message : 'Could not save the profile. Please try again.';
    } finally {
      busy = false;
    }
  }

  async function remove(p: SearchProfile) {
    if (!window.confirm(`Delete profile “${p.name}”?`)) return;
    try {
      await searchProfiles.remove(p.id);
      if (editingId === p.id) resetForm();
    } catch {
      formError = 'Could not delete the profile. Please try again.';
    }
  }
</script>

{#if !isAuthenticated()}
  <div class="flex flex-col items-center gap-3 py-12 text-center">
    <p class="text-sm text-muted-foreground">Sign in to create search profiles.</p>
    <Button variant="primary" onclick={() => openAuthDialog()}>Sign in</Button>
  </div>
{:else}
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-1">
      <h1 class="text-2xl font-semibold tracking-tight">Search profiles</h1>
      <p class="text-sm text-muted-foreground">
        Describe what you do — a specialization and your skills — and reuse it to find relevant work.
      </p>
    </div>

    <form onsubmit={submit} class="flex flex-col gap-4 rounded-lg border border-border p-4">
      <p class="text-sm font-medium">{editingId === null ? 'New profile' : 'Edit profile'}</p>

      <label class="flex flex-col gap-1">
        <span class="text-sm font-medium">Name</span>
        <Input bind:value={name} placeholder="e.g. Go backend" maxlength={100} class="w-full" />
      </label>

      <div class="flex flex-col gap-1.5">
        <span class="text-sm font-medium">Specialization</span>
        <SearchSelect
          options={CATEGORY_OPTIONS}
          selected={specialization ? [specialization] : []}
          placeholder="Search specializations"
          onToggle={toggleSpecialization}
        />
      </div>

      <div class="flex flex-col gap-1.5">
        <span class="text-sm font-medium">Skills</span>
        <SearchSelect
          options={skillOptions}
          selected={skills}
          placeholder="Search skills"
          onToggle={toggleSkill}
        />
      </div>

      {#if formError}
        <p class="text-sm text-destructive">{formError}</p>
      {/if}

      <div class="flex items-center gap-2">
        <Button variant="primary" type="submit" disabled={!canSubmit || busy}>
          {busy ? 'Saving…' : editingId === null ? 'Create profile' : 'Save changes'}
        </Button>
        {#if editingId !== null}
          <Button variant="ghost" type="button" onclick={resetForm}>Cancel</Button>
        {/if}
      </div>
    </form>

    {#if status === 'loading'}
      <States state="loading" />
    {:else if status === 'error'}
      <States state="error" message="Couldn't load your profiles." />
    {:else if profiles.length === 0}
      <States state="empty" message="No profiles yet. Create one above." />
    {:else}
      <ul class="flex flex-col divide-y divide-border rounded-lg border border-border">
        {#each profiles as profile (profile.id)}
          <li class="flex items-start justify-between gap-3 px-4 py-3">
            <div class="flex min-w-0 flex-col gap-1">
              <span class="truncate text-sm font-medium">{profile.name}</span>
              <span class="text-xs text-muted-foreground">{categoryLabel(profile.specialization)}</span>
              <div class="flex flex-wrap gap-1">
                {#each profile.skills as skill (skill)}
                  <span class="rounded bg-secondary px-1.5 py-0.5 text-xs text-secondary-foreground">{skill}</span>
                {/each}
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-1">
              <Button variant="ghost" size="sm" onclick={() => startEdit(profile)}>Edit</Button>
              <Button variant="ghost" size="sm" onclick={() => remove(profile)}>Delete</Button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
{/if}
