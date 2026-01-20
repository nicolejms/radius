#!/bin/bash
# Example script showing how to deploy and visualize the application graph
# This script is for demonstration purposes and requires Radius to be installed and configured

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}Radius Application Graph Example${NC}"
echo -e "${BLUE}Container + Redis Cache + Volume${NC}"
echo -e "${BLUE}================================================${NC}"
echo

# Configuration
APP_NAME="corerp-app-redis-volume"
BICEP_FILE="corerp-resources-container-redis-volume.bicep"
NAMESPACE="corerp-app-redis-volume-ns"

# Check if rad CLI is installed
if ! command -v rad &> /dev/null; then
    echo -e "${YELLOW}Error: rad CLI not found. Please install Radius first.${NC}"
    echo "Visit: https://docs.radapp.io/getting-started/"
    exit 1
fi

echo -e "${GREEN}Step 1: Deploy the application${NC}"
echo "Deploying Bicep template: ${BICEP_FILE}"
echo

# Deploy the application
# Note: Replace these parameters with actual values
rad deploy "${BICEP_FILE}" \
    --parameters magpieimage="ghcr.io/radius-project/magpiego:latest" \
    --parameters registry="ghcr.io/radius-project/dev" \
    --parameters version="latest"

echo
echo -e "${GREEN}Step 2: Verify the deployment${NC}"
echo "Checking application status..."
echo

# List the application
rad app show "${APP_NAME}"

echo
echo -e "${GREEN}Step 3: Visualize the application graph${NC}"
echo "Generating application graph..."
echo

# Show the application graph
rad app graph "${APP_NAME}"

echo
echo -e "${GREEN}Step 4: Inspect Kubernetes resources${NC}"
echo "Listing resources in namespace: ${NAMESPACE}"
echo

# List pods in the namespace
kubectl get pods -n "${NAMESPACE}"

echo
echo -e "${GREEN}Step 5: Check connections${NC}"
echo "Viewing resource details..."
echo

# Show container details
rad resource show containers redis-app-container -a "${APP_NAME}"

echo
echo -e "${GREEN}Step 6: Inspect the Redis cache${NC}"
echo "Viewing Redis cache details..."
echo

# Show Redis cache details
rad resource show extenders redis-cache -a "${APP_NAME}"

echo
echo -e "${BLUE}================================================${NC}"
echo -e "${GREEN}✓ Application graph example complete!${NC}"
echo -e "${BLUE}================================================${NC}"
echo
echo "The application includes:"
echo "  • A Kubernetes container (redis-app-container)"
echo "  • A Redis cache (redis-cache) provisioned via recipe"
echo "  • An ephemeral volume mounted at /var/cache"
echo "  • A connection from the container to the Redis cache"
echo
echo "To clean up:"
echo "  rad app delete ${APP_NAME}"
