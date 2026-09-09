/**
 * Jest configuration for the frontend.
 *
 * Uses @swc/jest to transform TS/TSX (typescript@7 rules out ts-jest), runs in
 * jsdom, maps CSS/SCSS imports to a proxy, and honours the `src/*` path alias
 * from tsconfig. `import.meta.env` is only touched in src/utils/env.ts, which
 * tests mock, so it is never transformed here.
 */
module.exports = {
    testEnvironment: "jsdom",
    setupFilesAfterEnv: ["<rootDir>/src/setupTests.ts"],
    moduleNameMapper: {
        "\\.(css|scss|sass)$": "identity-obj-proxy",
        "^react-cytoscapejs$": "<rootDir>/src/testStubs/reactCytoscape.tsx",
        "^src/(.*)$": "<rootDir>/src/$1",
    },
    transform: {
        "^.+\\.(t|j)sx?$": "<rootDir>/jest.swc-transform.cjs",
    },
    // react-router (and a few deps) ship ESM; let swc transform them too.
    transformIgnorePatterns: [
        "/node_modules/(?!(react-router|react-router-dom|@remix-run|react-cytoscapejs)/)",
    ],
    testMatch: ["<rootDir>/src/**/*.test.{ts,tsx}"],
    collectCoverageFrom: [
        "src/**/*.{js,jsx,ts,tsx}",
        "!<rootDir>/node_modules/",
        "!src/entities/**",
    ],
};
