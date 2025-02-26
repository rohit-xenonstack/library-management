import { z } from 'zod'

export interface UpdateBookData {
  isbn: string
  title: string
  authors: string
  publisher: string
  version: string
}

export interface UserData {
  id: string
  name: string
  email: string
  contact: string
  role: string
  library_id: string
}

export interface BookData {
  isbn: string
  title: string
  authors: string
  publisher: string
  version: string
  total_copies: number
  available_copies: number
}

export interface IssueRequestData {
  request_id: string
  isbn: string
  reader_id: string
  request_date: string
  book_title: string
  available_copies: number
}

export const createLibraryWithOwnerSchema = z.object({
  library_name: z.string(),
  name: z.string().min(3).max(50),
  email: z.string().email(),
  contact: z.string().min(10).max(10),
})

export type CreateLibraryWithOwnerData = z.infer<
  typeof createLibraryWithOwnerSchema
>

export const onboardAdminSchema = z.object({
  name: z.string().min(3).max(100),
  email: z.string().email(),
  contact: z.string().min(10).max(10),
  library_id: z.string().uuid(),
})

export type OnboardAdminData = z.infer<typeof onboardAdminSchema>
