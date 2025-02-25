import { api } from './config'
import type { Book, SearchBooksResponse } from './admin'

interface ApiResponse {
  status: string
  payload: string
}

export const searchBooks = async (
  searchString: string,
  field: 'title' | 'authors' | 'publisher',
): Promise<SearchBooksResponse> => {
  return api
    .post('protected/reader/books', {
      json: {
        search_string: searchString,
        field: field,
      },
    })
    .json<SearchBooksResponse>()
}

export const requestBook = async (
  isbn: string,
  email: string,
): Promise<ApiResponse> => {
  return api
    .post('protected/reader/request-issue', {
      json: { isbn, email },
    })
    .json<ApiResponse>()
}

export const checkAvailability = async (
  isbn: string,
): Promise<{
  status: string
  payload: string
}> => {
  return api
    .get(`protected/reader/latest/${isbn}`)
    .json<{ status: string; payload: string }>()
}

export const getBooks = async (): Promise<{
  status: string
  message: string
  payload?: Book[]
}> => {
  return api
    .get('protected/reader/get-books')
    .json<{ status: string; message: string; payload?: Book[] }>()
}
