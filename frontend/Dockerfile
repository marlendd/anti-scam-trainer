# === Stage 1: Build ===
FROM node:22-alpine AS build
WORKDIR /app

# Копируем файлы зависимостей и устанавливаем их
COPY package*.json ./
RUN npm ci

# Копируем исходный код и собираем dist
COPY . .
RUN npm run build


# === Stage 2: Runtime ===
FROM nginx:1.27-alpine AS runtime

ENV BACKEND_HOST=backend
ENV BACKEND_PORT=8080

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]