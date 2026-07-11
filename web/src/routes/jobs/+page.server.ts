import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// The job feed moved to the homepage. `/jobs` is a permanent (301) redirect that
// carries the query string through verbatim, so every saved filter/share link
// (e.g. /jobs?q=go&remote=true) lands on the same filtered feed at `/`. Child
// routes (/jobs/[slug], /jobs/swipe) have their own loads and are unaffected.
export const load: PageServerLoad = ({ url }) => {
  redirect(301, `/${url.search}`);
};
