FROM google/golang

# these are useful for "go get"
RUN apt-get update && apt-get install -yq bzr git mercurial

RUN go get github.com/dustin/go-humanize

WORKDIR /gopath/src/app
ADD . /gopath/src/app/
RUN go install app

# RUNTIME CONFIG
# Set the following here in the Dockerfile 
# or use docker run -e USER_CREDS=username:password
# (the -e option is safer, just in case you 
# accidentally ever push this image to the Index)
# ENV USER_CREDS=username:password
CMD []
ENTRYPOINT ["/gopath/bin/app"]

