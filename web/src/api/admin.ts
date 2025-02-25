import { api } from './config'
import type { AddBookData } from '../lib/schema'

export type ApiResponse = {
  status: 'success' | 'error'
  payload: string
}

export interface Book {
  isbn: string
  title: string
  authors: string
  publisher: string
  version: string
  total_copies: number
  available_copies: number
}

export interface SearchBooksResponse {
  status: string
  payload: Book[]
}

interface UpdateBookData {
  isbn: string
  title?: string
  authors?: string
  publisher?: string
  version?: string
}

export interface IssueRequest {
  request_id: string
  isbn: string
  reader_id: string
  request_date: string
  book_title: string
  available_copies: number
}

interface IssueRequestResponse {
  status: string
  payload: IssueRequest[]
}

export const addBook = async (data: AddBookData): Promise<ApiResponse> => {
  return api
    .post(`protected/admin/add-book`, {
      json: {
        email: data.email,
        isbn: data.isbn,
        title: data.title,
        authors: data.authors,
        publisher: data.publisher,
        version: data.version,
      },
      throwHttpErrors: false,
    })
    .json<ApiResponse>()
}

export const searchBooks = async (
  searchString: string,
  field: 'title' | 'authors' | 'publisher',
): Promise<SearchBooksResponse> => {
  return api
    .post('protected/admin/books', {
      json: {
        search_string: searchString,
        field: field,
      },
    })
    .json<SearchBooksResponse>()
}

export const decreaseBookCount = async (isbn: string): Promise<ApiResponse> => {
  return api
    .post('protected/admin/remove-book', {
      json: { isbn },
    })
    .json<ApiResponse>()
}

export const getBook = async (
  isbn: string,
): Promise<{ status: string; payload: Book }> => {
  return api
    .get(`protected/admin/books/${isbn}`, {})
    .json<{ status: string; payload: Book }>()
}

export const updateBook = async (
  data: UpdateBookData,
): Promise<ApiResponse> => {
  return api
    .patch('protected/admin/update-book', {
      json: data,
    })
    .json<ApiResponse>()
}

export const getIssueRequests = async (): Promise<IssueRequestResponse> => {
  return api.get('protected/admin/issue-requests').json<IssueRequestResponse>()
}

export const approveRequest = async (
  requestId: string,
  userId: string,
): Promise<ApiResponse> => {
  return api
    .post('protected/admin/approve-issue-request', {
      json: { request_id: requestId, user_id: userId },
    })
    .json<ApiResponse>()
}

export const rejectRequest = async (
  requestId: string,
  userId: string,
): Promise<ApiResponse> => {
  return api
    .post('protected/admin/reject-issue-request', {
      json: { request_id: requestId, user_id: userId },
    })
    .json<ApiResponse>()
}

export const getLatestBooks = async (): Promise<{
  status: string
  message: string
  payload?: Book[]
}> => {
  return api
    .get('protected/admin/list-books')
    .json<{ status: string; message: string; payload?: Book[] }>()
}
