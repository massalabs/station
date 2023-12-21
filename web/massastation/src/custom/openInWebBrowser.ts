import { BrowserOpenURL } from '../../wailsjs/runtime';

// Opens the URL in the default web browser
// If the URL does not contain the protocol and host parts, we prepend VITE_BASE_API
export const openInWebBrowser = (url: string) => {
  if (!url.includes('://')) {
    if (url.startsWith('/')) {
      url = `${import.meta.env.VITE_BASE_API}${url}`;
    } else {
      url = `${import.meta.env.VITE_BASE_API}/${url}`;
    }
  }

  // eslint-disable-next-line new-cap
  BrowserOpenURL(url);
};
