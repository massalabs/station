import React from 'react';

import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import '@massalabs/react-ui-kit/src/global.css';
import mockServer from './mirage/server.js';
import './index.css';
import Error from './pages/Error.tsx';
import { ENV } from './const/env/env';

const queryClient = new QueryClient();
// Add ENV.STANDALONE to the array to enable MirageJS
if ([ENV.DEV, ENV.TEST].includes(import.meta.env.VITE_ENV)) {
  mockServer(import.meta.env.VITE_ENV);
}
const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={import.meta.env.VITE_BASE_APP}>
      {/* routes for errors */}
      <Route path="error" element={<Error />} />
      <Route path="*" element={<Error />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} fallbackElement={<Error />} />
    </QueryClientProvider>
  </React.StrictMode>,
);
