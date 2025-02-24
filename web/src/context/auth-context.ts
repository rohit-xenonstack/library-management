import { createContext } from 'react'

import type { UserDataResponse } from '../types/response'

export interface AuthContextType {
  user: UserDataResponse | null
  login: (user: UserDataResponse) => void
  logout: () => void
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)
