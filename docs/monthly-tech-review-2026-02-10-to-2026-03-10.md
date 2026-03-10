# Radius Monthly Technical Review

**Review Period:** February 10, 2026 – March 10, 2026

## Executive Summary

This review period was highly productive for the Radius project, culminating in the **v0.55.0 release** on March 6, 2026. The team delivered significant feature work around sensitive data management, major build infrastructure modernization, comprehensive architecture documentation, and numerous CI/CD improvements. A total of **59 commits** landed across **123 unique files**, representing **~17,966 insertions** and **~12,540 deletions** from **11+ contributors**.

---

## 1. Major Release

### Radius v0.55.0 (Released March 6, 2026)

The v0.55.0 release was the primary milestone of this period, progressing through **8 release candidates** (rc1–rc8) before final release. Key release highlights include:

- **Sensitive data management** via the `x-radius-sensitive` annotation
- **Fixed environment resource type deployment** for `rad deploy` with Bicep templates
- **CLI fixes** for `rad app graph --output json`, `rad version --output json`, and `rad rollback kubernetes --list-revisions`
- **Dashboard improvements** for custom resource types with non-default API versions

Release PRs: #11335 (rc5), #11352 (rc6), #11372 (rc7), #11375 (rc8), #11392 (final)

---

## 2. Feature Updates

### 2.1 Sensitive Data Management (Secret Redaction)

**PRs:** #11195, #11233, #11348

The largest feature work this period was the implementation of secret redaction for user-defined resource types. This multi-PR effort introduced:

