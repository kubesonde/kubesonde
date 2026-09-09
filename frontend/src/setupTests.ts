// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
/// <reference types="node" />
import '@testing-library/jest-dom';
import { TextEncoder, TextDecoder } from 'util';

// jsdom lacks these; react-router needs them.
const g = globalThis as unknown as Record<string, unknown>;
g.TextEncoder ??= TextEncoder;
g.TextDecoder ??= TextDecoder;

// Vite's `import.meta.env` is rewritten to this global by jest.swc-transform.cjs.
// Individual tests may still mock ./config or ./env to override values.
g.__VITE_ENV__ ??= { VITE_API_SERVER: undefined, VITE_APP_VERSION: 'test' };
