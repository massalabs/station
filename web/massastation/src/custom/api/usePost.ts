// STYLES

// EXTERNALS
import axios from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

// LOCALS

export function usePost<T, P = unknown>(
  resource: string,
): UseMutationResult<
  P,
  unknown,
  { params?: Record<string, string>; payload?: T },
  unknown
> {
  const baseURL = import.meta.env.VITE_BASE_API;
  const url = `${baseURL}/${resource}`;

  return useMutation<
    P,
    unknown,
    { params?: Record<string, string>; payload?: T },
    unknown
  >({
    mutationKey: [resource],
    mutationFn: async ({ params, payload }) => {
      const queryParams = new URLSearchParams(params).toString();
      const fullURL = `${url}?${queryParams}`;
      const decodedURL = decodeURIComponent(fullURL);
      const { data: responseData } = await axios.post<P>(decodedURL, payload);
      return responseData;
    },
  });
}
