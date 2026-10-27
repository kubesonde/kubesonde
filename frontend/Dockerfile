# Build the Vite app. Pin to the native build platform so Node never runs under
# QEMU emulation (arm64 emulation crashes npm with "signal 4 / Illegal
# instruction"). The output is static HTML/JS/CSS and is architecture-agnostic,
# so building it once serves every runtime platform; only the nginx stage below
# stays multi-arch.
FROM --platform=$BUILDPLATFORM node:26.12.0-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .

ARG VITE_API_SERVER=http://localhost:2709
ENV VITE_API_SERVER=$VITE_API_SERVER

# The release version comes from the repo-root VERSION file (the single source
# of truth that the controller release and both frontend deploy paths read),
# injected by the workflows as a build arg. Local builds show "dev".
# package.json's version field is not used and is not bumped per release.
ARG VERSION=dev
RUN VITE_APP_VERSION="$VERSION" npm run build

# Serve the static SPA with nginx

FROM nginxinc/nginx-unprivileged:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/build /usr/share/nginx/html
# Run as an unprivileged user. Make nginx's writable dirs owned by nobody so it
# can create its temp/pid files without root.
USER root
# Pull patched OS packages so fixable CVEs in the base image are resolved at
# build time (the :alpine tag floats to the latest nginx, but not always to the
# latest alpine package patches). Then set up nginx's writable dirs.
RUN apk upgrade --no-cache && chown -R nobody /var/cache/nginx /tmp
USER nobody
EXPOSE 8088
CMD ["nginx", "-g", "daemon off;"]
