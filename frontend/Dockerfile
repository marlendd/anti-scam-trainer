FROM nginx:1.27-alpine AS runtime

ENV BACKEND_HOST=backend
ENV BACKEND_PORT=8080

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]