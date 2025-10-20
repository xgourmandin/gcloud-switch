FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY gcloud-switcher /usr/local/bin/gcloud-switcher

RUN chmod +x /usr/local/bin/gcloud-switcher

# Install gcloud CLI (optional, but recommended)
RUN apk add --no-cache python3 py3-pip curl bash && \
    curl https://sdk.cloud.google.com | bash && \
    ln -s /root/google-cloud-sdk/bin/gcloud /usr/local/bin/gcloud

ENTRYPOINT ["/usr/local/bin/gcloud-switcher"]
CMD ["--help"]

