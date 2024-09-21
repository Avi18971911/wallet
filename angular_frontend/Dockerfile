FROM node:20.14 AS build
WORKDIR /app

COPY ./package.json .
COPY ./package-lock.json .
COPY ./angular.json .
COPY ./tsconfig.json .
RUN npm install

ARG PROJECT_NAME

COPY projects/${PROJECT_NAME} ./projects/${PROJECT_NAME}

# Build the Angular application
RUN npm run ng build ${PROJECT_NAME}

# Stage 2: Serve the application using Nginx
FROM nginx:alpine
ARG PROJECT_NAME
ARG NGINX_CONF
COPY --from=build /app/dist/${PROJECT_NAME}/browser /usr/share/nginx/html
COPY ./nginx/${NGINX_CONF} /etc/nginx/nginx.conf
