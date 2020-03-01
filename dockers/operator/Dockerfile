ARG OPERATOR_SDK_VERSION="v0.15.2"
ARG VERSION="v1alpha1"
ARG GROUP="vald.vdaas.org"
ARG KIND="Vald"

FROM quay.io/operator-framework/helm-operator:${OPERATOR_SDK_VERSION}
ARG OPERATOR_SDK_VERSION
ARG VERSION
ARG GROUP
ARG KIND

RUN { \
        echo "---"; \
        echo "- version: ${VERSION}"; \
        echo "  group: ${GROUP}"; \
        echo "  kind: ${KIND}"; \
        echo "  chart: helm-charts/vald"; \
    } > ${HOME}/watches.yaml
COPY vald ${HOME}/helm-charts/vald
