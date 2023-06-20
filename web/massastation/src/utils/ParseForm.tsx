import { BaseSyntheticEvent } from 'react';

export function parseForm<T>(e: BaseSyntheticEvent) {
  const form = new FormData(e.target as HTMLFormElement);
  const formObject = Object.fromEntries(form.entries()) as T;

  return formObject;
}
