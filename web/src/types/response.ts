import { BookData, IssueRequestData, UserData } from './data'

export interface RequiredResponse {
  status: 'success' | 'error'
  message: string
}

export interface SignInResponse extends RequiredResponse {
  access_token?: string
  user?: UserData
}

export interface UserDetailsResponse extends RequiredResponse {
  user?: UserData
}

export interface RefreshAccessTokenResponse extends RequiredResponse {
  access_token?: string
}

export interface SearchBooksResponse extends RequiredResponse {
  books?: BookData[]
}

export interface SearchBookByISBNResponse extends RequiredResponse {
  book?: BookData
}

export interface IssueRequestResponse extends RequiredResponse {
  requests?: IssueRequestData[]
}

export interface GetBooksResponse extends RequiredResponse {
  books?: BookData[]
}

export interface CheckAvailabilityResponse extends RequiredResponse {
  date?: string
}
