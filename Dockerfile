# Start from the latest golang base image
FROM golang:1.14-alpine as base

# Add Maintainer Info
LABEL maintainer="The Basebandit <@the_basebandit>"

# Set the Current Working Directory inside the container
WORKDIR /api

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the api with "-ldflags" aka linker flags to reduce binary size
# -s = disable symbol table
# -w = disable DWARF generation
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main ./cmd/

# # Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fupisha ./cmd/

FROM aquasec/trivy:0.4.4 as trivy

# RUN executes commands on top of the current image as a new layer and commits the results
# Scan the golang alpine image before production use
RUN trivy --debug --timeout 4m golang:1.14-alpine && \
  echo "No image vulnerabilities" > result


######## Start a new stage from scratch #######
FROM base as prod 

ENV FUPISHA_USER fupisha
ENV FUPISHA_GROUP api
ENV FUPISHA_HOME /go/src/fupisha
# SMTP config
ENV FUPISHA_SMTP_PORT=587
ENV FUPISHA_SMTP_HOST=smtp.gmail.com
ENV FUPISHA_SMTP_USER=smtp_username
ENV FUPISHA_SMTP_PASSWORD=smtp_password
ENV FUPISHA_SMTP_FROM_NAME=Fupisha
ENV FUPISHA_SMTP_FROM_ADDRESS=no-reply@fupisha.io

#Store type
ENV FUPISHA_STORE_TYPE=mongo

#Mongo config
# mongodb is the name of the mongodb container in docker-compose file
ENV FUPISHA_STORE_MONGO_ADDRESS=mongodb:27017
ENV FUPISHA_STORE_MONGO_USERNAME=fupisha
ENV FUPISHA_STORE_MONGO_PASSWORD=fupisha
ENV FUPISHA_STORE_MONGO_DATABASE=fupisha

#Auth config
ENV FUPISHA_JWT_SECRET=b0f635be4e3fc72030c33f3d8011be5b8e966930108d8ee85ed9f0a43647a0ae
ENV FUPISHA_JWT_EXPIRE_DELTA=6

#Fupisha config
ENV FUPISHA_BASE_URL=https://fupisha.io
ENV FUPISHA_TITLE=Fupisha
ENV FUPISHA_LOG_LEVEL=info
ENV FUPISHA_TEXT_LOGGING=false



RUN apk update && apk --no-cache add ca-certificates
RUN apk --no-cache add  bash


COPY --from=trivy result secure
#Copy the email templates from the previous stage
COPY --from=base /api/templates .

# Copy the Pre-built binary file from the previous stage
COPY --from=base /api/main .

#Copy the wait script file from the previous stage
COPY --from=base /api/wait-for-it.sh .

# # Create a group and user
# RUN addgroup $FUPISHA_GROUP && adduser -D -G $FUPISHA_USER $FUPISHA_GROUP

# Create a new group and user, recursively change directory ownership, then give permission to run script
RUN addgroup fupisha && adduser -D -G fupisha fupisha \
  && chown -R fupisha:fupisha /api && \
  chmod +x ./wait-for-it.sh && \
  chmod +x ./main

# RUN chown -R fupisha:fupisha ./wait-for-it.sh && \
#   chown -R chown -R $FUPISHA_USER:$FUPISHA_USER ./fupisha && \
#   chmod +x ./wait-for-it.sh && \
#   chmod +x ./fupisha

# Tell docker that all future commands should run as the your user
USER fupisha

# Expose ports to the outside world
# port 8080 - api server
EXPOSE 8080

# Command to run the executable
CMD ["./fupisha","start"] 