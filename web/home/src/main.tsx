import React from 'react';
import ReactDOM from 'react-dom/client';
import GlobalStyles from './styles/GlobalStyles';
import App from './App';
import './index.css';
import { QueryClient, QueryClientProvider } from 'react-query';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <QueryClientProvider client={new QueryClient()}>
    <GlobalStyles />
    <App />
  </QueryClientProvider>,
);
