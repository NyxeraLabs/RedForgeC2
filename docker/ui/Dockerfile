FROM node:20-alpine AS builder
WORKDIR /app/ui

ARG VITE_BUILD_SHA=dev
ARG VITE_BUILD_EPOCH=0

COPY ui/package.json ui/package-lock.json* ./
RUN npm ci --silent

COPY ui ./
ENV VITE_BUILD_SHA=$VITE_BUILD_SHA
ENV VITE_BUILD_EPOCH=$VITE_BUILD_EPOCH
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/ui/dist ./dist

RUN npm install -g serve

EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "tcp://0.0.0.0:5173"]
