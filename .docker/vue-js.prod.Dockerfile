# ==========================================
# STAGE 1: Build / Kompilasi Aplikasi
# ==========================================
FROM node:24-alpine AS build-stage

WORKDIR /app

# Copy package management untuk optimasi cache layer
COPY ./frontend/package*.json ./

# Menggunakan npm ci (Clean Install) yang lebih cepat dan konsisten untuk production
RUN npm ci

# Copy source code frontend sepenuhnya
COPY ./frontend ./

# Kompilasi Vue menjadi file statis (HTML, JS, CSS) di folder /app/dist
RUN npm run build

# ==========================================
# STAGE 2: Web Server Runtime Minimalis
# ==========================================
FROM nginx:1.25-alpine AS production-stage

# Menyalin hasil build dari stage 1 ke direktori root Nginx
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY ./.docker/nginx/default.conf /etc/nginx/conf.d/default.conf
# Expose port HTTP default internal
EXPOSE 80

# Jalankan Nginx di foreground agar container tetap hidup
CMD ["nginx", "-g", "daemon off;"]
