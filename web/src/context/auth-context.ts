import { createContext } from 'react'
import { UserData } from '../types/data'

export interface AuthContextType {
  user: UserData | null
  login: (user: UserData) => void
  logout: () => void
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)
