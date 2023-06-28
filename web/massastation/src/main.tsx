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
import { Network, Base } from './components';
import { mockServer } from './mirage';
import { Error, NotFound, Index, Store, Search } from './pages/index';
import { PAGES } from './const/pages/pages';

const baseURL = import.meta.env.VITE_BASE_APP;
const baseENV = import.meta.env.VITE_ENV;

if ([ENV.DEV, ENV.TEST].includes(baseENV)) {
  mockServer(baseENV);
}

const queryClient = new QueryClient();

const BUILDNET = 'buildnet';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Base />}>
      <Route path="" element={<Navigate to={BUILDNET} />} />

      <Route path=":network" element={<Network />}>
        <Route path="" element={<Navigate to={PAGES.INDEX} />} />

        <Route path={baseURL}>
          {/* routes for pages */}
          <Route path={PAGES.INDEX} element={<Index />} />
          <Route path={PAGES.STORE} element={<Store />} />
          <Route path={PAGES.SEARCH} element={<Search />} />
        </Route>
      </Route>

      {/* go swagger will redirect to index.html to serve index.html file */}
      {/* here we match "index.html" to navigate to "" which is Index */}
      <Route
        path={`${PAGES.INDEX}.html`}
        element={<Navigate to={BUILDNET} />}
      />

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
