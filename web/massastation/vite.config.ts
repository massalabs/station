import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default ({ mode }) => {
  // loadEnv(mode, process.cwd()) will load the .env files depending on the mode
  // import.meta.env.VITE_BASE_APP available here with: process.env.VITE_BASE_APP
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return defineConfig({
    plugins: [react()],
    base: process.env.VITE_BASE_APP,
    build: {
      outDir: '../../int/api/dist/massastation',
      emptyOutDir: true,
      manifest: true,
      sourcemap: true,
      assetsDir: './', // put the assets next to the index.html file
    },
  });
};
