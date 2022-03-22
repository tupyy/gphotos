FROM quay.io/keycloak/keycloak:17.0.0 as builder

ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres

RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:17.0.0
COPY --from=builder /opt/keycloak/lib/quarkus/ /opt/keycloak/lib/quarkus/
COPY gophoto-realm-dev.json /opt/keycloak/gophoto-realm.json
WORKDIR /opt/keycloak
ENV KC_HOSTNAME=localhost:8080


