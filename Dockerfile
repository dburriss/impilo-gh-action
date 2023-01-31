# Specify the version of Go to use
FROM golang:1.19.5-alpine

# Install required dependencies
# npm
RUN apk update && apk add nodejs npm
RUN npm --version
# RUN npm -g install licensecheck

# Copy all the files from the host into the container
WORKDIR /src
COPY ./src/go.mod ./src/go.sum ./

# Download packages
RUN go mod download && go mod verify

# Compile the action
COPY /src .
RUN go build -v -o /usr/local/bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["action"]