import { defineConfig } from 'vite';
import reactRefresh from '@vitejs/plugin-react-refresh';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [reactRefresh()],
  optimizeDeps: {
    include: ['react', 'react-dom', 'react-router-dom'] // Убедитесь, что react-router-dom указан здесь
  },
  server: {
    host: '0.0.0.0'
  }
});
