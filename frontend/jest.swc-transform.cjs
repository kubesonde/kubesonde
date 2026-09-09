/**
 * Wraps @swc/jest to make `import.meta.env` usable under Jest's CommonJS output.
 *
 * Vite exposes config via `import.meta.env`, which swc cannot emit into
 * CommonJS. We rewrite occurrences to a global that setupTests.ts populates,
 * then hand the source to the normal swc transformer.
 */
const swcJest = require("@swc/jest");

const base = swcJest.createTransformer({
    jsc: {
        parser: { syntax: "typescript", tsx: true },
        transform: { react: { runtime: "automatic" } },
        target: "es2020",
    },
});

module.exports = {
    process(src, filename, options) {
        const patched = src.replace(/import\.meta\.env/g, "globalThis.__VITE_ENV__");
        return base.process(patched, filename, options);
    },
    getCacheKey(src, filename, options) {
        return base.getCacheKey(
            src.replace(/import\.meta\.env/g, "globalThis.__VITE_ENV__"),
            filename,
            options
        );
    },
};
