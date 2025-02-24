interface JWTPayload {
  exp: number
  sub: string
  role: string
}

export const isTokenExpired = (token: string): boolean => {
  try {
    const [, base64Payload] = token.split('.')
    const payload = JSON.parse(atob(base64Payload)) as JWTPayload
    return payload.exp * 1000 < Date.now()
  } catch {
    return true
  }
}
