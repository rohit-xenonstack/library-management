export const ACCESS_TOKEN = 'access_token'

export const DASHBOARD = '/' as const
export const LOGIN_PAGE = '/sign-in' as const

export const enum ROLE {
  ADMIN = 'admin',
  OWNER = 'owner',
  READER = 'reader',
}
