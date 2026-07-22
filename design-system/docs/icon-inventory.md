# Icon inventory

Living doc. Every icon in `web/`, classified **library** (from `@lucide/svelte`) or **bespoke** (custom inline SVG / brand mark), with a reason for each bespoke one.

Last updated: phase 1 investigation. Totals: **172 library usages** across 48 files + **13 bespoke SVG files** (covering ~15 distinct SVGs).

## Consolidation candidates

1. **`PipelineFunnel.svelte` + `HomeFunnel.svelte`** — near-identical Sankey funnels (~107 lines each, same geometry constants, same ribbon-path formula). `HomeFunnel`'s header comment already says: *"Seam: when it lands, the two can fold into one shared presentational component."* Extract a shared `SankeyFunnel.svelte` taking a `buckets: {key,label,color}[]` prop.
2. **`BoardCard.svelte` "has notes" glyph, `DocsNav.svelte` search glyph, `docs/api/+layout.svelte` hamburger + chevron** — four inline UI icons with direct Lucide equivalents already used elsewhere (`StickyNote`/`AlignLeft`, `Search`, `Menu`, `ChevronDown`). Replace with `@lucide/svelte` imports — pure consistency cleanup.
3. **`BrandMark.svelte` + `web/static/favicon.svg`** — identical path data maintained in two places. Static file is forced (browsers fetch `/favicon.svg`), but the path string could be hoisted to a shared `web/src/lib/brandPath.ts` constant.
4. **`NoteEditor.svelte` 7 toolbar glyphs** — already Lucide paths inlined as strings (EasyMDE takes HTML strings, not Svelte components). Leave as-is; the constraint is real and documented.

## Library — `@lucide/svelte`

172 usages across 48 files. Below: one row per file with its named imports.

| file | lucide imports |
|---|---|
| routes/jobs/[slug]/fit/+page.svelte | ArrowLeft, SquarePen |
| routes/my/+layout.svelte | User, Bot, LayoutList, Activity, Bell, Key, FileText, ScrollText, Inbox, Link2, PanelLeftClose, PanelLeft (+ type LucideIcon) |
| routes/my/profile/+page.svelte | ScanSearch, Trash2 |
| routes/my/assistant/+page.svelte | Bot, AlertTriangle, Terminal, ChevronRight, ArrowUp, Loader2, Trash2, Plus, RefreshCw, ExternalLink, X, FileText, PanelLeft |
| lib/components/JobFitAnalysis.svelte | ArrowRight, FileText, ScanSearch |
| lib/components/facets/SearchSelect.svelte | Plus |
| lib/components/VerdictView.svelte | Check, TrendingUp |
| lib/components/onboarding/OnboardingWizard.svelte | ArrowLeft, ArrowRight, FileUp, LoaderCircle, X |
| lib/components/CompanyLogo.svelte | Globe |
| lib/components/ApplicationLinkPicker.svelte | Search |
| lib/components/onboarding/OnboardingBanner.svelte | Target, X |
| lib/components/facets/FacetSection.svelte | X |
| lib/components/CompanyFollowButton.svelte | Bell |
| lib/components/facets/RemoteSearchSelect.svelte | X |
| lib/components/AnalysesView.svelte | ArrowRight |
| lib/components/onboarding/OnboardingAlertBanner.svelte | Bookmark, X |
| lib/components/JobView.svelte | ArrowRight, Bookmark, Check, CheckCircle2, Eye, Flag |
| lib/components/ATSReportView.svelte | Check, Copy, TriangleAlert, X |
| lib/components/JobRow.svelte | Bookmark |
| lib/components/HeaderSearch.svelte | Search, X |
| lib/components/filters/LocationPane.svelte | ChevronDown, Search, X |
| lib/components/SavedSearchesView.svelte | Check, Mail |
| lib/components/filters/FilterModal.svelte | Bell, UserRound |
| lib/components/HeaderLocationFilter.svelte | ChevronDown, Globe, House, Blend, Building |
| lib/components/ListToolbar.svelte | Layers, SlidersHorizontal |
| lib/components/filters/FilterModalShell.svelte | ArrowRight, LoaderCircle, X |
| lib/components/InboxView.svelte | Mail, AtSign, Copy, Search, RefreshCw, ChevronLeft, CheckCheck, Trash2 |
| lib/components/SwipeDeck.svelte | Heart, RotateCcw, SlidersHorizontal, X |
| lib/components/GmailConnectDialog.svelte | Mail, X, Lock, ExternalLink |
| lib/components/FilterEdgeTab.svelte | SlidersHorizontal |
| lib/components/filters/SaveSearchAlert.svelte | Bell, Bookmark, Check |
| lib/components/ReportDialog.svelte | ArrowLeft, Ban, BellOff, Check, ChevronRight, Clock, MoreHorizontal, Send, ShieldAlert, X |
| lib/components/HeaderMenu.svelte | Menu, X, Sun, Moon, Briefcase, Building2, CircleUser, Activity, ListChecks, BellRing, KeyRound, Inbox, FileText, SquarePlus, ShieldCheck, Layers, ChartColumn, TrendingUp, Info, LogOut, LogIn |
| lib/components/ProfileForm.svelte | ArrowUp, Check, X |
| lib/components/SavedSearches.svelte | Pencil, Trash2 |
| lib/components/filters/FacetHeader.svelte | X |
| lib/components/BoardCard.svelte | Mail |
| lib/components/ResumeStructuredView.svelte | Briefcase, GraduationCap, Languages, Link (as LinkIcon), Mail, MapPin, Phone |
| lib/components/JobFitFull.svelte | RefreshCw, FileText, Check, Loader, TriangleAlert |
| lib/components/HeaderListSearch.svelte | Search, X |
| lib/components/filters/FilterSummaryShell.svelte | SlidersHorizontal, X |
| lib/components/JobDrawer.svelte | Trash2, X, ExternalLink |
| lib/components/filters/CategoryPane.svelte | ChevronDown, Search |
| lib/components/cv/StringListEditor.svelte | X, Plus |
| lib/components/cv/CvList.svelte | FileText, Download, Trash2, Plus |
| lib/components/JobMatch.svelte | FileText, Lock |
| lib/components/JobsView.svelte | Layers |
| lib/components/cv/CvEditor.svelte | ArrowLeft, Download, Plus, Trash2 |

