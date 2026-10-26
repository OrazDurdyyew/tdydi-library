# syntax=docker/dockerfile:1.6

# ===== Build =====
FROM node:22-alpine AS builder

WORKDIR /build

# pnpm через corepack (лочка pnpm-lock.yaml уже в репозитории)
RUN corepack enable && corepack prepare pnpm@10.26.0 --activate

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# ===== Runtime =====
FROM node:22-alpine

RUN apk add --no-cache wget \
 && addgroup -S library && adduser -S -G library -h /app library

WORKDIR /app
COPY --from=builder --chown=library:library /build/.output /app/.output

USER library

# Версия для отображения в админке. Передаётся build-arg'ом из make-offline-bundle.sh
# и зашивается в ENV образа → Nuxt читает её как public.appVersion на старте.
ARG VERSION=dev
ENV NUXT_PUBLIC_APP_VERSION=${VERSION}

ENV HOST=0.0.0.0
ENV PORT=2026
ENV NODE_ENV=production

EXPOSE 2026

HEALTHCHECK --interval=20s --timeout=5s --start-period=15s --retries=3 \
  CMD wget -qO- http://127.0.0.1:2026/ > /dev/null 2>&1 || exit 1

CMD ["node", ".output/server/index.mjs"]
