import React from 'react';
import ReactDOM from 'react-dom/client';
import { Theme } from '@radix-ui/themes';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App.tsx';

import '@radix-ui/themes/styles.css';

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <Theme appearance="dark" accentColor="violet" grayColor="slate" radius="large">
        <App />
      </Theme>
    </QueryClientProvider>
  </React.StrictMode>
);