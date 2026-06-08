// Reactive paginator shared by views over any endpoint that can report whether
// more items remain. The fetch fn returns a `Slice` ({ items, hasMore }), so the
// "has more" rule — total count, short page, cursor — lives in the caller, not
// here. Owns the load / append / in-flight / error state; views only render.

import type { Slice } from './api';

type FetchSlice<T> = (limit: number, offset: number) => Promise<Slice<T>>;

export class Paginator<T> {
  // Slices are whole API responses, only ever reassigned (never mutated in
  // place), so `raw` skips the per-item proxy overhead of deep `$state`.
  items = $state.raw<T[]>([]);
  status = $state<'loading' | 'error' | 'ready'>('loading');
  loadingMore = $state(false);
  // A failed `loadMore` surfaces here instead of flipping `status`, so the
  // already-loaded items stay on screen while the error shows by the button.
  loadMoreError = $state(false);
  hasMore = $state(false);

  #fetch: FetchSlice<T>;
  #limit: number;

  constructor(fetch: FetchSlice<T>, limit = 20) {
    this.#fetch = fetch;
    this.#limit = limit;
  }

  /** Load the first page. Call once from the view's onMount (or an effect). */
  async start() {
    try {
      const slice = await this.#fetch(this.#limit, 0);
      this.items = slice.items;
      this.hasMore = slice.hasMore;
      this.status = 'ready';
    } catch {
      this.status = 'error';
    }
  }

  /** Append the next page; no-op while one is in flight or none remain. */
  async loadMore() {
    if (this.loadingMore || !this.hasMore) return;
    this.loadingMore = true;
    this.loadMoreError = false;
    try {
      const slice = await this.#fetch(this.#limit, this.items.length);
      this.items = [...this.items, ...slice.items];
      this.hasMore = slice.hasMore;
    } catch {
      this.loadMoreError = true;
    } finally {
      this.loadingMore = false;
    }
  }
}
