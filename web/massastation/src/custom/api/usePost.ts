// STYLES

// EXTERNALS
import axios from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

// LOCALS

export function usePost<T, P = unknown>(
  resource: string,
): UseMutationResult<T, unknown, T, unknown> {
  var url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<T, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.post<T>(url, payload);

      return data;
    },
  });
}
