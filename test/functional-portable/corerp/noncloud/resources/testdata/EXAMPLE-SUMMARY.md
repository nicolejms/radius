# Application Graph Visualization Example - Summary

## Overview

This implementation demonstrates building a Radius application definition that uses:
1. **Kubernetes Container** - Running the magpie test application
2. **Redis Cache** - Provisioned via a Bicep recipe from the test recipe registry
3. **Volume** - An ephemeral volume with disk storage

The application definition generates a complete **application graph** that visualizes the relationships between these components.

## What Was Created

### 1. Application Definition (`corerp-resources-container-redis-volume.bicep`)

A complete Bicep template that defines:
- **Environment** with Kubernetes compute and Redis recipe configuration
- **Application** with a custom namespace
- **Container** with:
  - Image configuration (magpie test app)
  - Port 3000 exposed with HTTP readiness probe
  - Ephemeral volume mounted at `/var/cache` (disk-backed)
  - Connection to the Redis cache
- **Redis Cache** (extender resource) provisioned via the `corerp-redis-recipe`

### 2. Test Implementation (`application_test.go`)

Added `Test_ApplicationGraphWithRedisAndVolume` which:
- Deploys the application using the Bicep template
- Verifies all resources are created (environment, app, container, Redis)
- Validates Kubernetes objects (pods, services, etc.)
- Tests the application graph API
- Verifies connections between resources
- Logs detailed graph structure for debugging

### 3. Documentation Files

#### README (`corerp-resources-container-redis-volume-README.md`)
- Architecture diagram showing component relationships
- Deployment instructions
- Graph visualization explanation
- Testing details

#### Expected Graph Output (`corerp-resources-container-redis-volume-graph.json`)
JSON representation showing:
- Two main resources (container and Redis cache)
- Outbound connection from container to Redis
- Inbound connection to Redis from container
- Generated Kubernetes resources (deployments, services, RBAC)

#### CLI Output Example (`corerp-resources-container-redis-volume-graph-output.txt`)
Text visualization showing how `rad app graph` displays the relationships:
```
Name: redis-app-container (Applications.Core/containers)
Connections:
  redis-app-container -> redis-cache (Applications.Core/extenders)
Resources:
  redis-app-container (apps/Deployment)
  [... other Kubernetes resources ...]

Name: redis-cache (Applications.Core/extenders)
Connections:
  redis-app-container (Applications.Core/containers) -> redis-cache
Resources:
  redis-cache (apps/Deployment)
  [... other Kubernetes resources ...]
```

#### Usage Script (`example-usage.sh`)
Step-by-step demonstration script showing:
1. How to deploy the application
2. How to verify the deployment
3. How to visualize the application graph
4. How to inspect Kubernetes resources
5. How to view resource connections
6. Cleanup instructions

## Key Features Demonstrated

### 1. Recipe Integration
- Environment configured with a Bicep recipe for Redis
- Recipe sourced from the test recipe registry
- Parameters passed to the recipe (redisName)

### 2. Resource Connections
- Container explicitly connects to Redis via `connections` property
- Graph shows both outbound (from container) and inbound (to Redis) connections
- Demonstrates service-to-service relationships

### 3. Volume Management
- Ephemeral volume with disk backing
- Mounted at `/var/cache` in the container
- Shows how Radius handles storage in containerized applications

### 4. Application Graph Visualization
- Graph API returns structured data about resources and connections
- CLI command (`rad app graph`) provides human-readable visualization
- Test validates graph structure programmatically

## How to Use

### Prerequisites
- Radius installed and configured
- Kubernetes cluster available
- Access to the test recipe registry

### Deployment
```bash
cd test/functional-portable/corerp/noncloud/resources/testdata/
rad deploy corerp-resources-container-redis-volume.bicep \
  --parameters magpieimage=<image> \
  --parameters registry=<recipe-registry> \
  --parameters version=<recipe-version>
```

### Visualization
```bash
# View the application graph
rad app graph corerp-app-redis-volume

# Expected output shows:
# - Container with outbound connection to Redis
# - Redis cache receiving inbound connection
# - All generated Kubernetes resources
```

### Testing
```bash
# Run the functional test
cd test/functional-portable/corerp/noncloud/resources/
go test -run Test_ApplicationGraphWithRedisAndVolume -v
```

## Architecture

```
┌────────────────────────────────────────┐
│  Radius Application                    │
│  ┌──────────────┐    Connection        │
│  │  Container   │ ─────────────┐       │
│  │  + Volume    │              ▼       │
│  └──────────────┘    ┌──────────────┐  │
│                      │  Redis Cache │  │
│                      │  (via Recipe)│  │
│                      └──────────────┘  │
└────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────┐
│  Generated Kubernetes Resources        │
│  • Deployments                         │
│  • Services                            │
│  • ServiceAccounts                     │
│  • RBAC (Roles, RoleBindings)         │
└────────────────────────────────────────┘
```

## Testing Strategy

The test validates:
1. **Resource Creation**: All expected resources are created
2. **Kubernetes Objects**: Pods and services are deployed correctly
3. **Graph Structure**: Application graph API returns correct data
4. **Connections**: Container-to-Redis connection exists
5. **Resource Presence**: Both container and Redis appear in the graph

## Value Demonstrated

This example showcases how Radius:
1. **Simplifies Infrastructure**: Recipes abstract away Redis deployment complexity
2. **Visualizes Dependencies**: Graph shows how services relate to each other
3. **Manages Kubernetes**: Automatically generates and manages K8s resources
4. **Connects Services**: Declaratively expresses service dependencies
5. **Integrates Storage**: Seamlessly adds volumes to containerized applications

## Future Enhancements

Potential additions to this example:
- Add more complex volume types (persistent volumes, ConfigMaps)
- Include multiple containers with different connection patterns
- Add HTTP routes and gateways
- Demonstrate recipe parameters and outputs
- Show environment variable injection from recipe outputs
