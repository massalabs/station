// STYLES

// EXTERNALS
import axios from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

// LOCALS

export function usePost<T, P = unknown>(
  resource: string,
): UseMutationResult<P, unknown, T, unknown> {
  var url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<P, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.post<P>(url, payload);

      return data;
    },
  });
}
