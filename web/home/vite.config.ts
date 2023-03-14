import { build, defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '',
  build: {
    outDir: "../../int/api/dist/home",
    emptyOutDir: true,
    assetsDir: './'
  }
})
