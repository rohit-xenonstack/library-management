import ky from 'ky'

import { afterResponse, beforeRequest } from '../lib/auth-interceptor'

export const baseApiUrl = import.meta.env.VITE_BASE_API_URL
console.log(baseApiUrl)

export const api = ky.create({
  prefixUrl: baseApiUrl,
  hooks: {
    beforeRequest: [beforeRequest],
    afterResponse: [afterResponse],
  },
  retry: 0,
  credentials: 'include',
})
