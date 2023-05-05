# Implementation Notes

## Project initialization

### NodeJS with Typescript

See
- [this video](https://www.youtube.com/watch?v=H91aqUHn8sE) and;
- [this page](https://fireship.io/lessons/typescript-nodejs-setup/)

```shell
$ npm init -y
$ : insert "type": "module" into package.json for ES modules
$ npm install typescript --save-dev
$ : insert "build": "tsc" into package.json:scripts for compiling
$ mkdir src
$ touch src/server.ts
$ npm i -D @types/node # possibly optional
$ echo '
> {
>   "compilerOptions": {
>     "module": "NodeNext",
>     "moduleResolution": "NodeNext",
>     "target": "ES2020",
>     "sourceMap": true,
>     "outDir": "dist",
>   },
>   "include": ["src/**/*"],
> }' > tsconfig.json
$ npm run build
$ node dist/server.js
$ npm i -D nodemon ts-node
$ : add "dev:server": "nodemon --watch './**/*.ts' --exec 'ts-node' server.ts" to package.json
```

### Express

```shell
$ npm install express @types/express
```