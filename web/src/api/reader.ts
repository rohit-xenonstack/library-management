import { RequestBookRequest } from '../types/request'
import {
  CheckAvailabilityResponse,
  GetBooksResponse,
  RequiredResponse,
} from '../types/response'
import { api } from './config'

export const requestBook = async (
  data: RequestBookRequest,
): Promise<RequiredResponse> => {
  return api
    .post('protected/reader/request-issue', {
      json: data,
      throwHttpErrors: false,
    })
    .json<RequiredResponse>()
}

export const checkAvailability = async (
  isbn: string,
): Promise<CheckAvailabilityResponse> => {
  return api
    .get(`protected/reader/latest/${isbn}`)
    .json<CheckAvailabilityResponse>()
}

export const getBooks = async (): Promise<GetBooksResponse> => {
  return api.get('protected/books').json<GetBooksResponse>()
}
