set -x

VERSION=$(npm version patch)
git add .
echo "Bump version to ${VERSION}" | git commit -F -
git tag "${VERSION}"
git push origin "${VERSION}"
git push origin dev
git checkout main
git merge origin/dev
git push origin main
git checkout dev
