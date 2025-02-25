import ky from 'ky'

import { beforeRequest } from '../lib/auth-interceptor'

export const baseApiUrl = import.meta.env.VITE_BASE_API_URL

export const api = ky.create({
  prefixUrl: baseApiUrl,
  hooks: {
    beforeRequest: [beforeRequest],
  },
  retry: {
    limit: 0,
  },
  credentials: 'include',
  throwHttpErrors: false,
})
