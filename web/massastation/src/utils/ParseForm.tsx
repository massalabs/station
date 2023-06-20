import { BaseSyntheticEvent } from 'react';

interface IForm {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [key: string]: any;
}

export function parseForm(e: BaseSyntheticEvent) {
  const form = new FormData(e.target as HTMLFormElement);
  const formObject: IForm = Object.fromEntries(form.entries());

  return formObject;
}
