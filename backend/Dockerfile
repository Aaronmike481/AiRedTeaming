FROM mcr.microsoft.com/playwright/python:v1.48.0-jammy

# Install Go
RUN apt-get update && apt-get install -y wget && \
    wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz && \
    rm go1.23.0.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app

# Copy Go server files (NO 'backend/' prefix because we are already inside backend)
COPY go-server/ ./go-server/
WORKDIR /app/go-server
RUN go mod tidy && go build -o server .

# Install Python dependencies
WORKDIR /app
COPY python-worker/ ./python-worker/
WORKDIR /app/python-worker
RUN pip install pyproc-worker requests openai playwright
RUN playwright install chromium

# Run Go server
WORKDIR /app/go-server
CMD ["./server"]