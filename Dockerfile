############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app/

# Fetch dependencies.
# Using go mod - only need go.mod and go.sum and need to cache based on these layers
COPY go.mod /app
COPY go.sum /app

RUN go mod download all
RUN go mod verify
COPY . /app

# Build the binary.
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -mod=mod -a -installsuffix -ldflags="-extldflags=-static" -tags timetzdata -o appbin ./cmd

############################
# STEP 2 build a small image
############################
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /app/appbin /app/appbin
# Use an unprivileged user.
USER appuser:appuser

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the hello binary.
ENTRYPOINT ["/app/appbin"]