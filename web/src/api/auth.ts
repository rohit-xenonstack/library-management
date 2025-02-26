import { api } from './config'
import { ACCESS_TOKEN } from '../lib/constants'
import type { RegisterReaderRequest, SignInRequest } from '../types/request'
import type { RequiredResponse, SignInResponse } from '../types/response'

export async function signIn(data: SignInRequest): Promise<SignInResponse> {
  const response = await api
    .post(`auth/login`, {
      json: data,
    })
    .json<SignInResponse>()

  if (response.access_token) {
    localStorage.setItem(ACCESS_TOKEN, response.access_token)
  }
  return response
}

export async function readerRegister(
  data: RegisterReaderRequest,
): Promise<RequiredResponse> {
  const response = await api
    .post('auth/register', {
      json: {
        name: data.name,
        email: data.email,
        contact: data.contact,
        library_id: data.library_id,
      },
    })
    .json<RequiredResponse>()

  return response
}
