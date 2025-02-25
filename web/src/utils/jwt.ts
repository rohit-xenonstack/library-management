interface JWTPayload {
  id: string // token ID
  user_id: string // user ID
  role: string // user role
  issued_at: string // token issue time
  expires: string // token expiry time
}

export const isTokenExpired = (token: string): boolean => {
  if (!token) return true

  try {
    // Split the token into parts
    const parts = token.split('.')
    if (parts.length !== 3) return true

    // Get the payload part
    const base64Payload = parts[1]

    // Add padding if needed
    const padding = '='.repeat((4 - (base64Payload.length % 4)) % 4)
    const base64 = (base64Payload + padding)
      .replace(/-/g, '+')
      .replace(/_/g, '/')

    // Decode and parse the payload
    const payload = JSON.parse(atob(base64)) as JWTPayload

    // Check if expiry time exists and is valid
    if (!payload.expires) return true

    // Convert expiry time to milliseconds and compare with current time
    const expiryTime = new Date(payload.expires).getTime()
    const currentTime = Date.now()

    return currentTime >= expiryTime
  } catch (error) {
    console.error('Error parsing JWT:', error)
    return true
  }
}

// Helper function to get payload from token
export const getTokenPayload = (token: string): JWTPayload | null => {
  if (!token) return null

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const base64Payload = parts[1]
    const padding = '='.repeat((4 - (base64Payload.length % 4)) % 4)
    const base64 = (base64Payload + padding)
      .replace(/-/g, '+')
      .replace(/_/g, '/')

    return JSON.parse(atob(base64)) as JWTPayload
  } catch {
    return null
  }
}
