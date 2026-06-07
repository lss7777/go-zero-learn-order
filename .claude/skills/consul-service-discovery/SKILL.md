---
name: consul-service-discovery
description: This skill should be used when the user asks to "add service discovery", "integrate consul", "register service with consul", "use consul for service discovery", "configure consul registry", or needs to connect go-zero services via Consul instead of hardcoded addresses. Covers both server-side registration (RPC services) and client-side discovery (API/RPC services calling upstream services).
version: 1.0.0
allowed-tools:
  - Read
  - Edit
  - Bash
  - Glob
---

# go-zero Consul Service Discovery & Registration

This skill covers integrating Consul as the service registry in go-zero microservices. It addresses two distinct roles:

- **Client-side discovery**: An API or RPC service that needs to *call* another service via Consul (most common)
- **Server-side registration**: An RPC service that *registers itself* with Consul so others can discover it

## Package

```
github.com/zeromicro/zero-contrib/zrpc/registry/consul
```

**Import alias**: Always use blank import `_ "..."` to trigger the `init()` function that registers the `consul://` gRPC resolver.

## Client-Side Discovery (Call Another Service via Consul)

This is the minimal approach — your service discovers upstream RPC services through Consul without registering itself.

### Steps

**1. Add dependency**

```bash
go get github.com/zeromicro/zero-contrib/zrpc/registry/consul
```

**2. Blank import in `main.go`** (or any file that runs before gRPC client init)

```go
import (
    _ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)
```

This import triggers `init()` in the consul package, which registers a gRPC resolver for the `consul://` URL scheme. Without it, gRPC won't understand `consul://` targets.

**3. Change RPC client target in `etc/*.yaml`**

```yaml
# Before (hardcoded address)
UserRpcConf:
  Target: 127.0.0.1:8080

# After (Consul discovery)
UserRpcConf:
  Target: consul://127.0.0.1:8500/user-rpc?wait=14s
```

### URL Format

```
consul://[user:passwd@]<consul-host>:<port>/<service-key>?<params>
```

| Parameter | Required | Description | Example |
|---|---|---|---|
| `consul-host:port` | Yes | Consul agent address | `127.0.0.1:8500` |
| `service-key` | Yes | Service name registered in Consul | `user-rpc` or `user.UserService` |
| `wait` | No | Long-polling wait time for Consul blocking queries | `14s` |
| `token` | No | Consul ACL token | `token=abc123` |
| `tag` | No | Filter by service tag | `tag=public` |
| `dc` | No | Datacenter | `dc=us-east` |

**With ACL token:**
```yaml
Target: consul://127.0.0.1:8500/user-rpc?wait=14s&token=f0512db6-76d6-f25e-f344-a98cc3484d42
```

**4. No changes needed in `config.go` or `servicecontext.go`**

The `zrpc.RpcClientConf` already handles target resolution. The blank import is sufficient.

---

## Server-Side Registration (Register RPC Service with Consul)

Use this when an RPC service needs to register itself so other services can discover it.

> **Note**: REST API services typically do NOT register themselves with Consul. Only gRPC/RPC services need registration.

### Steps

**1. Update `internal/config/config.go`**

```go
import (
    "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
    "github.com/zeromicro/zrpc"
)

type Config struct {
    zrpc.RpcServerConf
    Consul consul.Conf
}
```

**2. Add Consul section to `etc/*.yaml`**

```yaml
Consul:
  Host: 127.0.0.1:8500          # Consul agent address
  Key: user-rpc                  # Service name for discovery
  Token: ""                     # ACL token (optional)
  TTL: 20                        # Health check TTL in seconds (default: 20)
  Tag:                           # Service tags (optional)
    - rpc
    - v1
  Meta:                          # Metadata (optional)
    Protocol: grpc
```

**3. Register in `main.go`**

```go
import (
    _ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
    "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

func main() {
    // ... load config ...

    server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
        // register gRPC service handlers
    })

    // Register this service with Consul
    if err := consul.RegisterService(c.ListenOn, c.Consul); err != nil {
        panic(err)
    }

    server.Start()
}
```

The `RegisterService` call:
- Registers the service with Consul agent
- Creates a TTL-based health check
- Starts a background goroutine to periodically update TTL (`TTL - 1` second interval)
- Registers a shutdown listener to deregister on graceful stop

### `consul.Conf` Fields

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `Host` | `string` | Yes | — | Consul agent address (e.g., `127.0.0.1:8500`) |
| `Key` | `string` | Yes | — | Service name registered in Consul |
| `Token` | `string` | No | `""` | Consul ACL token |
| `Tag` | `[]string` | No | `nil` | Service tags |
| `Meta` | `map[string]string` | No | `nil` | Service metadata (e.g., `Protocol: grpc`) |
| `TTL` | `int` | No | `20` | Health check TTL in seconds |

---

## Key Mechanism

The blank import `_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"` triggers the `init()` function in `builder.go`:

```go
func init() {
    resolver.Register(&builder{})
}
```

This registers a custom gRPC resolver for the `consul` scheme. When gRPC encounters a target like `consul://127.0.0.1:8500/user-rpc`, it routes to the consul resolver which:
1. Parses the URL to extract service name and parameters
2. Watches the Consul Health API for healthy service instances
3. Populates gRPC endpoints dynamically

## Common Issues

**`field "xxx" is not set` on startup**
- The consul import must be a blank import in `main.go` (or the entry point package)
- Ensure `go mod tidy` has been run after adding the dependency

**Service discovery returns no endpoints**
- Verify the service is registered and healthy in Consul UI (`http://127.0.0.1:8500`)
- Check the `Key` in server-side config matches the service name in client-side target URL
- Ensure the service's health check is passing (TTL updates must succeed)

**ACL errors**
- Server (registration): token needs `service "xxx" { policy = "write" }`
- Client (discovery): token only needs `service "xxx" { policy = "read" }` and `node "consul-server" { policy = "read" }`
