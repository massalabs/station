import axios, { AxiosResponse } from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

/**
 * @param resource - path of the resource
 * @typeParam TBody - type of the request
 * @typeParam TResponse - type of the response
 * @returns
 */
export function usePut<TBody, TResponse = null>(
  resource: string,
): UseMutationResult<TResponse, unknown, TBody, unknown> {
  const url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<TResponse, unknown, TBody, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.put<TBody, AxiosResponse<TResponse>>(
        url,
        payload,
      );

      return data;
    },
  });
}
