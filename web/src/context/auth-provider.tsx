import { useCallback, useEffect, useState } from 'react'
import type { ReactNode } from 'react'

import { AuthContext } from '../context/auth-context'
import { ACCESS_TOKEN } from '../lib/constants'
import type { UserDataResponse, UserDetailsResponse } from '../types/response'
import { api } from '../api/config'

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

      if (!user) {
        try {
          const response = await api
            .get('protected/me')
            .json<UserDetailsResponse>()
          if (response.status === 'success' && response.payload) {
            login(response.payload)
          } else {
            logout
          }
        } catch (err) {
          console.error(err)
          logout
        }
      }
    }

    checkAuth()

    // Check token every 10 minutes instead of every minute
    const checkInterval = setInterval(checkAuth, 10 * 60 * 1000)

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
