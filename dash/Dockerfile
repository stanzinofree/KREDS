# Get Go image from DockerHub.
FROM golang:1.21 AS build

# Set working directory.
WORKDIR /

# Copy dependency locks so we can cache.
COPY go.mod go.sum main.go ./

# Get all of our dependencies.
RUN go mod download

# Build our application.
RUN CGO_ENABLED=0 GOOS=linux go build -o dash ./main.go

# Use 'scratch' image for super-mini build.
#FROM golang:1.21 AS dev
FROM scratch AS dev

# Set working directory for this stage.
WORKDIR /

# Copy our compiled executable from the last stage.
COPY --from=build /dash .
COPY ./.env ./

# Run application and expose port 8080.
EXPOSE 8080
CMD ["./dash"]