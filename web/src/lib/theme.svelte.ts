// Theme controller. Three modes — explicit `light` / `dark`, or `system` which
// tracks the OS preference. Persisted in localStorage under `hire.theme`.
// `main.ts` calls `initTheme()` once at boot; components read `themeStore` and
// call `setMode(...)`.

const STORAGE_KEY = 'hire.theme';

export type ThemeMode = 'light' | 'dark' | 'system';

const mq = window.matchMedia('(prefers-color-scheme: dark)');

function readStored(): ThemeMode {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (raw === 'light' || raw === 'dark' || raw === 'system') return raw;
  return 'system';
}

function apply(mode: ThemeMode) {
  const dark = mode === 'dark' || (mode === 'system' && mq.matches);
  document.documentElement.classList.toggle('dark', dark);
}

class ThemeStore {
  mode = $state<ThemeMode>(readStored());
  /** Live OS `prefers-color-scheme: dark` state, kept current by `initTheme`. */
  systemDark = $state(mq.matches);

  /** Effective dark state — explicit `dark`, or `system` resolving to the OS. */
  isDark = $derived(this.mode === 'dark' || (this.mode === 'system' && this.systemDark));

  setMode(next: ThemeMode) {
    this.mode = next;
    try {
      localStorage.setItem(STORAGE_KEY, next);
    } catch {
      // best-effort: private mode / quota
    }
    apply(next);
  }

  /** Binary flip: light <-> dark. `system` collapses to its effective value
   *  first, so the first click always yields a concrete choice. */
  toggle() {
    this.setMode(this.isDark ? 'light' : 'dark');
  }
}

export const themeStore = new ThemeStore();

/** Apply the stored theme and keep `system` mode tracking the OS preference. */
export function initTheme() {
  apply(themeStore.mode);
  const onChange = () => {
    themeStore.systemDark = mq.matches;
    if (themeStore.mode === 'system') apply('system');
  };
  mq.addEventListener('change', onChange);
  if (import.meta.hot) {
    import.meta.hot.dispose(() => mq.removeEventListener('change', onChange));
  }
}
