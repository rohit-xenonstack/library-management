import { api } from './config'
import type { CreateLibraryData, OnboardAdminData } from '../lib/schema'

export type ApiResponse = {
  status: 'success' | 'error'
  payload: string
}

export const createLibrary = async (
  data: CreateLibraryData,
): Promise<ApiResponse> => {
  return api
    .post(`protected/owner/create-library`, {
      json: {
        library_name: data.libraryName,
        name: data.ownerName,
        email: data.ownerEmail,
        contact: data.ownerContact,
      },
      throwHttpErrors: false,
    })
    .json<ApiResponse>()
}

export const onboardAdmin = async (
  data: OnboardAdminData,
): Promise<ApiResponse> => {
  return api
    .post(`protected/owner/onboard-admin`, {
      json: {
        name: data.adminName,
        email: data.adminEmail,
        contact: data.adminContact,
        library_id: data.libraryId,
      },
      throwHttpErrors: false,
    })
    .json<ApiResponse>()
}
