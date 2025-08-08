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
import { IconContext } from 'react-icons/lib';

import { ENV } from './const/env/env';
import {
  Error,
  NotFound,
  Index,
  Base,
  Store,
  SearchComingSoon,
  NetworkConfig,
} from './pages/index';
import { mockServer } from './mirage';
import { PAGES } from './const/pages/pages';

const baseURL = import.meta.env.VITE_BASE_APP;
const baseENV = import.meta.env.VITE_ENV;

if ([ENV.DEV, ENV.TEST].includes(baseENV)) {
  mockServer(baseENV);
}

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={baseURL} element={<Base />}>
      <Route path="" element={<Navigate to={PAGES.INDEX} />} />
      <Route path={PAGES.INDEX} element={<Index />} />

      {/* go swagger will redirect to index.html to serve index.html file */}
      {/* here we match "index.html" to navigate to "" which is Index */}
      <Route path={`${PAGES.INDEX}.html`} element={<Navigate to={baseURL} />} />

      {/* routes for pages */}
      <Route path={PAGES.STORE} element={<Store />} />
      <Route path={PAGES.SEARCH} element={<SearchComingSoon />} />
      <Route path={PAGES.CONFIG} element={<NetworkConfig />} />

      {/* routes for errors */}
      <Route path="error" element={<Error />} />
      <Route path="*" element={<NotFound />} />
    </Route>,
  ),
);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <IconContext.Provider value={{ size: '32' }}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} fallbackElement={<Error />} />
      </QueryClientProvider>
    </IconContext.Provider>
  </React.StrictMode>,
);
