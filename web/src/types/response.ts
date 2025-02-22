export type AccessToken = {
  accessToken: string,
}

export type RefreshResponse = {
  status: string,
  payload: AccessToken,
}