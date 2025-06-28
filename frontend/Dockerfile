# Use nginx as base image to serve static files
FROM nginx:alpine

# Copy the HTML file to nginx html directory
COPY index.html /usr/share/nginx/html/index.html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
