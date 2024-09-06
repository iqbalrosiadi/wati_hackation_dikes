# Stage 1: Build stage
FROM golang:1.23 AS build

# Set Go env variables
ARG GITHUB_TOKEN
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux TOKEN=$GITHUB_TOKEN
ENV GOPRIVATE=github.com/ClareAI/*
RUN git config --global url."https://${TOKEN}@github.com".insteadOf "https://github.com"

# Set up the working directory
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENV GOCACHE=/root/.cache/go-build
ENV GOBIN=/root/bin
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o ./build/app ./.

# Stage 2: Final image
FROM golang:1.23-alpine AS runtime
WORKDIR /app

COPY --from=build /app/build .

EXPOSE 9001/tcp

ENTRYPOINT [ "./app" ]
CMD [ "-c", "./config/server.yaml"]