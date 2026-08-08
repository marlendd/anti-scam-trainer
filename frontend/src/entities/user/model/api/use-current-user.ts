import { useQuery } from '@tanstack/react-query';

import { getCurrentUser } from './get-current-user';

export function useCurrentUser() {
  return useQuery({
    queryKey: ['current-user'],
    queryFn: getCurrentUser,
  });
}