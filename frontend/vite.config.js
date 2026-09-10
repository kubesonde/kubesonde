import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import viteTsconfigPaths from 'vite-tsconfig-paths';
import { fileURLToPath } from 'node:url';

const srcDir = fileURLToPath(new URL('./src', import.meta.url));

export default defineConfig(() => {
  return {
    build: {
      outDir: 'build',
    },
    plugins: [react(),viteTsconfigPaths()],
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ['import'],
          // Lets any component do `@use 'styles' as *;` to reach src/styles/index.scss
          loadPaths: [srcDir],
        },
      },
    },
  };
});