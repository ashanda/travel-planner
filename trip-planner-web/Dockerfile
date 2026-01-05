# ---------- BUILD ----------
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build


# ---------- RUN ----------
FROM node:20-alpine

WORKDIR /app
ENV NODE_ENV=production

# Copy output build
COPY --from=builder /app/.output ./.output

# Nuxt public runtime variable (overridden in docker-compose)
ENV NUXT_PUBLIC_API_BASE=/api

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
