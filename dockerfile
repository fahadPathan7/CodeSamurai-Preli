FROM golang:1.21.1

WORKDIR /app

# Copy the entire project directory into the Docker image
COPY . .

# Download and install dependencies
RUN go mod tidy
RUN go mod download

# Build the application
RUN go build -o bin .

# Expose port 5000 to the outside world
EXPOSE 5000

ENTRYPOINT [ "/app/bin" ]

# docker build . -t go-containerized:latest  # Build the image
# docker image ls | Select-String "go-containerized" # List the image
# docker run go-containerized:latest # Run the image


# docker build -t go-containerized:latest .
# docker run -p 5000:5000 go-containerized:latest # Run the image with port forwarding