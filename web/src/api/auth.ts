import { api } from './config'
import { ACCESS_TOKEN } from '../lib/constants'
import type { SignInRequest } from '../types/request'
import type { RegisterApiResponse, SignInResponse } from '../types/response'

// const baseApiUrl = import.meta.env.VITE_BASE_API_URL

export async function signIn(data: SignInRequest): Promise<SignInResponse> {
  const apiResponse = await api
    .post(`auth/login`, {
      json: { email: data.email },
      throwHttpErrors: false,
    })
    .json<SignInResponse>()

  if (
    apiResponse.status === 'success' &&
    apiResponse.payload?.access_token &&
    apiResponse.payload?.user
  ) {
    localStorage.setItem(ACCESS_TOKEN, apiResponse.payload.access_token) // Use ACCESS_TOKEN constant
    return {
      status: 'success',
      payload: {
        message: apiResponse.payload.message,
        user: apiResponse.payload.user,
      },
    }
  }

  return {
    status: 'error',
    payload: {
      message: apiResponse.payload?.message || "Couldn't sign in",
    },
  }
}

export function signUp(
  name: string,
  email: string,
  contact: string,
  library_id: string,
): Promise<RegisterApiResponse> {
  return api
    .post('auth/api/register', {
      json: {
        name: name,
        email: email,
        contact: contact,
        library_id: library_id,
      },
    })
    .json<RegisterApiResponse>()
}

interface RefreshResponse {
  status: 'success' | 'error'
  payload?: {
    access_token: string
    message: string
  }
}

export async function refreshToken() {
  try {
    const response = await api
      .post('auth/refresh', {
        throwHttpErrors: false,
      })
      .json<RefreshResponse>()

    if (response.status === 'success' && response.payload?.access_token) {
      return response
    }

    return null
  } catch (error) {
    console.error('Refresh token error:', error)
    return null
  }
}

interface RegisterReaderData {
  name: string
  email: string
  contact: string
  library_id: string
}

export const registerReader = async (
  data: RegisterReaderData,
): Promise<RegisterApiResponse> => {
  return api
    .post('auth/register', {
      json: data,
    })
    .json<RegisterApiResponse>()
}
