import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import svgr from 'vite-plugin-svgr';
import * as path from 'path';

export default ({ mode }) => {
  // loadEnv(mode, process.cwd()) will load the .env files depending on the mode
  // import.meta.env.VITE_BASE_APP available here with: process.env.VITE_BASE_APP
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return defineConfig({
    plugins: [react(), svgr()],
    base: process.env.VITE_BASE_APP,
    build: {
      outDir: '../../int/api/dist/massastation',
      emptyOutDir: true,
      manifest: true,
      sourcemap: true,
      assetsDir: './', // put the assets next to the index.html file
    },
    resolve: {
      alias: [
        { find: '@', replacement: path.resolve(__dirname, 'src') },
        { find: '@/shared', replacement: path.resolve(__dirname, '../shared') },
        // // Local UI Kit development
        // { 
        //   find: '@massalabs/react-ui-kit', 
        //   replacement: path.resolve(__dirname, '../../../ui-kit') 
        // },
      ],
    },
    server: {
      fs: {
        // to allow server ui kit asset like font files
        allow: ['../../..'],
      },
    },
  });
};
