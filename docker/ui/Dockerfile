FROM node:20-alpine AS builder
WORKDIR /app/ui

COPY ui/package.json ui/package-lock.json* ./
RUN npm ci --silent

COPY ui ./
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/ui/dist ./dist

RUN npm install -g serve

EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "tcp://0.0.0.0:5173"]
