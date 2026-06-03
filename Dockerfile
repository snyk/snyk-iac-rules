FROM scratch
WORKDIR /app
# dockers_v2 organises build artifacts by platform in the build context,
# so the binary lives under the target platform directory.
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/snyk-iac-rules /usr/local/bin/
ENTRYPOINT ["snyk-iac-rules"]
