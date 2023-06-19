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
import './index.css';

import { ENV } from './const/env/env';
import { Error, NotFound, Index, Base, Store } from './pages/index';
import mockServer from './mirage/server.js';

const baseURL = import.meta.env.VITE_BASE_APP;
const baseENV = import.meta.env.VITE_ENV;

if ([ENV.DEV, ENV.TEST].includes(baseENV)) {
  mockServer(baseENV);
}

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={baseURL} element={<Base />}>
      {/* routes for pages */}
      <Route path="index" element={<Index />} />
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
