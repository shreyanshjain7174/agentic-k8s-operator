# OSS Scope

`agentic-operator-core` is the open-source foundation for agent workload orchestration.

## In Scope (OSS)

- CRD and controller baseline behavior
- Protocol integrations and abstractions (`pkg/mcp`, `pkg/opa`, `pkg/argo`)
- Reference Helm charts and deployment paths
- Public API stability and compatibility commitments

## Out of Scope (Private)

- Billing and monetization logic
- Multi-tenant SLA enforcement
- Enterprise autoscaling/routing optimization
- Advanced evaluation pipelines tied to paid tiers
- Commercial licensing enforcement

## Contribution Rule

If a feature is directly tied to paid tier differentiation or enterprise-only outcomes, open it in `agentic-operator-private` instead of this repository.
