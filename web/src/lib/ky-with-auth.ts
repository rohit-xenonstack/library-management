import ky, { HTTPError } from 'ky'

import { ACCESS_TOKEN } from './constants'
import type { RefreshResponse } from '../types/response'

const kyWithAuth = ky.create({
  prefixUrl: import.meta.env.VITE_BASE_API_URL,
  credentials: 'include',
  hooks: {
    beforeRequest: [
      (request) => {
        const accessToken = localStorage.getItem(ACCESS_TOKEN)
        if (!accessToken) return

        request.headers.set('Authorization', `Bearer ${accessToken}`)
      },
    ],
    beforeRetry: [
      async ({ request, error, options }) => {
        if (error instanceof HTTPError && error.response.status === 401) {
          const data = await ky
            .get('token/refresh', { ...options, retry: 0 })
            .json<RefreshResponse>()

          const newAccessToken = data.payload.accessToken

          localStorage.setItem(ACCESS_TOKEN, newAccessToken)

          request.headers.set('Authorization', `Bearer ${newAccessToken}`)
        }
      },
    ],
  },
  retry: {
    limit: 1,
    statusCodes: [401, 403, 500, 504],
  },
})

export { kyWithAuth as ky }