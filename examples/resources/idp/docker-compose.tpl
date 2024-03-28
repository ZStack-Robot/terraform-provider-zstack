version: '3.8'
services:
  keycloak:
    image: keycloak-zsthemev2.1:21.1.1
    network_mode: "host"
    volumes:
      - ./cache-custom-jgroups-tcp.xml:/opt/keycloak/conf/cache-custom-jgroups-tcp.xml
    command: start --health-enabled=true --metrics-enabled=true --cache=ispn --hostname-strict=false  --spi-theme-welcome-theme=zstack-theme
    environment:
      KC_CACHE: ispn
      KC_DB: ${db_name}
      KC_DB_URL_HOST: ${db_url}
      KC_DB_USERNAME: ${db_user}
      KC_DB_PASSWORD: ${db_pwd}
      KC_DB_URL_DATABASE: ${db_name}
      KC_PROXY: edge
      KEYCLOAK_ADMIN: ${idp_admin}
      KEYCLOAK_ADMIN_PASSWORD: ${idp_password}
      KC_CACHE_CONFIG_FILE: cache-custom-jgroups-tcp.xml
      KC_HTTP_ENABLED: true
      DEBUG: "true"
      DEBUG_PORT: "*:8787"
    extra_hosts:
%{ for ip_index,ip_value  in cluster ~}
      - idp${ip_index}.zstack.io:${ip_value}
%{ endfor ~}

  lb:
    image: nginx:alpine
    network_mode: "host"
    volumes:
      - ./idp.conf:/etc/nginx/conf.d/default.conf
    extra_hosts:
%{ for ip_index,ip_value  in cluster ~}
      - idp${ip_index}.zstack.io:${ip_value}
%{ endfor ~}