# Step 1: Use an official Node.js runtime as a base image
FROM node:16-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Step 4: Copy the rest of the application files
COPY . .

# Step 5: Expose the port the app runs on
EXPOSE 3000

# Step 6: Run the application
CMD ["npm", "start"]
