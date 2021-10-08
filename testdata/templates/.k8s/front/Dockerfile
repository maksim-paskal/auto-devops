FROM git.test.com:5005/devops/docker/nginx:1.17-production

ENV HOST="127.0.0.1"
ENV BACKEND="127.0.0.1"

ADD ./public /app/public

COPY ./.k8s/front/default.conf /etc/nginx/conf.d/default.conf.template
