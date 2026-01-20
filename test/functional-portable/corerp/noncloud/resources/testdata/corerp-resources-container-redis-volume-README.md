# Application Graph Example: Container with Redis Cache and Volume

This example demonstrates creating a Radius application that includes:
- A Kubernetes container
- A Redis cache (provisioned via a Bicep recipe)
- An ephemeral volume

The example showcases how Radius creates an application graph that visualizes the relationships between these resources.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│  Application: corerp-app-redis-volume               │
│  Namespace: corerp-app-redis-volume-ns              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────────────────────────────────┐       │
│  │  Container: redis-app-container         │       │
│  │  ────────────────────────────────────   │       │
│  │  • Image: magpie                        │       │
│  │  • Port: 3000                           │       │
│  │  • Readiness: /healthz                  │       │
│  │  • Volume: /var/cache (ephemeral/disk) │       │
│  └────────────┬────────────────────────────┘       │
│               │                                     │
│               │ Connection                          │
│               │ (outbound)                          │
│               ▼                                     │
│  ┌─────────────────────────────────────────┐       │
│  │  Redis Cache: redis-cache               │       │
│  │  ────────────────────────────────────   │       │
│  │  • Provisioned via Recipe               │       │
│  │  • Recipe: corerp-redis-recipe          │       │
│  │  • Parameter: redisName=myredis         │       │
│  └─────────────────────────────────────────┘       │
│                                                     │
└─────────────────────────────────────────────────────┘

Generated Kubernetes Resources:
├── Container Resources:
│   ├── Deployment/redis-app-container
│   ├── Service/redis-app-container
│   ├── ServiceAccount/redis-app-container
│   ├── Role/redis-app-container
│   └── RoleBinding/redis-app-container
└── Redis Resources (from recipe):
    ├── Deployment/redis-cache
    └── Service/redis-cache
```

## Files

- **corerp-resources-container-redis-volume.bicep**: The main Bicep template that defines the application
- **Test_ApplicationGraphWithRedisAndVolume**: Go test function in `application_test.go`

## Application Components

### 1. Environment
The environment is configured with:
- Kubernetes compute (`kind: 'kubernetes'`)
- A Redis recipe from the test recipe registry (`corerp-redis-recipe`)

### 2. Application
The application (`corerp-app-redis-volume`) includes:
- A custom Kubernetes namespace (`corerp-app-redis-volume-ns`)

### 3. Container
The container (`redis-app-container`) features:
- A container image (magpie test image)
- Port 3000 exposed
- HTTP readiness probe on `/healthz`
- An ephemeral volume mounted at `/var/cache` (using disk storage)
- A connection to the Redis cache

### 4. Redis Cache
The Redis cache (`redis-cache`) is provisioned via a recipe:
- Uses the `corerp-redis-recipe` from the test recipe registry
- Configured with a custom Redis name parameter

## How to Deploy (when Radius is installed)

```bash
# Deploy the application
rad deploy test/functional-portable/corerp/noncloud/resources/testdata/corerp-resources-container-redis-volume.bicep \
  --parameters magpieimage=<your-magpie-image> \
  --parameters registry=<recipe-registry> \
  --parameters version=<recipe-version>
```

## Visualizing the Application Graph

Once deployed, you can visualize the application graph using:

```bash
# Show the application graph
rad app graph corerp-app-redis-volume
```

This will display:
1. **redis-app-container** (Applications.Core/containers)
   - Connections: Outbound connection to redis-cache
   - Resources: Kubernetes Deployment, Service, ServiceAccount, Role, RoleBinding

2. **redis-cache** (Applications.Core/extenders)
   - Resources: Provisioned by the Bicep recipe (e.g., Redis deployment, service)

## Expected Graph Structure

The application graph shows:
- The container has an **outbound connection** to the Redis cache
- Both resources generate Kubernetes resources (deployments, services, etc.)
- The volume is part of the container's configuration (not a separate graph node)

## Testing

The test validates:
- All resources are created successfully
- The application graph API returns the expected resources
- The container has a connection to Redis
- Both the container and Redis cache appear in the graph
