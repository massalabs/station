import axios, { AxiosResponse } from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

/**
 * @param resource - path of the resource
 * @typeParam TBody - type of the request
 * @typeParam TResponse - type of the response
 * @returns
 */
export function usePost<TBody, TResponse = null>(
  resource: string,
): UseMutationResult<
  TResponse,
  unknown,
  { params?: Record<string, string>; payload?: TBody },
  unknown
> {
  const baseURL = import.meta.env.VITE_BASE_API;
  const url = `${baseURL}/${resource}`;

  return useMutation<
    TResponse,
    unknown,
    { params?: Record<string, string>; payload?: TBody },
    unknown
  >({
    mutationKey: [resource],
    mutationFn: async ({ params, payload }) => {
      const queryParams = new URLSearchParams(params).toString();
      const fullURL = `${url}?${queryParams}`;
      const decodedURL = decodeURIComponent(fullURL);
      const { data: responseData } = await axios.post<
        TBody,
        AxiosResponse<TResponse>
      >(decodedURL, payload);
      return responseData;
    },
  });
}
