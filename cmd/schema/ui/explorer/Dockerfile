FROM node:12
WORKDIR /build

# install yarn
RUN apt-get update && apt-get install -y curl
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get install -y yarn

# install node deps
RUN yarn global add @vue/cli @vue/cli-service @vue/cli-plugin-babel @vue/cli-plugin-unit-jest

RUN chown -R 1000:1000 /build
USER 1000