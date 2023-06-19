import React from 'react';
import ReactDOM from 'react-dom/client';
import {
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  Navigate,
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
      <Route path="" element={<Index />} />

      {/* go swagger will redirect to index.html to serve index.html file */}
      {/* here we match index.html en navigate to "" which is Index */}
      <Route path={'index.html'} element={<Navigate to={baseURL} />} />

      {/* routes for pages */}
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
