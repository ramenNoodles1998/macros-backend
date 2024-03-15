# Specifies a parent image
FROM golang:1.22.1-alpine3.19
 
# Creates an src directory to hold your app’s source code
WORKDIR /src
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go mod download
 
# Builds your app with optional configuration
RUN go -C ./cmd/macros-backend build 

 
# Tells Docker which network port your container listens on
EXPOSE 3030

RUN cd ./cmd/macros-backend/
 
# Specifies the executable command that runs when the container starts
CMD [ “macros-backend” ]
