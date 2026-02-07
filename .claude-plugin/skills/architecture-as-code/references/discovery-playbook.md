---
tags: [discovery, exploration, grep, tracing, agents]
related: [bootstrap.md, templates.md]
---

# Discovery Playbook

Reusable exploration techniques for mapping system architecture. Used during both initial bootstrap and ongoing maintenance.

## Service Boundary Discovery

Map each service's role, entry points, and external interfaces.

| What to Find | Search Strategy |
|--------------|-----------------|
| Handler entry points | `grep -rn "func.*Handler\|func.*handler\|func Handle" --include="*.go"` |
| RPC/IDL definitions | `find . -name "*.thrift" -o -name "*.proto"` |
| HTTP endpoints | `grep -rn "router\.\|mux\.\|Handle\(\|HandleFunc\(" --include="*.go"` |
| gRPC services | `grep -rn "RegisterServer\|pb.Register" --include="*.go"` |
| Service dependencies (outbound) | `grep -rn "rpc\.\|http\.\|grpc\." dal/ --include="*.go"` |
| DI/wiring config | `find . -name "wire.go" -o -name "providers.go" -o -name "inject.go"` |
| Main entry point | `find . -name "main.go" -path "*/cmd/*"` |

**Language-specific variants**: Adapt patterns for non-Go repos. For Python: `def.*handler`, `@app.route`. For Java: `@Controller`, `@Service`, `@Repository`. For Node.js: `router.get`, `app.use`.

## Storage Tracing

Map all storage systems and who reads/writes them.

| Storage Type | Discovery Pattern | Key Output |
|-------------|-------------------|------------|
| Redis/KV | Grep for `Set`, `Get`, `HMSet`, `HGetAll` + key prefix patterns | Reader/writer matrix with key prefixes |
| MongoDB/Document | Find collection names, `Find`, `Insert`, `Aggregate` calls | Collection ownership map |
| Object Storage | Grep for bucket names, `Upload`, `Download`, `PutObject` | Bucket-to-service mapping |
| Message Queue | Grep for topic names, `Publish`, `Subscribe`, `Produce`, `Consume` | Topic ownership + consumer list |
| Elasticsearch | Find index patterns, `Search`, `Index`, `Bulk` calls | Index-to-service mapping |
| SQL databases | Find table names, migrations, ORM models | Table ownership map |

**Key output**: A reader/writer matrix showing which service reads/writes each storage resource. This is the foundation for "if I change X, what breaks?" questions.

Example matrix format:
```
| Storage Resource | Writers | Readers |
|-----------------|---------|---------|
| kv:prefix:factors | service-a | service-b, service-c |
| mongodb:scan_results | service-a | frontend-bff |
| tos:screenshots/ | sandbox-service | frontend-bff |
```

## Data Flow Tracing

Trace data from entry to output through the full call chain.

### Technique: Follow the Call Chain

1. Start at the handler entry point
2. Follow: handler -> service -> operator -> dal
3. At each dal call, note: what data is written, where, with what key
4. At each incoming read, note: what data is read, from where, by what key
5. Map the contract between services: "service A writes X, service B reads X"

### Technique: Trace a Specific Field

When you need to understand how a specific data field flows:

1. Find where the field is first set (origin)
2. Grep for the field name across all repos
3. Map: origin -> persistence -> consumers
4. Document in a data flow diagram:
   ```
   field_name: service-a (creates) -> storage-key (persists) -> service-b (reads)
   ```

### Technique: Contract Discovery

Identify implicit contracts between services:

| Contract Type | How to Identify |
|---------------|-----------------|
| Shared storage key | Same key prefix/collection in multiple repos |
| RPC contract | Shared Thrift/Proto struct referenced by caller and callee |
| HTTP contract | Request/response structs matching across services |
| Event contract | Message queue topic with multiple producers or consumers |
| Config contract | Same config key read by multiple services |

Key question: "If I change X in repo A, what breaks in repo B?"

## External Config Discovery

Identify external systems that affect runtime behavior.

| Config Type | Discovery | Impact |
|-------------|-----------|--------|
| Rule engines | Grep for SDK calls (Dolphin, Drools, decision tables) | Changes take effect without deploy |
| LLM prompts | Grep for prompt keys, model API calls | Output changes without code changes |
| Feature flags | Grep for config center reads (TCC, LaunchDarkly, etc.) | Behavior changes without deploy |
| Whitelists/blocklists | Grep for allowlist/denylist reads from external source | Filtering changes without deploy |

For each external config:
- Document the config-to-code mapping (which config key controls which behavior)
- Determine if config changes need deployment or take effect immediately
- Decide if snapshots should be stored in git (see `references/maintenance.md` > Config Snapshots)

## Explore Agent Prompts

Pre-built prompts for parallel discovery with Explore agents.

### Service Mapping Agent

```
Explore {repo_path}. Map the service architecture:
1. Handler entry points (file:line for each)
2. Service layer orchestration (what calls what, in what order)
3. Operator/business logic (key algorithms and decisions)
4. DAL/external calls (what external systems are called, with what data)
5. DI wiring (how dependencies are injected)
Output: service map with key files, function names, and line numbers.
```

### Storage Mapping Agent

```
In {repo_path}, find ALL storage operations:
- KV/Redis: key prefixes, Set/Get/HMSet patterns
- Document DB: collection names, query patterns
- Object storage: bucket names, upload/download paths
- Message queues: topic names, publish/subscribe
- Search indexes: index names, query/index operations
For each: note the key/prefix/collection/bucket/topic, operation type
(read/write/both), and file:line. Output a storage matrix.
```

### Contract Discovery Agent

```
Compare {repo_A} and {repo_B} for shared contracts:
- Same storage key prefixes used by both
- Shared IDL/Thrift types referenced by both
- Same config keys read by both
- Any RPC calls between them (find caller and callee)
Output: contract map showing what's shared, which side owns it,
and what would break if either side changed independently.
```

### Config Mapping Agent

```
In {repo_path}, find all references to external configuration:
- Rule engine SDK calls (identify rule groups and field references)
- LLM/AI prompt keys (identify prompt identifiers and where results are used)
- Feature flag reads (identify flag names and controlled behaviors)
- Config center reads (identify config keys and their consumers)
For each: document the config key, the code that reads it, and the
behavior it controls. Output a config-to-code mapping table.
```
