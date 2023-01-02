ARG BASIC_IMAGE=ghcr.io/hankei6km/h6-dev-containers:dind
FROM ${BASIC_IMAGE} as go

ENV H6_IMAGE_TYPE=go

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

USER root
RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list \
    && sudo apt update \
    && sudo apt install -y --no-install-recommends goreleaser \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

USER vscode
RUN cd /tmp \
    && mkdir mage && cd mage \
    && curl -s https://api.github.com/repos/magefile/mage/releases/latest | jq .assets[].browser_download_url | grep _Linux-64bit.tar.gz | xargs -I '{}' curl -sL '{}' | tar -zxf - \
    && cp mage "/home/vscode/.local/bin/mage" \
    && cd .. && rm -rf mage