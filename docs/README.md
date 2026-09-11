# Kubesonde documentation site

The Kubesonde documentation, built with [VitePress](https://vitepress.dev/).

Content is organised with the [Diátaxis](https://diataxis.fr/) framework:

- `guide/` — introduction, quickstart, installation
- `tutorials/` — learning-oriented, start-to-finish lessons
- `how-to/` — task-oriented recipes
- `concepts/` — understanding-oriented background
- `reference/` — the `Kubesonde` resource, HTTP API, env vars, kubectl commands

## Develop

```bash
npm install
npm run docs:dev      # http://localhost:5173
```

## Build & preview

```bash
npm run docs:build
npm run docs:preview
```

Static output is written to `.vitepress/dist/`.

## Images

Static assets live in `public/`:

- `kubesonde.png` — the component/architecture diagram
- `kubesonde.pdf` — design notes / paper
- `logo.png` — the Kubesonde logo
- `screenshots/` — captures of the viewer UI
