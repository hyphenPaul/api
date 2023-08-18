# Use a Golang base image for M1 Mac
FROM arm64v8/golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go dependencies
RUN go mod download && go mod verify && go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application
#
# CGO_ENABLED=0
# It's worth noting that not all Go applications can be built with CGO_ENABLED=0,
# especially if they rely on certain C dependencies. However, if your application
# doesn't require any C bindings, setting CGO_ENABLED=0 is generally recommended to
# minimize the size and complexity of the resulting Docker image.
#
# GOARCH=arm64
# In summary, setting GOARCH=arm64 in a Go Dockerfile indicates that you want to build
# the Go application specifically for the ARM64 architecture, making it compatible with
# ARM-based devices or platforms.
#
# GOOD=linux
# Setting GOOS=linux is common when building Go applications for containerization because
# Docker containers typically run on Linux-based hosts. By building the application specifically
# for Linux, you ensure that it is compatible with the underlying Linux kernel and can
# run in a Linux container environment.
#
# It's worth noting that Go supports building applications for multiple operating systems and
# architectures. By setting the appropriate GOOS and GOARCH values, you can cross-compile your
# Go application for different target platforms, such as Windows, macOS, or various architectures
# like ARM or AMD64.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app

# Set the command to run when the container starts
CMD go run api
