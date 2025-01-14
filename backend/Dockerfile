# Use the official Go image to build the Go application
FROM golang:1.23.4 AS build

# Set the working directory for your app
WORKDIR /app

# Copy the Go app's source code into the container
COPY . .

# Install dependencies and build the Go app
# RUN apt-get update && apt-get install -y libc6
RUN go mod tidy
RUN go build -o main .

# Use the same Go image to run the Go application
FROM golang:1.23.4

# Set the working directory for the final image
WORKDIR /root/

# Copy the Go binary from the build stage
COPY --from=build /app/main .

# Expose the app's port
EXPOSE 4000

# Run the Go app
CMD ["./main"]
