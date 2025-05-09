FROM node:20-alpine as builder

WORKDIR /app


COPY package.json ./

# Install dependencies without creating a package-lock.json
# Force esbuild reinstall to ensure binary matches
RUN npm install --no-package-lock && \
    npm uninstall esbuild && \
    npm install --no-package-lock esbuild@latest

COPY . .

ARG GOOGLE_MAPS_API_KEY
ENV VITE_GOOGLE_MAPS_API_KEY=${GOOGLE_MAPS_API_KEY}

# Build without using platform-specific optimizations
RUN npm run build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]