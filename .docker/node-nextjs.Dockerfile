FROM node:24-alpine AS development
WORKDIR /app
COPY package*.json ./

RUN npm install
COPY . .

CMD ["npm", "run", "dev"]

FROM development AS builder
RUN npm run build

FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV HOSTNAME="0.0.0.0"

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000
CMD ["node", "server.js"]
