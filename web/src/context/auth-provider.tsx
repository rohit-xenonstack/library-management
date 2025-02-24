import { useCallback, useEffect, useState } from 'react'
import type { ReactNode } from 'react'

import { AuthContext } from '../context/auth-context'
import { ACCESS_TOKEN } from '../lib/constants'
import { router } from '../lib/router'
import { isTokenExpired } from '../utils/jwt'
import type { UserDataResponse } from '../types/response'

const getInitialState = () => {
  const currentUser = sessionStorage.getItem('user')
  return currentUser ? JSON.parse(currentUser) : null
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserDataResponse | null>(getInitialState)

  const login = useCallback((userData: UserDataResponse) => {
    setUser(userData)
  }, [])

  useEffect(() => {
    sessionStorage.setItem('user', JSON.stringify(user))
  }, [user])

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem(ACCESS_TOKEN)
      if (!token) {
        setUser(null)
        return
      }

      // Check if token is expired
      if (isTokenExpired(token)) {
        setUser(null)
        localStorage.removeItem(ACCESS_TOKEN)
        sessionStorage.removeItem('user')

        router.navigate({
          to: '/sign-in',
          search: { redirect: window.location.pathname },
        })
      }
    }

    checkAuth()

    // Check token every 5 minutes instead of every minute
    const checkInterval = setInterval(checkAuth, 5 * 60 * 1000)

    return () => clearInterval(checkInterval)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem(ACCESS_TOKEN)
    sessionStorage.removeItem('user')
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
