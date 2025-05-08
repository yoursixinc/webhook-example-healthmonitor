FROM node:22 AS frontend-builder

# Build frontend
WORKDIR /build
COPY package.json .
COPY package-lock.json .
COPY index.html .
COPY src src
COPY public public
RUN npm install && npx vite build

FROM golang:1.24-bookworm AS backend-builder

# Build backend
WORKDIR /build
COPY --from=frontend-builder /build/dist ./dist
COPY go.mod .
COPY go.sum .
COPY main.go .
COPY internal internal
RUN go build

FROM debian:bookworm-slim

# Add run user
RUN groupadd -r app && \
  useradd -r -s /bin/false -g app app

# Install dependencies
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates tzdata && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* && \
  update-ca-certificates

# Install applications
USER app
WORKDIR /app
COPY --from=backend-builder /build/dist ./dist
COPY --from=backend-builder /build/healthmonitor .

# Run
CMD ["./healthmonitor"]
