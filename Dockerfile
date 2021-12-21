FROM scratch
WORKDIR /app
# Copy from the tmp directory created as part of the Docker release
# this is equivalent to the /dist dir created by goreleaser.
COPY ./snyk-iac-rules /usr/local/bin/
ENTRYPOINT ["snyk-iac-rules"]
