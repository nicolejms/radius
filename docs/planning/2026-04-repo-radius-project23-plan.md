# Sprint Plan: Repo Radius & Application Graph — GitHub Project 23

**Date:** 2026-04-17  
**Planning horizon:** April 2026 – July 2026 (5 × 2-week sprints)  
**Team capacity:** 7 developers per sprint  
**Source backlog:** Issues associated with [radius-project/radius Project 23](https://github.com/orgs/radius-project/projects/23)  
**Working repository:** [nicolejms/radius](https://github.com/nicolejms/radius) (fork — planning work isolated here to avoid cluttering upstream PR queue)

---

## Table of Contents

1. [Overview](#1-overview)
2. [Assumptions and Conventions](#2-assumptions-and-conventions)
3. [Epic Map: Grouping Work into Relatable Units](#3-epic-map-grouping-work-into-relatable-units)
4. [Sprint Plan](#4-sprint-plan)
   - [Sprint 1 — Foundations](#sprint-1--foundations-apr-21--may-2-2026)
   - [Sprint 2 — Repo Radius Usable Loop](#sprint-2--repo-radius-usable-loop-may-4--may-15-2026)
   - [Sprint 3 — App Graph Preview Path](#sprint-3--app-graph-preview-path-may-18--may-29-2026)
   - [Sprint 4 — Deployments UX & Deployed Graph](#sprint-4--deployments-ux--deployed-graph-jun-1--jun-12-2026)
   - [Sprint 5 — Hardening, Performance & Compliance](#sprint-5--hardening-performance--compliance-jun-15--jun-26-2026)
5. [Risks and Mitigations](#5-risks-and-mitigations)
6. [Progress Tracking](#6-progress-tracking)
7. [Out of Scope / Deferred](#7-out-of-scope--deferred)

---

## 1. Overview

This document is a **sprint-by-sprint execution plan** for the features and improvements tracked in [radius-project/radius Project 23](https://github.com/orgs/radius-project/projects/23). The two primary themes are:

- **Repo Radius** — running Radius entirely inside a GitHub Actions runner (ephemeral, persistent, lifecycle-aware), powered by a cloud-or-git-backed state store.
- **Application Graph in GitHub** — visualizing the app topology from `app.bicep` and live deployments directly inside the GitHub UI via a browser extension/app.

Supporting work spans security hardening (OIDC, encryption, SBOM), CLI backwards-compatibility (preview flag, deprecation warnings), infrastructure-as-code settings lifecycle, and CI/workflow health.

---

## 2. Assumptions and Conventions

| Assumption | Value |
|---|---|
| Sprint length | 2 weeks |
| Developers per sprint | 7 |
| Target velocity per dev | ~4–5 story points / sprint (stories 1–3 pts) |
| Branch strategy | Feature branches → `main` of `nicolejms/radius`; upstream PRs opened separately when ready |
| Definition of Done | Code merged, CI green, unit tests pass, feature-level docs updated |
| Preview features | Gated behind `--preview` flag until stabilized |
| Canonical issue URLs | `https://github.com/radius-project/radius/issues/<N>` |
| Hard deadline | AutoRest deprecation: **July 1, 2026** (see [#11425](https://github.com/radius-project/radius/issues/11425)) |
| Cross-repo work | Dashboard, `resource-types-contrib`, and `radius-project/extensions` are tracked by owners; integration contracts are defined here |

---

## 3. Epic Map: Grouping Work into Relatable Units

### Epic A — Repo Radius: State, Storage & Lifecycle

**Goal:** Radius can run inside a GitHub Actions runner, persist its state across runs, and support clean lifecycle operations.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| A1 | [#11318](https://github.com/radius-project/radius/issues/11318) | Repo Radius Serialization — define and implement the state schema and read/write path | M |
| A2 | [#11308](https://github.com/radius-project/radius/issues/11308) | Implement Storage Provider abstraction (interface + one concrete impl) | M |
| A3a | [#11604](https://github.com/radius-project/radius/issues/11604) | Repo Radius Git storage backend | M |
| A3b | [#11605](https://github.com/radius-project/radius/issues/11605) | Repo Radius Cloud storage backend | M |
| A4 | [#11606](https://github.com/radius-project/radius/issues/11606) | `rad startup` / `rad shutdown` lifecycle commands | S |
| A5 | [#11367](https://github.com/radius-project/radius/issues/11367) | Background Cleanup Job (GC for stale state/artifacts) | S |
| A6 | [#11493](https://github.com/radius-project/radius/issues/11493) | PR review state serialization | S |

---

### Epic B — Preview Mode & CLI Compatibility

**Goal:** CLI supports opt-in preview semantics and does not break existing flows.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| B1 | [#11637](https://github.com/radius-project/radius/issues/11637) | Implement preview functionality in `rad init` while retaining current behavior | M |
| B2 | [#11490](https://github.com/radius-project/radius/issues/11490) | Support both new and old environments with `rad init --preview` | S |
| B3 | [#11677](https://github.com/radius-project/radius/issues/11677) | Update `rad app *` commands with `--preview` flag | M |
| B4 | [#11489](https://github.com/radius-project/radius/issues/11489) | Update tests to use new types (excluding tf/bicep config tests) | S |
| B5 | [#11488](https://github.com/radius-project/radius/issues/11488) | Update `app.bicep` to point to new namespace | S |
| B6 | [#11486](https://github.com/radius-project/radius/issues/11486) | Add deprecation warning for old Radius types | S |

---

### Epic C — Application Graph

**Goal:** Graph is deterministically produced from app definitions and rendered in the GitHub UI.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| C0 | [#11319](https://github.com/radius-project/radius/issues/11319) | App graph in GitHub: finalize updated spec and plan | S |
| C1 | [#11309](https://github.com/radius-project/radius/issues/11309), [#11574](https://github.com/radius-project/radius/issues/11574) | Graph representation of `app.bicep` + params; deep-dive into design | M |
| C2 | [#11583](https://github.com/radius-project/radius/issues/11583) | Preview/dry-run version of app graph visualization | M |
| C3 | [#11584](https://github.com/radius-project/radius/issues/11584) | Deployed version of app graph visualization | M |
| C4 | [#11581](https://github.com/radius-project/radius/issues/11581) | GUI Extension for Graph — shell (nav, auth, loading states) | M |
| C5 | [#11582](https://github.com/radius-project/radius/issues/11582) | Visualizing `app.bicep` in GUI extension (render pipeline) | M |
| C6 | [#208](https://github.com/radius-project/dashboard/issues/208) | Bug: property references should create links in app graph | S |
| C7 | [#11656](https://github.com/radius-project/radius/issues/11656) | Graph interaction outside of control plane | L |

---

### Epic D — GitHub Deployments Integration & Repo UI

**Goal:** Developers can discover applications, trigger deployments, and view deployment status without leaving GitHub.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| D1 | [#11653](https://github.com/radius-project/radius/issues/11653) | Backend APIs / interfaces to enable UI elements for Repo Radius | L |
| D2 | [#11652](https://github.com/radius-project/radius/issues/11652) | GUI extension for applications section on repository page | M |
| D3 | [#11577](https://github.com/radius-project/radius/issues/11577) | GUI extension for "Define an application" button | S |
| D4 | [#11580](https://github.com/radius-project/radius/issues/11580) | GUI extension for dispatching a deployment workflow | M |
| D5 | [#11660](https://github.com/radius-project/radius/issues/11660) | Radius integration with GitHub Deployments API | L |
| D6 | [#11657](https://github.com/radius-project/radius/issues/11657) | Radius deploys to a Kubernetes cluster defined on the environment | M |

---

### Epic E — Performance & Persistence in GitHub Runners

**Goal:** Radius is fast and reliable in ephemeral CI environments.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| E1 | [#11655](https://github.com/radius-project/radius/issues/11655) | Persist Radius data store between executions | M |
| E2 | [#11659](https://github.com/radius-project/radius/issues/11659) | Optimize Radius performance when running in GH action runner | M |

---

### Epic F — Security, Auth & Compliance

**Goal:** Workflows authenticate securely, secrets are encrypted, and the project meets baseline security standards.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| F1 | [#11492](https://github.com/radius-project/radius/issues/11492) | Add OIDC configuration to GitHub workflow | S |
| F2 | [#11654](https://github.com/radius-project/radius/issues/11654) | Radius leverages OIDC to authenticate to AWS (Core) | M |
| F3 | [#11446](https://github.com/radius-project/radius/issues/11446) | Add an SBOM (Software Bill of Materials) | S |
| F4 | [#11273](https://github.com/radius-project/radius/issues/11273) | Security: Data store and queue hardening for UCP | M |
| F5 | [#11271](https://github.com/radius-project/radius/issues/11271) | Security: Add Kubernetes NetworkPolicies for UCP and Radius components | M |
| F6 | [#11258](https://github.com/radius-project/radius/issues/11258) | Add a security insights file for OpenSSF tooling | S |
| F7 | [#11095](https://github.com/radius-project/radius/issues/11095) | Backend encryption updates | M |
| F8 | [#11307](https://github.com/radius-project/radius/issues/11307) | Implement secrets changes for `bicepSettings` and `terraformSettings` | M |

---

### Epic G — Configuration & Settings Lifecycle

**Goal:** Declarative configuration of Radius is supported end-to-end, including Bicep/Terraform settings.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| G1 | [#11658](https://github.com/radius-project/radius/issues/11658) | Fully declarative Radius configuration | L |
| G2 | [#11306](https://github.com/radius-project/radius/issues/11306) | Design doc: using `Radius.Core/Security` resource for `bicepSettings`/`terraformSettings` | S |
| G3 | [#11354](https://github.com/radius-project/radius/issues/11354) | Terraform/Bicep Settings Lifecycle — Handover | M |
| G4 | [#10615](https://github.com/radius-project/radius/issues/10615) | Terraform/Bicep Settings Lifecycle (broader tracking) | L |

---

### Epic H — CI/Workflow Health & Developer Experience

**Goal:** CI is reliable, developer workflows are smooth, and automation is in place for recurring tasks.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| H1 | [#11482](https://github.com/radius-project/radius/issues/11482) | Fix Close Stale Pull Requests workflow (broken) | XS |
| H2 | [#11425](https://github.com/radius-project/radius/issues/11425) | AutoRest deprecation migration (deadline: July 1, 2026) | L |
| H3 | [#11487](https://github.com/radius-project/radius/issues/11487) | Automate publishing extensions | M |
| H4 | [#11424](https://github.com/radius-project/radius/issues/11424) | Add agentic workflow to keep architecture docs up to date | M |
| H5 | [#11320](https://github.com/radius-project/radius/issues/11320) | Local functional tests (ability to run functional tests locally) | M |
| H6 | [#11321](https://github.com/radius-project/radius/issues/11321) | E2E developer flow refinements | M |
| H7 | [#11317](https://github.com/radius-project/radius/issues/11317) | Prototype for `resource-types-contrib` | M |

---

### Epic I — Demo & Documentation

**Goal:** The work done is demonstrable, discoverable, and documented for contributors and users.

| Slice | Issue(s) | Description | Size |
|---|---|---|---|
| I1 | [#11331](https://github.com/radius-project/radius/issues/11331) | Create demo | M |
| I2 | [#11310](https://github.com/radius-project/radius/issues/11310) | Demo (end-to-end walkthrough) | M |

---

## 4. Sprint Plan

### Capacity model

- 7 developers × 2 weeks = ~14 dev-weeks capacity
- Estimated usable capacity after meetings/review: ~11–12 dev-weeks
- Slices are sized: XS=0.25dw, S=0.5dw, M=1dw, L=2dw

---

### Sprint 1 — Foundations (Apr 21 – May 2, 2026)

**Theme:** Create the "platform rails" so all subsequent work is unblocked. Establish state persistence, OIDC auth, and the preview flag scaffold.

**Sprint goal:** By end of Sprint 1, a minimal state round-trip (write → read → survive restart) is implemented, the `--preview` flag is wired end-to-end in the CLI, and CI workflows can authenticate via OIDC.

#### Committed work (≈11.75 dev-weeks)

| # | Slice | Description | Owner notes | Size |
|---|---|---|---|---|
| 1 | A1 | Repo Radius Serialization — state schema, read/write path, migration stubs | Blocks A3a, A3b, A4, E1 | M |
| 2 | A2 | Implement Storage Provider abstraction (interface + in-memory impl for testing) | Blocks A3a, A3b | M |
| 3 | F1 | Add OIDC configuration to GitHub workflow ([#11492](https://github.com/radius-project/radius/issues/11492)) | Unblocks F2 | S |
| 4 | F6 | Add security insights file for OpenSSF ([#11258](https://github.com/radius-project/radius/issues/11258)) | Small, high-visibility | S |
| 5 | B1 | Implement preview functionality scaffold; `--preview` flag accepted, wired into pipeline | Unblocks B2, B3 | M |
| 6 | B6 | Add deprecation warnings for old `Applications.*` types ([#11486](https://github.com/radius-project/radius/issues/11486)) | Non-breaking; consistent UX | S |
| 7 | H1 | Fix Close Stale PRs workflow ([#11482](https://github.com/radius-project/radius/issues/11482)) | CI health; blocks contributor velocity | XS |
| 8 | C0 | Finalize app graph updated spec and plan ([#11319](https://github.com/radius-project/radius/issues/11319)) | Design; produces contracts for C1–C5 | S |
| 9 | G2 | Design doc: `Radius.Core/Security` resource for settings ([#11306](https://github.com/radius-project/radius/issues/11306)) | Design; blocks G3, F8 | S |
| 10 | B5 | Update `app.bicep` to new namespace ([#11488](https://github.com/radius-project/radius/issues/11488)) | Small; needed to unblock tests | S |
| 11 | H2 (start) | **AutoRest migration spike**: inventory all AutoRest usages and produce a migration plan before the July 1 deadline | Risk mitigation | M |

#### Exit criteria — Sprint 1

- [ ] **State persistence MVP:** A unit test proves state can be written, read back, and survives a simulated "restart" using the in-memory storage provider
- [ ] **Storage interface:** `StorageProvider` interface is code-reviewed and merged; in-memory implementation passes its test suite
- [ ] **OIDC auth:** At least one GitHub Actions workflow in the fork successfully authenticates to AWS (or Azure) using OIDC with least-privilege IAM role; documented in PR
- [ ] **Preview flag:** `rad init --preview` is accepted by the CLI without errors; produces a clearly-labelled plan-like output even if full functionality is not yet implemented
- [ ] **Deprecation warning:** `rad deploy` or `rad run` emits a visible, actionable deprecation warning when old `Applications.*` types are detected
- [ ] **CI green** in `nicolejms/radius` including updated stale-PR workflow
- [ ] **AutoRest spike output:** A markdown doc listing all files/generators that depend on AutoRest, proposed migration path (e.g., to TypeSpec-based generation), and owners. Shared with team before Sprint 2 planning
- [ ] App graph schema draft circulated for team review

#### Sprint 1 dependencies

- OIDC trust policy must be configured in the cloud account before F1 can be tested
- App graph spec (C0) must be completed before Sprint 3 C1 implementation begins

---

### Sprint 2 — Repo Radius Usable Loop (May 4 – May 15, 2026)

**Theme:** Make Repo Radius "repeatable": start, stop, persist state, and store artifacts in Git and cloud backends.

**Sprint goal:** By end of Sprint 2, `rad startup` → work → `rad shutdown` → restart is functional with both Git and Cloud storage backends. Test suite is updated for new types.

#### Committed work (≈11.5 dev-weeks)

| # | Slice | Description | Owner notes | Size |
|---|---|---|---|---|
| 1 | A3a | Repo Radius Git storage backend ([#11604](https://github.com/radius-project/radius/issues/11604)) | Depends on A1, A2 | M |
| 2 | A3b | Repo Radius Cloud storage backend ([#11605](https://github.com/radius-project/radius/issues/11605)) | Depends on A1, A2 | M |
| 3 | A4 | `rad startup` / `rad shutdown` commands ([#11606](https://github.com/radius-project/radius/issues/11606)) | Depends on A1 | S |
| 4 | B2 | `rad init --preview` supports both old and new env types ([#11490](https://github.com/radius-project/radius/issues/11490)) | Depends on B1 | S |
| 5 | B4 | Update tests to use new types (excl. tf/bicep config) ([#11489](https://github.com/radius-project/radius/issues/11489)) | Test hygiene | S |
| 6 | F2 | Radius leverages OIDC to authenticate to AWS Core ([#11654](https://github.com/radius-project/radius/issues/11654)) | Depends on F1 | M |
| 7 | G2→G3 | Terraform/Bicep Settings Lifecycle handover ([#11354](https://github.com/radius-project/radius/issues/11354)) | Depends on G2 design | M |
| 8 | H5 | Local functional tests infrastructure ([#11320](https://github.com/radius-project/radius/issues/11320)) | Developer experience | M |
| 9 | H6 | E2E developer flow refinements ([#11321](https://github.com/radius-project/radius/issues/11321)) | Paper-cuts from Sprint 1 work | M |
| 10 | H2 (cont.) | AutoRest migration — begin file-by-file migration for highest-risk generators | Critical path to July 1 deadline | M |

#### Exit criteria — Sprint 2

- [ ] **Lifecycle commands:** `rad startup` and `rad shutdown` execute without errors on a sample app and are documented in the CLI reference
- [ ] **Git storage:** Can write and read Repo Radius state to/from a Git commit; integration test passes in CI
- [ ] **Cloud storage:** Can write and read Repo Radius state to/from the configured cloud storage backend; integration test passes in CI
- [ ] **Persistence round-trip:** Running `rad startup` after shutdown restores prior state from both storage backends (demonstrated in a test or demo script)
- [ ] **OIDC for AWS:** An end-to-end functional test deploys a Radius resource using OIDC-provisioned AWS credentials; no static credentials in code or CI secrets
- [ ] **Test suite:** Test files updated for new types; no test regressions
- [ ] **Local functional tests:** Developers can run the functional test suite locally without a full cloud environment (at minimum a documented subset)
- [ ] **AutoRest migration:** ≥50% of identified AutoRest-dependent generators migrated or migration PRs open for review

---

### Sprint 3 — App Graph Preview Path (May 18 – May 29, 2026)

**Theme:** Produce and render the preview (pre-deploy) application graph inside the GitHub UI.

**Sprint goal:** By end of Sprint 3, running `rad deploy --preview` produces a schema-valid graph artifact for a sample app, and the GUI extension can load and display it.

#### Committed work (≈11.5 dev-weeks)

| # | Slice | Description | Owner notes | Size |
|---|---|---|---|---|
| 1 | C1 | Graph representation of `app.bicep` + params ([#11309](https://github.com/radius-project/radius/issues/11309), [#11574](https://github.com/radius-project/radius/issues/11574)) | Core data model; blocks C2–C5 | M |
| 2 | C2 | Preview/dry-run app graph generation ([#11583](https://github.com/radius-project/radius/issues/11583)) | Depends on C1 | M |
| 3 | C4 | GUI Extension for Graph shell — nav, auth, loading states ([#11581](https://github.com/radius-project/radius/issues/11581)) | Parallelizable with C2 | M |
| 4 | C5 | Visualizing `app.bicep` in GUI extension ([#11582](https://github.com/radius-project/radius/issues/11582)) | Depends on C2, C4 | M |
| 5 | B3 | Update `rad app *` commands with `--preview` flag ([#11677](https://github.com/radius-project/radius/issues/11677)) | Depends on B1 | M |
| 6 | C6 | Bug: property references should create links in app graph ([#208](https://github.com/radius-project/dashboard/issues/208)) | Bug; fix with regression test | S |
| 7 | A6 | PR review state serialization ([#11493](https://github.com/radius-project/radius/issues/11493)) | Needed for Repo Radius in PRs | S |
| 8 | F8 | Implement secrets changes for `bicepSettings` / `terraformSettings` ([#11307](https://github.com/radius-project/radius/issues/11307)) | Depends on G2 | M |
| 9 | H2 (cont.) | Complete AutoRest migration for all identified generators | Hard deadline: July 1, 2026 | M |

#### Exit criteria — Sprint 3

- [ ] **Graph data model:** The graph schema is finalized, versioned, and documented; a sample `app.bicep` can be parsed into a valid graph object and serialized to JSON
- [ ] **Preview graph:** `rad deploy --preview` produces a graph artifact that passes schema validation for a reference app; artifact is accessible to the GUI extension
- [ ] **GUI extension shell:** Extension loads in a browser, authenticates, and displays a loading/error state gracefully for missing data
- [ ] **Graph rendering:** Extension renders the preview graph for the reference app; nodes and edges are clickable
- [ ] **Property-reference bug fixed:** A regression test confirms that resources referenced as properties (not explicit connections) appear as linked nodes in the graph
- [ ] **`rad app *` preview:** `rad app graph --preview` and `rad app list --preview` are consistent with the preview contract established in Sprint 1
- [ ] **AutoRest migration:** All identified generators migrated or replacement PRs merged; CI no longer depends on deprecated AutoRest tooling

---

### Sprint 4 — Deployments UX & Deployed Graph (Jun 1 – Jun 12, 2026)

**Theme:** End-to-end deployment dispatch from the GitHub UI; display the post-deploy application graph.

**Sprint goal:** By end of Sprint 4, a developer can click "Deploy" in the GitHub UI, the deployment is dispatched via the GitHub Deployments API, and the resulting deployed graph is visible in the extension.

#### Committed work (≈11.5 dev-weeks)

| # | Slice | Description | Owner notes | Size |
|---|---|---|---|---|
| 1 | C3 | Deployed version of app graph visualization ([#11584](https://github.com/radius-project/radius/issues/11584)) | Depends on C1 | M |
| 2 | D1 | Backend APIs / interfaces for Repo Radius UI elements ([#11653](https://github.com/radius-project/radius/issues/11653)) | Blocks D2–D5 | L |
| 3 | D2 | GUI extension: applications section on repository page ([#11652](https://github.com/radius-project/radius/issues/11652)) | Depends on D1 | M |
| 4 | D3 | GUI extension: "Define an application" button ([#11577](https://github.com/radius-project/radius/issues/11577)) | Depends on D1 | S |
| 5 | D4 | GUI extension: dispatch deployment workflow ([#11580](https://github.com/radius-project/radius/issues/11580)) | Depends on D1 | M |
| 6 | D5 | Radius integration with GitHub Deployments API ([#11660](https://github.com/radius-project/radius/issues/11660)) | Depends on D1 | L |
| 7 | D6 | Radius deploys to K8s cluster defined on the environment ([#11657](https://github.com/radius-project/radius/issues/11657)) | Depends on D5 | M |
| 8 | H3 | Automate publishing extensions ([#11487](https://github.com/radius-project/radius/issues/11487)) | CI; needed before demo | M |

#### Exit criteria — Sprint 4

- [ ] **Deployed graph:** After a successful `rad deploy`, the extension displays the live deployed application graph; nodes reflect actual deployed resource IDs/status
- [ ] **Backend APIs:** The Repo Radius UI APIs are versioned, have OpenAPI specs, and pass contract tests
- [ ] **Repo page section:** The GitHub UI extension shows an "Applications" section on the repository page listing apps from the current repo
- [ ] **Deployment dispatch:** A user can click "Deploy" in the UI, choose an environment, and a GitHub Actions workflow run is triggered; status is reflected back in the UI
- [ ] **GitHub Deployments integration:** Radius creates and updates a GitHub Deployment record on workflow dispatch; deployment status (in-progress / success / failure) is visible in the GitHub repo's Deployments section
- [ ] **K8s from environment:** `rad deploy` correctly targets the Kubernetes cluster referenced in the Radius environment definition (not hardcoded)
- [ ] **Extension publishing:** Extension build is automated in CI; a tagged release produces a distributable artifact without manual steps
- [ ] **Happy-path demo script** documented (steps + expected screenshots) and successfully executed against the fork environment

---

### Sprint 5 — Hardening, Performance & Compliance (Jun 15 – Jun 26, 2026)

**Theme:** Make everything reliable, fast, and secure in CI environments; complete compliance items.

**Sprint goal:** By end of Sprint 5, performance benchmarks in GitHub runners are published, the encryption stack is merged, the SBOM is generated in CI, and backlog deferred items are closed out.

#### Committed work (≈11.5 dev-weeks)

| # | Slice | Description | Owner notes | Size |
|---|---|---|---|---|
| 1 | E1 | Persist Radius data store between executions ([#11655](https://github.com/radius-project/radius/issues/11655)) | Depends on A1, A2 | M |
| 2 | E2 | Optimize Radius performance in GH action runner ([#11659](https://github.com/radius-project/radius/issues/11659)) | Depends on E1 | M |
| 3 | A5 | Background Cleanup Job ([#11367](https://github.com/radius-project/radius/issues/11367)) | Depends on A1 | S |
| 4 | F3 | Add SBOM ([#11446](https://github.com/radius-project/radius/issues/11446)) | Compliance | S |
| 5 | F4 | Security: data store and queue hardening for UCP ([#11273](https://github.com/radius-project/radius/issues/11273)) | Security | M |
| 6 | F5 | Security: Add Kubernetes NetworkPolicies ([#11271](https://github.com/radius-project/radius/issues/11271)) | Security | M |
| 7 | F7 | Backend encryption updates ([#11095](https://github.com/radius-project/radius/issues/11095)) | Security | M |
| 8 | H4 | Agentic workflow to keep architecture docs up to date ([#11424](https://github.com/radius-project/radius/issues/11424)) | DevEx | M |
| 9 | I1/I2 | Create and record end-to-end demo ([#11331](https://github.com/radius-project/radius/issues/11331), [#11310](https://github.com/radius-project/radius/issues/11310)) | Stakeholder comms | M |

#### Exit criteria — Sprint 5

- [ ] **Persistence:** A GitHub Actions workflow that runs twice in sequence proves state is restored from the previous run; documented cache invalidation strategy for secrets boundaries
- [ ] **Performance benchmarks:** Wall-clock time and peak memory for a reference app deploy in a GitHub-hosted runner are measured before and after optimizations; results published in the repo (e.g., `docs/planning/benchmarks-sprint5.md`)
- [ ] **Cleanup job:** Background cleanup runs on a schedule, respects configurable age threshold, and has a `--dry-run` mode; no production data is deleted without explicit configuration
- [ ] **SBOM:** SBOM artifact is generated on every release CI run and attached to the GitHub Release; format is CycloneDX or SPDX
- [ ] **Encryption stack:** `bicepSettings` and `terraformSettings` secrets are encrypted at rest using the agreed-upon encryption provider; no plaintext secrets in state store
- [ ] **Network policies:** Kubernetes NetworkPolicies for UCP and Radius components are defined in Helm chart and tested in a functional environment
- [ ] **Architecture docs:** Agentic workflow PR is open and demonstrated with at least one auto-updated architecture document
- [ ] **Demo:** Recorded demo (video or screenshot walkthrough) covers the full Repo Radius + App Graph user journey and is linked from the repository README

---

## 5. Risks and Mitigations

### Risk 1 — AutoRest Deprecation (Critical, Jul 1, 2026 hard deadline)

| Field | Detail |
|---|---|
| **Issue** | [#11425](https://github.com/radius-project/radius/issues/11425) |
| **Deadline** | July 1, 2026 (AutoRest tool deprecated by Microsoft) |
| **Impact** | Any generator still using AutoRest will break CI and block releases |
| **Probability** | High if unmitigated |
| **Mitigation** | Start spike in Sprint 1 to inventory usage; target full migration by end of Sprint 3 (≥3 weeks before deadline). Assign a dedicated migration lead. Track progress in the spike doc. |
| **Contingency** | If migration is not complete by Sprint 3 end, freeze feature work for 1–2 devs in Sprint 4 to close remaining gaps before the deadline. |

---

### Risk 2 — Cross-Repo Dependencies (Dashboard, resource-types-contrib, extensions)

| Field | Detail |
|---|---|
| **Issue** | Multiple slices (C4, C5, D2–D5) require changes in `radius-project/dashboard` and `resource-types-contrib` |
| **Impact** | Integration delays; extension features can't ship if the dashboard repo is blocked |
| **Probability** | Medium |
| **Mitigation** | Define and freeze integration contracts (graph JSON schema, API specs) at end of Sprint 1. Use API mocks in the fork to allow parallel development. Pin explicit interface versions. |
| **Contingency** | If upstream repos are blocked, ship preview-gated "stub" UI that surfaces the contract, allowing backend work to continue. |

---

### Risk 3 — Preview Contract Instability

| Field | Detail |
|---|---|
| **Issue** | If the `--preview` output contract changes between sprints, UI code and tests built on it churn |
| **Impact** | High rework cost in Sprints 3–5 |
| **Probability** | Medium |
| **Mitigation** | Explicitly version the preview output schema (v1alpha1 etc.) in Sprint 1. Gate breaking changes behind a new version. Require schema review for any API changes. |
| **Contingency** | Maintain a compatibility shim for at least one version behind the current schema. |

---

### Risk 4 — OIDC & Cloud Auth Complexity

| Field | Detail |
|---|---|
| **Issue** | [#11492](https://github.com/radius-project/radius/issues/11492), [#11654](https://github.com/radius-project/radius/issues/11654) |
| **Impact** | Broken auth blocks all cloud-dependent tests and demo scenarios |
| **Probability** | Medium (OIDC trust policies are easy to misconfigure) |
| **Mitigation** | Document the required IAM role trust policy and minimum permissions. Test in a sandbox environment in Sprint 1 before gating other work. Add a dedicated "auth health check" workflow that runs nightly. |
| **Contingency** | Fall back to temporary service-principal credentials (scoped, short-lived) while OIDC is investigated; replace before release. |

---

### Risk 5 — CI Flakiness Blocking Sprint Velocity

| Field | Detail |
|---|---|
| **Issue** | [#11275](https://github.com/radius-project/radius/issues/11275) (flaky test), [#11482](https://github.com/radius-project/radius/issues/11482) (broken workflow) and similar |
| **Impact** | Flaky CI forces re-runs, slows reviews, and obscures real failures |
| **Probability** | High (observable today) |
| **Mitigation** | Fix H1 in Sprint 1. Assign a rotating "CI health" dev each sprint. Quarantine (skip + file issue) any test that flakes 2+ times in a sprint. |
| **Contingency** | Introduce a flakiness dashboard; block merges on a flaky test appearing more than N times per week. |

---

### Risk 6 — Ephemeral Runner Security (Secrets in Cache)

| Field | Detail |
|---|---|
| **Issue** | [#11655](https://github.com/radius-project/radius/issues/11655) |
| **Impact** | Caching Radius state between runs may inadvertently persist credentials or sensitive data |
| **Probability** | Medium-high if not designed carefully |
| **Mitigation** | Design the persistence layer to explicitly separate mutable state from credentials. Credentials must never enter the cache. Security review the cache write path before E1 merges. |
| **Contingency** | Disable cross-run caching for secrets-bearing state; require re-authentication on each run even if state is cached. |

---

### Risk 7 — Scope Creep from 189 Issues

| Field | Detail |
|---|---|
| **Issue** | Project 23 has 189 issues; this plan covers the highest-priority 60–70 |
| **Impact** | Unplanned work displaces committed sprint items |
| **Probability** | High |
| **Mitigation** | Gate sprint entry: any issue not in the plan must be approved by the team lead before entering the sprint. Reserve ~1 dev-week per sprint as a buffer. |
| **Contingency** | Defer low-priority items to Sprint 6 (post-July) or mark as backlog. |

---

## 6. Progress Tracking

### Sprint board

Each sprint's committed work items map to issues in `radius-project/radius`. Track progress by:

1. **Labels:** Apply `sprint-1`, `sprint-2`, … `sprint-5` labels to issues in the upstream repo (or in the fork's issue tracker) as issues are picked up.
2. **Milestones:** Create a GitHub Milestone per sprint in `nicolejms/radius` and attach fork-side issues/PRs.
3. **Weekly sync:** 15-minute team sync every Monday to review completed items, surface blockers, and update the risk register.

### Sprint-to-issue mapping (quick reference)

| Sprint | Slices | Key issue numbers |
|---|---|---|
| Sprint 1 | A1, A2, F1, F6, B1, B6, H1, C0, G2, B5, H2(start) | #11318, #11308, #11492, #11258, #11637, #11486, #11482, #11319, #11306, #11488, #11425 |
| Sprint 2 | A3a, A3b, A4, B2, B4, F2, G3, H5, H6, H2(cont.) | #11604, #11605, #11606, #11490, #11489, #11654, #11354, #11320, #11321, #11425 |
| Sprint 3 | C1, C2, C4, C5, B3, C6, A6, F8, H2(done) | #11309, #11583, #11581, #11582, #11677, #208, #11493, #11307, #11574, #11425 |
| Sprint 4 | C3, D1, D2, D3, D4, D5, D6, H3 | #11584, #11653, #11652, #11577, #11580, #11660, #11657, #11487 |
| Sprint 5 | E1, E2, A5, F3, F4, F5, F7, H4, I1/I2 | #11655, #11659, #11367, #11446, #11273, #11271, #11095, #11424, #11331, #11310 |

### Definition of a "closed" issue (for this plan)

An issue is considered **done** when:

- Code is merged to `main` of `nicolejms/radius` (or upstream, if appropriate)
- Unit tests added or updated
- Any new CLI commands have at least a one-sentence entry in the relevant help text
- Any new APIs have an OpenAPI spec entry
- CI is green

---

## 7. Out of Scope / Deferred

The following open issues from Project 23 are acknowledged but not scheduled in this plan. They should be revisited in the Sprint 6+ backlog grooming:

| Issue | Title | Reason for deferral |
|---|---|---|
| [#11658](https://github.com/radius-project/radius/issues/11658) | Fully declarative Radius configuration | Large, architectural; needs design before implementation |
| [#11656](https://github.com/radius-project/radius/issues/11656) | Graph interaction outside of control plane | Scope not yet defined; needs spike |
| [#11590](https://github.com/radius-project/radius/issues/11590) | Agent Skill for generating app.bicep | Nice-to-have; depends on graph work completing |
| [#11441](https://github.com/radius-project/radius/issues/11441) | Guided Radius install mission in KubeStellar Console | External dependency; tracked separately |
| [#11217](https://github.com/radius-project/radius/issues/11217) | Cosmetic bug: confusing error in `rad app list` | Low priority; can be batched |
| [#11218](https://github.com/radius-project/radius/issues/11218) | `rad upgrade kubernetes` resets to chart defaults | Medium priority; no Sprint 1–5 dependency |
| [#11147](https://github.com/radius-project/radius/issues/11147) | Add `--force` to `rad resource delete` | Good DX improvement; low urgency |
| [#11108](https://github.com/radius-project/radius/issues/11108) | Automatically sync core resource manifests | Automation; deferred until after configuration lifecycle is stable |
| [#9830](https://github.com/radius-project/radius/issues/9830) | Integrate new core RRTs into Radius | Depends on #11317 prototype | 
| [#11253](https://github.com/radius-project/radius/issues/11253) | Shared PVC for Terraform storage across RPs | Infrastructure; needs capacity after Sprint 5 |
| [#11252](https://github.com/radius-project/radius/issues/11252) | Fix TerraformResource deployment + private module deployment | Bug; triage for Sprint 6 |
| [#11292](https://github.com/radius-project/radius/issues/11292) | Insufficient permissions on dynamic-rp SA for Dapr Terraform Recipe | Bug; triage for Sprint 6 |
| [#11291](https://github.com/radius-project/radius/issues/11291) | Dashboard image missing linux/arm64 manifest | Cross-repo bug; track in dashboard repo |

---

*Last updated: 2026-04-17*  
*Maintainer: planning tracked in [nicolejms/radius](https://github.com/nicolejms/radius) — please open issues or PRs there for corrections.*
