extension radius

@description('Specifies the location for resources.')
param location string = 'local'

@description('Specifies the image for the container resource.')
param magpieimage string

@description('The OCI registry for test Bicep recipes.')
param registry string

@description('The OCI tag for test Bicep recipes.')
param version string

// Environment with Redis recipe configured
resource env 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'corerp-app-redis-volume-env'
  location: location
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'corerp-app-redis-volume-env'
    }
    recipes: {
      'Applications.Core/extenders': {
        default: {
          templateKind: 'bicep'
          templatePath: '${registry}/test/testrecipes/test-bicep-recipes/corerp-redis-recipe:${version}'
        }
      }
    }
  }
}

resource app 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'corerp-app-redis-volume'
  location: location
  properties: {
    environment: env.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'corerp-app-redis-volume-ns'
      }
    ]
  }
}

// Container with ephemeral volume
resource container 'Applications.Core/containers@2023-10-01-preview' = {
  name: 'redis-app-container'
  location: location
  properties: {
    application: app.id
    container: {
      image: magpieimage
      ports: {
        web: {
          containerPort: 3000
        }
      }
      readinessProbe: {
        kind: 'httpGet'
        containerPort: 3000
        path: '/healthz'
      }
      volumes: {
        cache: {
          kind: 'ephemeral'
          managedStore: 'disk'
          mountPath: '/var/cache'
        }
      }
    }
    connections: {
      redis: {
        source: redisCache.id
      }
    }
  }
}

// Redis cache using recipe
resource redisCache 'Applications.Core/extenders@2023-10-01-preview' = {
  name: 'redis-cache'
  location: location
  properties: {
    application: app.id
    environment: env.id
    recipe: {
      name: 'default'
      parameters: {
        redisName: 'myredis'
      }
    }
  }
}
