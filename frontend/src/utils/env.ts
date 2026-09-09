// Thin indirection around Vite's `import.meta.env` so that other modules
// (and their tests) do not have to touch `import.meta` directly. Under Jest
// `import.meta` is not available, so tests mock this module instead.
export const getRawApiServer = (): string | undefined =>
    import.meta.env.VITE_API_SERVER;
