# Step 1: Use an official Node.js runtime as a base image
FROM node:20-alpine AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy only package.json and package-lock.json for efficient caching
COPY package.json package-lock.json ./

# Step 4: Install production dependencies only
RUN npm ci --omit=dev

# Step 5: Copy the rest of the application files
COPY . .

# Step 6: Expose the application port
EXPOSE 3000

# Step 7: Set the command to start the app
CMD ["npm", "start"]
