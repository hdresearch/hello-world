FROM node:22.14.0-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json build.mjs ./
COPY src ./src
RUN npm run build

FROM nginx:1.27.4-alpine3.21
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY require-env.sh /docker-entrypoint.d/05-require-env.sh
COPY --from=build /app/dist /usr/share/nginx/html
RUN chmod +x /docker-entrypoint.d/05-require-env.sh
EXPOSE 80
HEALTHCHECK --interval=10s --timeout=2s --retries=6 \
  CMD wget -qO- http://127.0.0.1/healthz >/dev/null || exit 1
