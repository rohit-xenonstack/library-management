import { SearchBookRequest } from '../types/request'
import { SearchBooksResponse, GetBooksResponse } from '../types/response'
import { api } from './config'

export const searchBooks = async (
  data: SearchBookRequest,
): Promise<SearchBooksResponse> => {
  return api
    .post('protected/book', {
      json: data,
    })
    .json<SearchBooksResponse>()
}

export const getBooks = async (): Promise<GetBooksResponse> => {
  return api.get('protected/books').json<GetBooksResponse>()
}
