package deps

//go:generate curl -o preact.js https://unpkg.com/preact@10.5.14/dist/preact.module.js?module
//go:generate curl -o preact.d.ts https://unpkg.com/preact@10.5.14/src/index.d.ts
//go:generate curl -o jsx.d.ts https://unpkg.com/preact@10.5.14/src/jsx.d.ts
//go:generate curl -o preact-license.txt https://unpkg.com/preact@10.5.14/LICENSE
