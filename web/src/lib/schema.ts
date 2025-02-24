import { z } from 'zod'

export const signInFormSchema = z.object({
  email: z.string().email(),
})

export type SignInFormSchema = z.infer<typeof signInFormSchema>

export const signUpFormSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  contact: z.string().min(10).max(10),
})

export type SignUpFormSchema = z.infer<typeof signUpFormSchema>

export const createLibrarySchema = z.object({
  libraryName: z.string().min(3).max(30),
  ownerName: z.string().min(3).max(100),
  ownerEmail: z.string().email(),
  ownerContact: z.string().min(10).max(10),
})

export type CreateLibraryData = z.infer<typeof createLibrarySchema>

export const onboardAdminSchema = z.object({
  adminName: z.string().min(3).max(100),
  adminEmail: z.string().email(),
  adminContact: z.string().min(10).max(10),
  libraryId: z.string(),
})

export type OnboardAdminData = z.infer<typeof onboardAdminSchema>

export const addBookSchema = z.object({
  email: z.string().email(),
  isbn: z.string().min(1, 'ISBN is required'),
  title: z.string().min(1, 'Title is required'),
  authors: z.string().min(1, 'Authors are required'),
  publisher: z.string().min(1, 'Publisher is required'),
  version: z.string().min(1, 'Version is required'),
})

export type AddBookData = z.infer<typeof addBookSchema>
