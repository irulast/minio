# MinIO maintenance

## Ownership and provenance

The source base is archived upstream commit `7aac2a2c5b7c882e68c1ce017d8256be2feea27f`.
Security backports come from the AGPL-licensed PGSTY source, identified by full
commit in `security-backports.json`. Copyright and license notices are retained.
We own patch selection, builds, regression validation, release approval, and
GitOps promotion. There is no automatic tracking of PGSTY's server releases.
The MinIO module name, binary, endpoint names, and operator contracts are retained.

The patched `github.com/pgsty/silo-pkg/v3` v3.11.0 dependency supplies IAM
condition-key and bucket-resource matching fixes and the LDAP TLS correction,
using Go's replacement mechanism to preserve historical import paths. The server
and console are not switched to the Silo product.

## Candidate security scope

The upstream base includes CVE-2024-55949 and CVE-2025-27414, CVE-2025-31489,
and CVE-2025-62506. Backports address CVE-2026-33322, CVE-2026-33419,
CVE-2026-34204, CVE-2026-39414, CVE-2026-41145, CVE-2026-40344,
CVE-2026-42600, and CVE-2026-97731. Internode path/metadata guards and IAM
condition hardening include the documented follow-up corrections. This is a
patch inventory, not a claim that a candidate has passed validation or has no
other vulnerabilities. The dependency graph uses Go 1.27.3 and updated libraries;
the release must retain fresh govulncheck and full image scan results.

OIDC STS now requires JWKS-backed asymmetric signing. LDAP throttling is
per-node and uses the peer address unless trusted proxies are configured.
S3 clients must sign operation-bearing x-amz headers. IAM policies granting
sensitive bucket writes need the bare bucket ARN as well as object ARNs.
The removed ReadMultiple endpoint has no production caller in the deployed
2024-11 source; mixed-version operation still requires an upgrade rehearsal.

## Build and promotion

The private CI repository accepts an explicit 40-character source commit and
builds it on `mke-builds`. It runs ordinary storage/serialization tests, a compile
check of all test packages, the upstream blocking govulncheck, and the existing
shared BuildKit signing, SBOM, and Trivy gates. Scan thresholds, base-layer
handling, and waivers are not changed. Candidate tags start with `security-` and
cannot match production release policies.

Before release: review the diff, pass CI and fresh image scans, export config
and IAM securely, preserve the exact old image, and rehearse old-to-new startup
and reads/writes against isolated restored storage. Verify ordinary, versioned,
large/multipart, and encrypted objects as applicable. Record test limitations.
Backups or snapshots without a successful restore are not recovery evidence.

Promote only the tested digest through the bare-metal-k8s-setup GitOps repository.
Upgrade the shared backup store and Harbor store in separate observed steps.
Pre-pull both candidate and recovery images before touching Harbor's own store.
Do not assume a binary downgrade is safe after storage metadata changes; use the
verified recovery procedure. Keep production credentials out of CI and logs.

## Ongoing maintenance

At each explicitly requested security maintenance, check vendor advisories,
Go vulnerability data, dependency and base-image fixes, and relevant reviewed
community patches. Record CVEs, source commits, scans, functional checks, and
rollout evidence. There is no new recurring scan or suppression policy.
