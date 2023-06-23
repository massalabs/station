import axios, { AxiosResponse } from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

/**
 * @param resource - path of the resource
 * @typeParam TBody - type of the request
 * @typeParam TResponse - type of the response
 * @returns
 */
export function usePut<TBody, TResponse = null, TError = unknown>(
  resource: string,
  headers: object = {},
): UseMutationResult<TResponse, TError, TBody, unknown> {
  const url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<TResponse, TError, TBody, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.put<TBody, AxiosResponse<TResponse>>(
        url,
        payload,
        { headers },
      );

      return data;
    },
  });
}
