export interface UserDataResponse {
  user_id: string
  name: string
  email: string
  contact: string
  role: string
  library_id?: string
}

interface LoginResponsePayload {
  message: string
  access_token?: string
  user?: UserDataResponse
}

export interface SignInResponse {
  status: 'success' | 'error'
  payload: LoginResponsePayload
}

export interface RegisterApiResponse {
  status: 'success' | 'error'
  payload: string
}
