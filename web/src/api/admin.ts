import { api } from './config'
import type { AddBookData } from '../lib/schema'
import {
  IssueRequestResponse,
  RequiredResponse,
  SearchBookByISBNResponse,
} from '../types/response'
import { UpdateBookData } from '../types/data'
import {
  ApproveRequest,
  RejectRequest,
  RemoveBookRequest,
} from '../types/request'

export async function addBook(data: AddBookData): Promise<RequiredResponse> {
  return await api
    .post(`protected/admin/add-book`, {
      json: data,
    })
    .json<RequiredResponse>()
}

export const decreaseBookCount = async (
  data: RemoveBookRequest,
): Promise<RequiredResponse> => {
  return api
    .post('protected/admin/remove-book', {
      json: data,
    })
    .json<RequiredResponse>()
}

export const getBook = async (
  isbn: string,
): Promise<SearchBookByISBNResponse> => {
  return api.get(`protected/book/${isbn}`, {}).json<SearchBookByISBNResponse>()
}

export const updateBook = async (
  data: UpdateBookData,
): Promise<RequiredResponse> => {
  return api
    .patch('protected/admin/update-book', {
      json: data,
    })
    .json<RequiredResponse>()
}

export const getIssueRequests = async (): Promise<IssueRequestResponse> => {
  return api.get('protected/admin/issue-requests').json<IssueRequestResponse>()
}

export const approveRequest = async (
  data: ApproveRequest,
): Promise<RequiredResponse> => {
  return await api
    .post('protected/admin/approve-issue-request', {
      json: data,
    })
    .json<RequiredResponse>()
}

export const rejectRequest = async (
  data: RejectRequest,
): Promise<RequiredResponse> => {
  return api
    .post('protected/admin/reject-issue-request', {
      json: data,
    })
    .json<RequiredResponse>()
}
