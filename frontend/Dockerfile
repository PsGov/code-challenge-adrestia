# Use node image to build the frontend
FROM node:16 as build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Use nginx to serve the frontend
FROM nginx:alpine

# Copy built assets from the previous stage
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the port for the frontend
EXPOSE 8080

# Start NGINX server
CMD ["nginx", "-g", "daemon off;"]
