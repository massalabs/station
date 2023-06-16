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

import { ENV } from './const/env/env';
import './index.css';
import mockServer from './mirage/server.js';
import Error from './pages/Error.tsx';
import NotFound from './pages/NotFound.tsx';
import Base from './pages/Base/Base.tsx';
import { Store } from './pages/store/Store.tsx';

const baseURL = import.meta.env.VITE_BASE_APP;

// Add ENV.STANDALONE to the array to enable MirageJS
if ([ENV.DEV, ENV.TEST].includes(baseURL)) {
  mockServer(baseURL);
}

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={baseURL} element={<Base />}>
      {/* routes for pages */}
      <Route path="index" element={null} />
      <Route path="store" element={<Store />} />

      {/* routes for errors */}
      <Route path="error" element={<Error />} />
      <Route path="*" element={<NotFound />} />
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
