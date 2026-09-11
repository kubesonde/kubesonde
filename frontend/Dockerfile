# Build the Vite app. Pin to the native build platform so Node never runs under
# QEMU emulation (arm64 emulation crashes npm with "signal 4 / Illegal
# instruction"). The output is static HTML/JS/CSS and is architecture-agnostic,
# so building it once serves every runtime platform; only the nginx stage below
# stays multi-arch.
FROM --platform=$BUILDPLATFORM node:24.20.0-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .

ARG VITE_API_SERVER=http://localhost:2709
ENV VITE_API_SERVER=$VITE_API_SERVER

RUN VITE_APP_VERSION="$(node -p "require('./package.json').version")" npm run build

# Serve the static SPA with nginx

FROM nginxinc/nginx-unprivileged:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/build /usr/share/nginx/html
# Run as an unprivileged user. Make nginx's writable dirs owned by nobody so it
# can create its temp/pid files without root.
USER root
RUN chown -R nobody /var/cache/nginx /tmp
USER nobody
EXPOSE 8088
CMD ["nginx", "-g", "daemon off;"]
