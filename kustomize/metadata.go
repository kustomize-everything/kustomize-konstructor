package kustomize

import (
)

// Functions for annotating the kustomizations with standard metadata

// Existing tool adds the following annotations:
// kustomize edit add annotation env-branch:"${ENV_BRANCH}"
// kustomize edit add annotation env-branch-url:"${ENV_BRANCH_URL}"
// kustomize edit add annotation deployment-repo:"${GITHUB_REPOSITORY}"
// kustomize edit add annotation deployment-repo-url:"${DEPLOY_REPO_URL}"
// kustomize edit add buildmetadata originAnnotations,managedByLabel
