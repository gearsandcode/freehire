# Component inventory

Living doc. One row per `.svelte` component under `web/src/lib/components/` (root + `cv/`, `facets/`, `filters/`, `onboarding/`). Classified **systemized** (generic, reusable primitive or pattern; will be built from the shared primitive layer) or **bespoke** (one-off, page-specific, or content-bound; stays as-is with a stated reason).

Last updated: phase 1 investigation. Total: **99 components** (23 systemized, 76 bespoke).

## Migration-candidate patterns

These are the recurring shapes the design-system primitive layer (phase 04) should absorb. Listed by highest consolidation value.

1. **Dialogs/modals — 4 implementations.** `AuthDialog`, `ReportDialog`, `GmailConnectDialog` each hand-roll backdrop + Escape + scroll-lock. `FilterModalShell` is already a clean reusable shell. The three one-offs fold into a shared `<Dialog>` primitive.
2. **Hand-built SVG charts — 6 implementations.** `ActivityBars`, `GrowthArea`, `HomeFunnel`, `PipelineFunnel`, `RateDonut`, `FacetBreakdown`. `RateDonut` is the only generic one; the rest are bound to specific data shapes.
3. **Empty/loading states — scattered.** `States` is the shared primitive, but several charts/views inline their own empty message (`GrowthArea`, `ActivityBars`, `HomeFunnel`). Consolidate onto `States`.
4. **Searchable selects / dropdowns — 5 variants.** `SearchSelect`, `RemoteSearchSelect`, `TokenInput`, `ApplicationLinkPicker`, `HeaderSearch` results dropdown. Overlapping "type → filter → pick" pattern with stale-request guards.
5. **Tabs — 5 inline implementations.** `JobRelated`, `JobDrawer`, `ModerationView`, `VerdictView`, `ProfileForm` each hand-roll a tab strip with bespoke active styling. A shared `<Tabs>` primitive absorbs all.
6. **List-row + paginator composition — 3 near-identical.** `SavedJobs`, `JobHistory`, and the pattern inside `JobsView`/`CompaniesView` repeat `Paginator` + `States` + `JobRow` list + `LoadMore`. A `<PaginatedList>` wrapper dedupes.
7. **Marketing "steps" sections — 5 copies.** `AboutValues`, `HomeView` (`sourced`), `ContributeLandingView`, `ForCompaniesView`, `RecruitersView` each render a numbered `01/02/03` definition-list with the same markup. Candidate `<NumberedSteps>` block.

## Systemized (23)

| file | subdir | reason |
|---|---|---|
| Avatar.svelte | | Generic email-initial avatar circle with deterministic color. Reusable anywhere a user face is shown. |
| ChipFacet.svelte | filters | Generic chip-facet wrapper (FacetHeader + PillGroup) driven by a facet param. Already flagged as migration candidate. |
| CompanyLogo.svelte | | Logo-with-monogram-fallback primitive reused across job row, job page, company page, header search, board card. |
| DocsCodeBlock.svelte | | Generic copyable code block (label + pre/Shiki + copy button). No doc-specific data coupling. |
| FacetHeader.svelte | filters | Reusable label + Clear header shared by ChipFacet, CategoryPane, FacetSection. |
| FacetSection.svelte | facets | One facet section (header + control dispatch) reused across job modal, company modal, sidebar. |
| FilterEdgeTab.svelte | | Generic floating edge-tab button reused by /jobs, /companies, account profile. |
| FilterModalShell.svelte | filters | Domain-agnostic two-pane filter-modal chrome (backdrop/rail/footer/deferred-apply) shared by FilterModal and CompanyFilterModal. |
| FilterSummaryShell.svelte | filters | Reusable filter-summary sidebar shell shared by FilterSummary and CompanyFilterSummary. |
| InfiniteScroll.svelte | | Pure IntersectionObserver sentinel trigger reused across jobs and companies lists. |
| InsightsPageShell.svelte | | Shared page chrome (breadcrumb, H1, intro, rail) reused by all salary/skills/roles insight pages. |
| ListToolbar.svelte | | Generic mobile list toolbar + scroll-revealed floating tabs reused by JobsView and CompaniesView. |
| LoadMore.svelte | | Generic "Load more" button + error line. Already flagged as migration candidate. |
| PillGroup.svelte | facets | Stateless three-state pill group reused by ChipFacet, FacetSection, CategoryPane. |
| ProviderIcon.svelte | | Brand SVG icon set (Google/GitHub/Telegram/LinkedIn) reused by AuthDialog, Footer, GithubStars. |
| RateDonut.svelte | | Purely presentational donut chart with generic percent/label props. No data-source coupling. |
| RemoteSearchSelect.svelte | facets | Generic server-backed searchable multi-select with chips. Reusable across entity facets. |
| SaveSearchAlert.svelte | filters | Centralized save-search + Telegram-alert control reused by sidebar, onboarding banner, modal, account page. |
| SearchSelect.svelte | facets | Generic searchable multi-select of three-state pills. Reused by job modal, sidebar, ProfileForm, OnboardingWizard. |
| Seo.svelte | | Generic per-page `<svelte:head>` metadata primitive used on every route. |
| States.svelte | | Shared loading/empty/error rendering used by nearly every data view. Already flagged as migration candidate. |
| StringListEditor.svelte | cv | Generic add/remove list-of-strings editor over Input + Button. Reusable for any token list. |
| TokenInput.svelte | facets | Generic free-text chip input (Enter add / Backspace remove) for open-vocabulary facets. |

