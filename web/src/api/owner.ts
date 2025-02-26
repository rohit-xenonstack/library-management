import { api } from './config'
import type { OnboardAdminData } from '../types/data'
import { RequiredResponse } from '../types/response'
import { CreateLibraryWithOwnerData } from '../types/data'

export const createLibrary = async (
  data: CreateLibraryWithOwnerData,
): Promise<RequiredResponse> => {
  return api
    .post(`create-library`, {
      json: data,
    })
    .json<RequiredResponse>()
}

export const onboardAdmin = async (
  data: OnboardAdminData,
): Promise<RequiredResponse> => {
  return api
    .post(`protected/owner/onboard-admin`, {
      json: data,
    })
    .json<RequiredResponse>()
}
