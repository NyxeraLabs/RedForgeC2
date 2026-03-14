FROM node:20-alpine AS builder
WORKDIR /app/ui

ARG VITE_BUILD_SHA=dev
ARG VITE_BUILD_EPOCH=0
ARG VITE_TEAMSERVER_URL=https://localhost:9080

COPY ui/package.json ui/package-lock.json* ./
RUN npm ci --silent

COPY ui ./
ENV VITE_BUILD_SHA=$VITE_BUILD_SHA
ENV VITE_BUILD_EPOCH=$VITE_BUILD_EPOCH
ENV VITE_TEAMSERVER_URL=$VITE_TEAMSERVER_URL
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/ui/dist ./dist

RUN npm install -g serve

# Copy SSL certificates
COPY docker/certs/server.crt /app/server.crt
COPY docker/certs/server.key /app/server.key

EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "tcp://0.0.0.0:5173", "--ssl-cert", "/app/server.crt", "--ssl-key", "/app/server.key"]
