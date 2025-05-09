FROM node:18-alpine as builder

WORKDIR /app

COPY package.json package-lock.json
RUN npm ci

COPY . .

ARG GOOGLE_MAPS_API_KEY
ENV VITE_GOOGLE_MAPS_API_KEY="YOUR_GOOGLE_MAPS_API_KEY"

RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]