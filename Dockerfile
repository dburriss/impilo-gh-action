# # Container image that runs your code
# FROM alpine:3.17

# # Install required dependencies
# # npm
# RUN apk add --update nodejs npm
# RUN npm -g install licensecheck

# # dotnet
# RUN apk add bash icu-libs krb5-libs libgcc libintl libssl1.1 libstdc++ zlib
# RUN apk add dotnet7-sdk

# # Copies your code file from your action repository to the filesystem path `/` of the container
# COPY entrypoint.sh /entrypoint.sh

# # Code file to execute when the docker container starts up (`entrypoint.sh`)
# ENTRYPOINT ["/entrypoint.sh"]

# Specify the version of Go to use
FROM golang:1.13

# Install required dependencies
# npm
RUN apk add --update nodejs npm
# RUN npm -g install licensecheck

# Copy all the files from the host into the container
WORKDIR /src
COPY . .

# Enable Go modules
ENV GO111MODULE=on

# Compile the action
RUN go build -o /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]