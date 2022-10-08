FROM nginx:latest


EXPOSE 80/tcp


ENV CONSUL_TEMPLATE_VERSION=0.29.5
ENV CONSUL_URL consul:8500
ENV TEMPLATE_PATH /scripts/nginx.conf.tmpl
# Update apt
RUN apt-get update \
	&& apt-get install -y -q wget curl unzip lsb-release runit
# Download consul-template
#ENV CT_URL https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip
#RUN curl -O $CT_URL
## Unzip consul-template
COPY nginx/consul-template_0.29.5_linux_amd64.zip consul-template_0.29.5_linux_amd64.zip
RUN unzip consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip -d /usr/bin \
  && rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip

#Create entrypoint
RUN mkdir -p /scripts/
COPY nginx/start.sh /scripts/start.sh
COPY nginx/nginx.conf.tmpl /scripts/nginx.conf.tmpl
COPY nginx/consul-template-config.hcl /scripts/consul-template-config.hcl


RUN ["chmod", "+x", "/scripts/start.sh"]
ENTRYPOINT ["/scripts/start.sh"]