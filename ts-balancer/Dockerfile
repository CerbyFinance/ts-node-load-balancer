FROM node:16-alpine3.12

WORKDIR /app
COPY . .

RUN npm install
RUN npm run build
RUN npm run tsc 

CMD ["node", "main.js"]