- **Backend, GET, and LIST controllers** for secret decryption and redaction in the dynamic resource provider (PR #11195 — 22 files, ~2,952 insertions). Fields marked with `x-radius-sensitive` are encrypted at rest and redacted (set to `null`) on GET/LIST operations.
- **Sensitive string schema validation** updates to support both string and encrypted object representations (PR #11233 — 782 insertions).
- **End-to-end functional tests** validating sensitive field encryption on creation/update and redaction on retrieval (PR #11348 — 400 insertions with new Bicep templates and resource type schemas).

Key files changed: `pkg/dynamicrp/frontend/getresource.go`, `pkg/dynamicrp/frontend/listresources.go`, `pkg/portableresources/backend/controller/createorupdateresource.go`, `pkg/schema/annotations.go`, `pkg/crypto/encryption/sensitive.go`

### 2.2 Deployment Engine Image Publishing

**PR:** #11365

Refactored the Deployment Engine image publishing workflow to use a remote repository dispatch approach, improving the release pipeline for the `bicep-de` container image. Updated both the `publish-de-image.yaml` and `release.yaml` workflows (118 insertions).

### 2.3 RBAC Permissions for Dynamic-RP

**PR:** #11350

Added missing RBAC permissions to the `dynamic-rp` service account ClusterRole to support Dapr Recipe deployments:

- `apiextensions.k8s.io/customresourcedefinitions` — `get`, `list`, `watch` (for CRD GVK resolution)
- `dapr.io` resources — full CRUD for `components`, `subscriptions`, `configurations`, `resiliencies`

### 2.4 Gopls MCP Server Integration

**PR:** #11270

Added a `gopls` MCP (Model Context Protocol) server configuration to the development environment, improving Go language server support for contributors using AI-assisted development tools.

### 2.5 CODEOWNERS Update for On-Call

**PR:** #11342

Added on-call team management for GitHub Actions updates in the CODEOWNERS file, improving the review workflow for CI/CD changes.

---

## 3. Defect Fixes

### 3.1 CLI Output Format Error Message

**PR:** #11220 (10 files, 304 insertions)

Fixed a confusing error message when passing unsupported output formats to CLI commands. The fix:

- Added `NormalizeFormat()` to map `"plain-text"` → `"table"` for backwards compatibility
- Rewrote `RequireOutput()` with unified validation and clear error messages listing supported formats
- Removed `plain-text` as an advertised format while maintaining backwards compatibility

### 3.2 JSON Output for `rad app graph`

**PR:** #11207 (5 files, 121 insertions)

Fixed a bug where `rad app graph -o json` was not outputting valid JSON. Added proper JSON formatter support for the application graph command.

### 3.3 Helm Upgrade Value Preservation

**PR:** #6 (3 files, 139 insertions)

Fixed `rad upgrade kubernetes` to preserve existing Helm release values during upgrades by setting `ReuseValues=true` in the Helm client, preventing loss of user-customized configuration.

### 3.4 LRT Database Configuration on Upgrade

**PR:** #11371

Fixed the long-running test workflow to set `database.enabled=false` during LRT upgrade operations, preventing test failures related to database configuration conflicts.

### 3.5 RC-to-RC Upgrade Support

**PR:** #11351 (2 files, 53 insertions)

Fixed upgrade logic to allow successful upgrades between release candidate versions, which was previously failing during the rc5–rc8 release candidate cycle.

### 3.6 Namespace Deletion Timeout

**PR:** #11374

Increased timeout for namespace deletion operations in tests, resolving intermittent failures caused by slow Kubernetes namespace cleanup.

### 3.7 Docker Command Ordering

**PR:** #11210

Corrected the ordering of Docker commands in documentation/scripts for accuracy.

### 3.8 Gosec SARIF Output

**PR:** #11386

Rolled back the Gosec security scanner version to fix invalid SARIF output that was breaking the security scanning pipeline.

---

## 4. Test Reliability Improvements

### 4.1 Flaky Test Fixes for Race Conditions

**PR:** #11339 (11 files, 60 insertions)

Comprehensive fix for multiple flaky test issues:

- Replaced `require.NoError` with `panic` in non-test goroutines to avoid data races on `testing.T`
- Reduced poll intervals from 1s to 200ms across controller reconciler and flux controller tests
- Fixed proxy async test to accept both `Accepted` and `Updating` as valid intermediate provisioning states
- Replaced in-memory filesystem with OS filesystem in chart validation tests
- Removed `-v` from default unit test invocation to reduce CI noise

### 4.2 Recipe Reconciler Secret Deletion Fix

**PR:** #11276

Fixed flaky test assertions for recipe reconciler secret deletion by improving timing and assertion logic.

### 4.3 Consolidated Preflight Upgrade Tests

**PR:** #11202 (211 insertions, 274 deletions)

Consolidated and simplified preflight upgrade tests with tighter polling intervals, reducing test execution time and improving reliability.

---

## 5. Build System & CI/CD Improvements

### 5.1 Bicep-Types Submodule Removal and pnpm Migration

**PR:** #11139 (36 files, 7,853 insertions, 9,897 deletions)

**The largest change this period.** Removed the `bicep-types` git submodule and migrated all JavaScript/TypeScript package management from npm to pnpm:

- Removed `.gitmodules` and the `bicep-types` submodule reference
- Migrated `typespec/`, `hack/bicep-types-radius/src/generator/`, and `hack/bicep-types-radius/src/autorest.bicep/` from npm to pnpm
- Updated all GitHub workflows, Makefiles, and devcontainer configuration
- Added a developer migration guide at `docs/contributing/bicep-types-migration-guide.md`
- Added a root `package.json` for consistent dev dependency management

### 5.2 Bicep Types Publishing Workflow Refactor

**PRs:** #11257, #11263, #11264

Complete refactoring of the bicep types publishing workflow across three PRs:

- Restructured the workflow for better maintainability (PR #11257 — 191 insertions)
- Updated configurations for the new pnpm-based build (PR #11263)
- Added dispatch-based publishing in functional tests (PR #11264)

### 5.3 Go 1.26.0 Upgrade

**PR:** #11227 (6 files)

Bumped the Go version to 1.26.0 across all Dockerfiles, CI workflows, and development configuration.

### 5.4 Unit Test Workflow Improvements

**PR:** #11230 (2 files, 18 insertions, 67 deletions)

Improved the unit test workflow for fork compatibility and direct triggering, making it easier for external contributors to validate their changes.

### 5.5 Doc-Only Change Detection

**PR:** #11283

Added change detection to `unit-tests.yaml` to skip test execution on documentation-only changes, saving CI compute time.

### 5.6 Workflow Summary Job for PR Status

**PR:** #11236 (68 insertions)

Added a summary job to the build workflow that aggregates status across all PR checks, providing a single status check for branch protection rules.

### 5.7 Functional Test Cloud Artifact Handling

**PR:** #11232

Improved artifact handling in the `functional-test-cloud` workflow for more reliable test artifact management.

### 5.8 pnpm Version Correction

**PR:** #11228

Fixed the pnpm version setup in the `generate-pnpm-installed` Make target after the npm-to-pnpm migration.

### 5.9 CodeQL Workflow Fix

**PR:** #11226

Fixed an issue with the CodeQL workflow action configuration that was preventing security scanning from running correctly.

### 5.10 Azure Tenant Migration

**PR:** #11269 (4 files, 127 insertions)

Migrated CI/CD infrastructure to a new Azure tenant with Terraform-managed resources, improving the fork-ability of the repository for external contributors.

### 5.11 Long-Running Test `--set` Flags

**PR:** #11219

Added `--set` flags to `rad upgrade kubernetes` in the long-running test workflow for proper test configuration.

---

## 6. Documentation & Developer Experience

### 6.1 Architecture Documentation

Two comprehensive architecture documents were added:

- **Deployment Engine Architecture** (PR #11385 — 649 lines): Covers the `bicep-de` container, deployment flow from template submission through resource creation, client integration, internal package structure, and debugging guidance. Includes Mermaid diagrams.
- **State Persistence Architecture** (PR #11284 — 659 lines): Documents the three pluggable state persistence subsystems (`database.Client`, `secret.Client`, `queue.Client`) with interface descriptions, implementation details, and data flow diagrams.

### 6.2 Copilot Skills

Three new Copilot agent skills were introduced:

- **Architecture Documenter** (PR #11281 — 352 insertions): A skill for generating architecture overviews, component diagrams, and sequence diagrams using Mermaid notation from codebase analysis.
- **Contributing Docs Updater** (PR #11401 — 311 insertions): A skill for maintaining contributor documentation with modes for rewrite, create, gap analysis, doc review, and code review impact assessment.
- **Issue Investigator Agent** (PR #11221 — 148 insertions): Ported from the Drasi repository for analyzing and providing technical context on GitHub issues.

### 6.3 Code Review Instructions

**PR:** #11213

Added PR title review guidance to the code review instructions, ensuring reviewers check that PR titles accurately describe the changes.

---

## 7. Dependency Updates

Dependabot automated **10 dependency update PRs** during this period:

| Category | Updates |
|----------|---------|
| Go dependencies | 4 PRs (~69 package updates across go.mod files) |
| GitHub Actions | 4 PRs (~20 action version bumps) |
| npm dependencies | 1 PR (17 package updates across 3 directories) |
| Security (go-git) | 1 PR (v5.16.4 → v5.16.5) |

---

## 8. Key Metrics

| Metric | Value |
|--------|-------|
| Total commits | 59 |
| Unique files changed | 123 |
| Lines added | ~17,966 |
| Lines removed | ~12,540 |
| PRs merged | ~50 |
| Release candidates | 8 (rc1–rc8) |
| Final release | v0.55.0 |
| Contributors | 11+ |

### Top Contributors (by commit count)

| Contributor | Commits |
|-------------|---------|
| Brooke Hamilton | 13 |
| Dariusz Porowski | 13 |
| dependabot[bot] | 10 |
| Copilot | 5 |
| Nicole James | 4 |
| Will Smith | 3 |
| Lakshmi Javadekar | 2 |
| Nithya Subramanian | 2 |
| Nell Shamrell-Harrington | 1 |
| Shruthi Kumar | 1 |

---

## 9. Looking Ahead

Based on the work completed this period, the following areas are expected to continue in the next cycle:

- **Sensitive data management** — Additional testing and potential schema extensions for the `x-radius-sensitive` annotation
- **Build infrastructure stabilization** — Continued refinement of the pnpm migration and bicep-types publishing workflows
- **Architecture documentation** — Expansion of the architecture documentation series with additional Mermaid diagrams
- **CI/CD optimization** — Further improvements to test reliability, fork compatibility, and workflow efficiency
- **Copilot agent capabilities** — Continued development of AI-assisted development tools and agent skills

---

*Document generated from git commit history for the `radius` repository, covering the period February 10, 2026 through March 10, 2026.*
