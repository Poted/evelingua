# Use the official Node.js image to build the Vue.js application
FROM node:18 AS build

# Set the working directory for your app
WORKDIR /app

# Copy the package.json and package-lock.json files first
COPY package.json package-lock.json ./

# Install Vue.js app dependencies
RUN npm install

# Copy the rest of the Vue.js source code into the container
COPY . .

# Build the Vue.js application
RUN npm run build

# Use a smaller image to serve the Vue.js app
FROM nginx:alpine

# Copy the build output to the Nginx server
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the app's port
EXPOSE 80

# Run the Nginx server
CMD ["nginx", "-g", "daemon off;"]
