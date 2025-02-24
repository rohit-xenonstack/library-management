import type { AfterResponseHook, BeforeRequestHook } from 'ky'

import { refreshToken } from '../api/auth'
import { ACCESS_TOKEN } from '../lib/constants'
import { router } from '../lib/router'
import { isTokenExpired } from '../utils/jwt'

export const beforeRequest: BeforeRequestHook = async (request) => {
  const token = localStorage.getItem(ACCESS_TOKEN)
  if (!token) return

  if (isTokenExpired(token)) {
    const newToken = await refreshToken()
    if (newToken) {
      request.headers.set('Authorization', `Bearer ${newToken}`)
      return
    }
    router.navigate({
      to: '/sign-in',
      search: { redirect: window.location.pathname },
    })
    return
  }

  request.headers.set('Authorization', `Bearer ${token}`)
}

export const afterResponse: AfterResponseHook = async (
  request,
  options, // eslint-disable-line @typescript-eslint/no-unused-vars
  response,
) => {
  if (response.status === 401) {
    const newToken = await refreshToken()

    if (!newToken || !newToken.payload || newToken?.status === 'error') {
      router.navigate({
        to: '/sign-in',
        search: { redirect: window.location.pathname },
      })
      return response
    }

    localStorage.setItem(ACCESS_TOKEN, newToken.payload?.access_token)
    request.headers.set('Authorization', `Bearer ${newToken}`)
  }
  return response
}