## Bespoke — inline `<svg>`

| file | icon | reason |
|---|---|---|
| BrandMark.svelte | freehire ring + diamond mark (currentColor, 512 viewBox) | Brand mark. Lucide has no freehire logo. currentColor tracks theme alongside wordmark. |
| ProviderIcon.svelte | Google "G" (4-color) | Brand logo with official colors. Lucide ships no brand marks. |
| ProviderIcon.svelte | GitHub Octocat | Brand logo. Lucide ships no brand marks. |
| ProviderIcon.svelte | Telegram paper plane | Brand logo. Same reason. |
| ProviderIcon.svelte | LinkedIn "in" | Brand logo. Same reason. |
| CompanyLogo.svelte | Monogram tile (rect + hashed-hue fill + text initial) | Generated glyph — fill color hashed from company name. Not an icon; a fallback tile. (Last-resort branch uses Lucide Globe.) |
| RateDonut.svelte | Donut chart (two circles with stroke-dasharray + center text) | Data-driven chart — geometry is a function of percent prop. No static icon equivalent. |
| GrowthArea.svelte | Cumulative area chart (polygon + polyline + axis text) | Data-driven chart scaled to point series. |
| ActivityBars.svelte | Grouped bar chart (per-period rect pairs + tooltip) | Data-driven chart with hover model. |
| PipelineFunnel.svelte | Single-level Sankey (source rect + ribbon path + node rect + text) | Data-driven chart; ribbons are computed cubic paths. |
| HomeFunnel.svelte | Single-level Sankey (same shape as PipelineFunnel) | Data-driven chart. Near-duplicate of PipelineFunnel — consolidation candidate. |
| JobFitFull.svelte | Radial gauge (two circles with stroke-dasharray driven by overall_score) | Data-driven gauge — arc length is a function of score. |
| BoardCard.svelte | "Has notes" glyph (three horizontal lines) | UI glyph. Replaceable by Lucide StickyNote/AlignLeft. Consolidation candidate. |
| DocsNav.svelte | Search glyph (circle + line) | UI glyph. Directly replaceable by Lucide Search. Consolidation candidate. |
| routes/docs/api/+layout.svelte | Hamburger (three horizontal lines) | UI glyph. Directly replaceable by Lucide Menu. Consolidation candidate. |
| routes/docs/api/+layout.svelte | Chevron down | UI glyph. Directly replaceable by Lucide ChevronDown. Consolidation candidate. |
| NoteEditor.svelte | 7 EasyMDE toolbar glyphs (bold, italic, heading, lists, link, preview) | Forced pattern: EasyMDE takes HTML string, not Svelte component. Already Lucide paths inlined as strings. Leave as-is. |
| static/favicon.svg | Ring + diamond (same path as BrandMark.svelte + prefers-color-scheme) | Static asset — must be standalone file. Cannot import Svelte component. |
| static/favicon.png | PNG rendering of brand mark | Required PNG format for favicon. |
| static/apple-touch-icon.png | PNG rendering of brand mark | Required PNG format for apple-touch-icon. |
