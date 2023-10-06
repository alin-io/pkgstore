import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { LibConfig } from './src/components';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: LibConfig.urlPrefix,
  build: {
    outDir: '../cmd/server/ui',
    emptyOutDir: true,
  },
  server: {
    port: 3000,
  },
});
