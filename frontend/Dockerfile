# Stage 1: Build the Vue application
FROM node:20-alpine AS build-stage

WORKDIR /app

# Copy package.json and lockfile
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build the frontend (Vite)
# Force build since it's production
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine AS production-stage

# Copy the built files from build-stage to Nginx public folder
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
# We'll create this file next
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
