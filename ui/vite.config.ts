import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import fs from 'fs';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    strictPort: true,
    host: true,
    https: tryLoadCerts() || undefined,
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    css: true,
    restoreMocks: true,
    clearMocks: true,
    mockReset: true,
  },
});

function tryLoadCerts() {
  // Try local docker cert directory first
  const localKeyPath = '../docker/certs/server.key';
  const localCertPath = '../docker/certs/server.crt';
  
  try {
    if (fs.existsSync(localKeyPath) && fs.existsSync(localCertPath)) {
      return {
        key: fs.readFileSync(localKeyPath),
        cert: fs.readFileSync(localCertPath),
      };
    }
  } catch (e) {
    // Fall through
  }
  
  // Try absolute paths (Docker)
  try {
    if (fs.existsSync('/certs/server.key') && fs.existsSync('/certs/server.crt')) {
      return {
        key: fs.readFileSync('/certs/server.key'),
        cert: fs.readFileSync('/certs/server.crt'),
      };
    }
  } catch (e) {
    // Fall through
  }
  
  return null;
}
