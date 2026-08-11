# Reusable GitHub Actions CI/CD Pipeline

Three files, meant to sit in a repo exactly at these paths:

```
.github/
├── actions/setup-env/action.yml     # composite action: install Node + cache deps
└── workflows/
    ├── reusable-ci-cd.yml           # the reusable pipeline (workflow_call)
    └── caller-example.yml           # example: how another workflow invokes it
```

## What it does

`reusable-ci-cd.yml` is not triggered directly — it exposes a `workflow_call`
interface and runs six jobs in dependency order:

1. **lint** — ESLint + Prettier, uploads a report artifact even on failure.
2. **test** — matrix across Node 18/20/22 and OS, unit + integration tests,
   coverage upload, Codecov.
3. **security** — `npm audit`, CodeQL static analysis, Trivy filesystem scan,
   results published to the repo's Security tab.
4. **build** — multi-arch (amd64/arm64) Docker build via Buildx, GHA layer
   caching, SBOM + provenance attestation, pushes to GHCR, then scans the
   built image with Trivy and fails the job on critical/high CVEs.
5. **deploy** — gated by `run-deploy` input and a GitHub **Environment**
   (configure required reviewers / wait timers there, not in YAML), OIDC
   federation to AWS (no long-lived secrets), ECS rolling deploy, smoke test,
   automatic rollback on failure.
6. **notify** — Slack status message, runs regardless of upstream outcome.

## Design choices worth keeping if you adapt this

- **`permissions: contents: read`** at the top, escalated per-job only where
  needed (`packages: write` for build, `id-token: write` for OIDC).
- **OIDC over static AWS keys** (`aws-actions/configure-aws-credentials` with
  `role-to-assume`) — nothing long-lived to leak.
- **`concurrency`** cancels superseded staging runs but never cancels an
  in-flight production deploy.
- **Environment protection rules** (manual approval for prod) live in repo
  Settings → Environments, not in the workflow — keeps the gate enforceable
  even if someone edits the YAML.
- **Composite action** for setup avoids duplicating cache/install logic
  across five jobs.

## Adapting to your stack

- Swap `npm ci` / `npm run test:*` in `setup-env/action.yml` and the `test`
  job for your language's equivalents.
- Replace the ECS `aws ecs update-service` block in `deploy` with your
  target (Kubernetes `kubectl set image`, Terraform, Helm, etc.).
- Required repo secrets: `REGISTRY_USERNAME`, `REGISTRY_PASSWORD`,
  `AWS_ROLE_ARN` (and `AWS_ROLE_ARN_PROD` if you split roles per env),
  `SLACK_WEBHOOK_URL` (optional).

## Reusing across repos

Because this lives in `workflow_call`, another repository can call it too —
just point the `uses:` at the full ref instead of a relative path:

```yaml
uses: your-org/your-repo/.github/workflows/reusable-ci-cd.yml@main
```
