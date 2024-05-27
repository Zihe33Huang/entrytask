# Use an official Node.js 14 image as a base
FROM node:14

# Set the working directory to /app
WORKDIR /app

# Copy the package.json file
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the frontend code
COPY frontend/. .

# Expose port 3000 for the React frontend
EXPOSE 3000

# Run the command to start the React frontend
CMD ["npm", "start"]

# Create a new directory for the Go backend
WORKDIR /backend

# Copy the Go code
COPY backend/tcp /backend/tcp
COPY backend/http /backend/http

# Install Go dependencies
RUN go mod download

# Compile the Go code
RUN go build -o tcp-server main.go
RUN go build -o http-server main.go

# Expose ports 8888 and 8080 for the Go backend
EXPOSE 8888
EXPOSE 8080

# Run the command to start the Go backend
CMD ["tcp-server", "8888"]
CMD ["http-server", "8888"]