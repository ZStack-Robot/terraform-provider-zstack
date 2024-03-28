ARG TERRAFORM_VERSION=1.5.7

#
# Fetch
#
FROM --platform=$TARGETPLATFORM hashicorp/terraform:${TERRAFORM_VERSION} AS fetch

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN set -eo pipefail; \
    apk add -U --no-cache \
      ca-certificates \
      curl openssh \
      unzip \
      git \
      jq \
    ; \
    rm -rf /var/cache/apk/*

ENV TF_LOG="ERROR"
WORKDIR /workspace

# download templates
RUN set -eo pipefail; \
    echo "walrus-catalog walrus-catalog-sandbox" | tr -s '[:blank:]' '\n' | \
    while read -r org _; do \
      curl -sSL "https://api.github.com/orgs/$org/repos" | jq -r '.[].name' | \
      while read -r repo _; do \
        git clone "https://github.com/$org/$repo" "$org"_"$repo" --depth 1; \
      done; \
    done

# mirror plugins
## cache plugins to reduce network latency
ENV TF_PLUGIN_CACHE_DIR="/workspace/.terraform.d/plugin-cache" \
    TF_PLUGIN_MIRROR_DIR="/workspace/.terraform.d/plugins"
RUN set -eo pipefail; \
    mkdir -p $TF_PLUGIN_CACHE_DIR; \
    mkdir -p $TF_PLUGIN_MIRROR_DIR; \
    echo -e "provider_installation {\n \
    dev_overeides {\n \
      \"zstack.io/terraform-provider-zstack/zstack\" = \"/usr/local/bin/\" \n \
      \"hashicorp/local\" = \"/usr/local/bin/\" \n \
       \"hashicorp/local\" = \"/usr/local/bin/\" \n \
       }\n\
      filesystem_mirror {\n \
        path = \"$TF_PLUGIN_MIRROR_DIR\"\n \
      }\n \
      direct {} \n \
    }\n" > /root/.terraformrc && \
    find . -maxdepth 1 -type d -name 'walrus-catalog*' -exec sh -c 'terraform -chdir="$1" init && terraform -chdir="$1" providers mirror $TF_PLUGIN_MIRROR_DIR' _ {} \;
## remove non-plugin files to prevent annoying message
RUN set -eo pipefail; \
    find $TF_PLUGIN_MIRROR_DIR -type f ! -name "terraform-provider-*" -delete

#
# Release
#
FROM --platform=$TARGETPLATFORM hashicorp/terraform:${TERRAFORM_VERSION}
LABEL maintainer="Seal Engineer Team <engineering@seal.io>"

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN set -eo pipefail; \
    apk add -U --no-cache \
      ca-certificates \
      openssl \
      curl unzip \
    ; \
    rm -rf /var/cache/apk/*

# set locale
RUN set -eo pipefail; \
    apk add -U --no-cache \
      tzdata \
    ; \
    rm -rf /var/cache/apk/*
ENV LANG='en_US.UTF-8' \
    LANGUAGE='en_US:en' \
    LC_ALL='en_US.UTF-8'

# get kubectl for gavinbunney/kubectl provider
RUN KUBECTL_VER="v1.25.5"; \
    curl -sfL https://dl.k8s.io/${KUBECTL_VER}/kubernetes-client-${TARGETOS}-${TARGETARCH}.tar.gz | \
        tar -xvzf - --strip-components=3 --no-same-owner -C /usr/bin/ kubernetes/client/bin/kubectl

# get terraform plugins
COPY --from=fetch /workspace/.terraform.d/plugins /usr/share/terraform/providers/plugins

ENV TF_LOG=INFO \
    TF_IN_AUTOMATION=1 \
    TF_PLUGIN_MIRROR_DIR="/usr/share/terraform/providers/plugins"
COPY terraform /usr/local/bin/terraform
COPY terraform-provider-zstack /usr/local/bin/terraform-provider-zstack
COPY terraform-provider-template /usr/local/bin/terraform-provider-template
COPY terraform-provider-local /usr/local/bin/terraform-provider-local
ENTRYPOINT ["terraform"]