## Bespoke (76)

| file | subdir | reason |
|---|---|---|
| ATSReportView.svelte | | Renders backend `ATSReport` shape with five fixed weighted categories and AI suggestions. |
| AboutValues.svelte | | Hardcoded marketing values block specific to /about. |
| ActivityBars.svelte | | Hand-built grouped bar chart bound to `ActivityPoint` catalogue-flow series. |
| AnalysesView.svelte | | Tracking "AI fit" tab listing `MyAnalysisItem` rows tied to analyses API. |
| AnalyticsView.svelte | | /analytics page wiring FilterModal + FacetBreakdown to facet-counts endpoint. |
| ApiKeysView.svelte | | /my/api-keys page — create form, one-time reveal, key list tied to API-key resource. |
| ApplicationLinkPicker.svelte | | One-off searchable popover used only by InboxView to link email to tracked application. |
| AuthDialog.svelte | | Modal hand-rolled for auth form. Bound to login/register/OAuth logic. |
| BoardCard.svelte | | Tracking-board card bound to `MyJob` + stage + email-count semantics. |
| BoardColumn.svelte | | Drag-and-drop column bound to `BoardColumnId` and svelte-dnd-action. Only used by JobBoard. |
| BrandMark.svelte | | freehire brand mark SVG. Brand asset, not a generic primitive. |
| CategoryPane.svelte | filters | Specialization pane tied to `CATEGORY_GROUPS` and `category` facet. |
| ChatGptView.svelte | | Marketing page for ChatGPT GPT with hardcoded hero, copy, mock chat. |
| CliView.svelte | | Marketing page for CLI with hardcoded install one-liner and command reference. |
| CompaniesView.svelte | | /companies list page composing its own paginator, filters, header-scope wiring. |
| CompanyAbout.svelte | | "About" card bound to `Company.company_info.description` field. |
| CompanyFacts.svelte | | Facts definition-list bound to `Company` scalar fields and curated YC badges. |
| CompanyFollowButton.svelte | | Subscribe button wired to saved-search + Telegram-notification stores for one company. |
| CompanyFilterModal.svelte | filters | Company-facets wrapper over FilterModalShell. Bound to `COMPANY_FACETS` grouping. |
| CompanyFilterSummary.svelte | filters | Company-facets summary bound to `CompanyFilterStore` and `COMPANY_FACETS`. |
| CompanyHeader.svelte | | Company identity header card bound to `Company` entity and follow button. |
| CompanyView.svelte | | Company detail page composition (header + facts + streamed JobsView). |
| ContributeLandingView.svelte | | Marketing/landing page for contributions with hardcoded steps and ATS list. |
| ContributeView.svelte | | /contribute form page wired to submission API and reward points. |
| DocsEndpoint.svelte | | Renders one API endpoint from docs `Endpoint` spec with method/path/params. |
| DocsNav.svelte | | Docs navigation rail with scroll-spy tied to `CONCEPTS`/`NAV` doc registries. |
| FacetBreakdown.svelte | | Bar chart bound to `FacetDef` distribution and `FilterStore` drill-down. |
| FilterModal.svelte | filters | Job-filters wrapper over FilterModalShell. Bound to job `RAIL` and `StagedFilters`. |
| FilterSummary.svelte | filters | Job-filters sidebar bound to `FilterStore`, `FACETS`, freshness/salary controls. |
| Footer.svelte | | Site footer with hardcoded navigation groups and social links. |
| ForCompaniesView.svelte | | Marketing page for companies with hardcoded benefits/freshness copy. |
| GithubStars.svelte | | Star-count badge bound to `githubStars` store and freehire repo URL. |
| GmailConnectDialog.svelte | | One-off modal with hardcoded pipeline-step marketing copy for Gmail connect. |
| GrowthArea.svelte | | Hand-built area chart bound to `UserGrowthPoint` member-growth series. |
| HeaderListSearch.svelte | | Header text input wired to `listSearchTarget` store and `/` / Cmd-K hotkeys. |
| HeaderLocationFilter.svelte | | Header location popover tied to `JOBS_SCOPE`/`COMPANIES_SCOPE` and LocationPane. |
| HeaderMenu.svelte | | Site nav + account menu + theme toggle with hardcoded primary/account link sets. |
| HeaderSearch.svelte | | Global launcher dropdown wired to jobs + companies search endpoints. |
| HomeFunnel.svelte | | Decorative Sankey chart for homepage with inlined bucket vocabulary. |
| HomeView.svelte | | Homepage: hardcoded hero, sourced values, illustrative feed marquee, FAQ. |
| InboxView.svelte | | /inbox page: Gmail/mailbox triage with own tabs, search, link picker. |
| JobBoard.svelte | | Tracking board page owning columns, deep-link drawer, dnd state. |
| JobDrawer.svelte | | Tracking application drawer with bespoke tabs (application/fit/description/emails). |
| JobFitAnalysis.svelte | | Compact fit-summary block bound to cached `JobFitResponse` and quota. |
| JobFitFull.svelte | | Full AI-fit report + live SSE stream tied to fit endpoint and reducer. |
| JobHistory.svelte | | "Viewed jobs" list bound to `viewed` my-jobs endpoint. |
| JobMatch.svelte | | Profile-match sidebar bound to `JobMatch`, profile store, match-state logic. |
| JobRelated.svelte | | "More like this" tabbed block (similar/copies) tied to job's related data. |
| JobRow.svelte | | Content-bound job card tied to `Job` shape, enrichment, save/view stores. |
| JobView.svelte | | Job detail page with view/apply/save/report interactions wired to slug. |
| JobsView.svelte | | /jobs list page owning filters, onboarding, infinite scroll, header-scope wiring. |
| LocationPane.svelte | filters | Region→country tree + cities list bound to `COUNTRY_REGION_MAP` and location facets. |
| ModerationView.svelte | | Moderator page with queue/reports tabs bound to submission/report APIs. |
| MySubmissionsView.svelte | | "My submissions" list bound to user's `Submission` records. |
| NoteEditor.svelte | | EasyMDE wrapper mounted only inside JobDrawer for per-application notes. |
| OnboardingAlertBanner.svelte | onboarding | One-off post-onboarding nudge hosting quick SaveSearchAlert. |
| OnboardingBanner.svelte | onboarding | One-off pre-onboarding nudge banner above /jobs feed. |
| OnboardingWizard.svelte | onboarding | /jobs onboarding overlay tied to facet registries and onboarding lifecycle. |
| PipelineFunnel.svelte | | Sankey chart bound to `PIPELINE_BUCKETS` and pipeline stats payload. |
| PipelineView.svelte | | Pipeline page composing RateDonut + PipelineFunnel from `getMyPipeline`. |
| ProfileForm.svelte | | Profile editor form bound to `UserProfile`, specialization caps, CV upload. |
| RealityBadge.svelte | | Badge bound to `Reality` signal and `realityBadge`/`postingContrast` helpers. |
| RecruitersView.svelte | | Marketing page for recruiters with hardcoded benefits/steps copy. |
| ReportDialog.svelte | | Multi-step report-a-job modal wired to report reasons API. |
| ReportQueue.svelte | | Moderator report queue bound to `Report` resource and resolve/dismiss actions. |
| ResumeStructuredView.svelte | | Read-only view bound to `ResumeStructured` parsed-CV shape. |
| SavedJobs.svelte | | "Saved jobs" list bound to `saved` my-jobs endpoint. |
| SavedSearches.svelte | | "My filters" modal tab bound to `savedSearches` store and staged filters. |
| SavedSearchesView.svelte | | /my/searches account page bound to saved searches + Telegram connection. |
| StatusBoard.svelte | | /status page bound to `IngestStatus` rollup and provider-kind taxonomy. |
| SubmitView.svelte | | /submit form wired to job-submission API. |
| SwipeDeck.svelte | | /jobs/swipe page owning drag physics, queue, filters, save/dismiss flow. |
| TopBar.svelte | | Site header: logo + context-swapping search + menu, auth-redirect handling. |
| VerdictView.svelte | | Coverage/verdict body bound to `Verdict` payload and gap/skill tabs. |
| CvEditor.svelte | cv | CV section editor bound to `Document` shape and CV API. |
| CvList.svelte | cv | CV builder landing bound to CV list/create/delete API. |
