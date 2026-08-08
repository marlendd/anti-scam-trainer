import { useMutation, useQueryClient } from '@tanstack/react-query';

import { currentUserQueryKey } from '@/entities/user';

import { logout } from '../api/logout.ts';

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.setQueryData(currentUserQueryKey, null);
    },
  });
}