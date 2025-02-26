import ky, { type BeforeRequestHook } from 'ky'

import { ACCESS_TOKEN } from '../lib/constants'
import { router } from './router'
import { isTokenExpired } from '../utils/jwt'
import { RefreshAccessTokenResponse } from '../types/response'

const baseApiUrl = import.meta.env.VITE_BASE_API_URL

export const beforeRequest: BeforeRequestHook = async (request) => {
  var token = localStorage.getItem(ACCESS_TOKEN)
  if (!token) return
  if (isTokenExpired(token)) {
    const data = await ky
      .get(`${baseApiUrl}/auth/refresh`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem(ACCESS_TOKEN)}`,
        },
        retry: 0,
      })
      .json<RefreshAccessTokenResponse>()

    if (!data.access_token) {
      localStorage.removeItem(ACCESS_TOKEN)
      sessionStorage.removeItem('user')
      router.navigate({ to: '/sign-in' })
      return
    }

    token = data.access_token
    console.log('Access token refreshed:', token)
    localStorage.setItem(ACCESS_TOKEN, token)
  }
  request.headers.set('Authorization', `Bearer ${token}`)
  return ky(request)
}
