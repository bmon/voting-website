# Accept the Go version for the image to be set as a build argument.
ARG GO_VERSION=1.13

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.* ./
RUN go mod download

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app ./cmd/app/

# Final stage: the running container.
FROM scratch AS final

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# Import the root ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
ENV PORT 8080

# Run the compiled binary.
ENTRYPOINT ["/app"]
