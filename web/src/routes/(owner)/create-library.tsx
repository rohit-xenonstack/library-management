import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'
import type { FormEvent } from 'react'

import { createLibrary } from '../../api/owner'
import { fallback } from '../../lib/constants'
import { createLibrarySchema } from '../../lib/schema'
import styles from '../../styles/modules/create-library.module.scss'
import type { CreateLibraryData } from '../../lib/schema'

export const Route = createFileRoute('/(owner)/create-library')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context, search }) => {
    if (!context.auth.user) {
      throw redirect({
        to: '/sign-in',
      })
    }
    if (context.auth.user && context.auth.user.role !== 'owner') {
      throw redirect({
        to: search.redirect || fallback,
      })
    }
  },
  component: CreateLibrary,
})

function CreateLibrary() {
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState<CreateLibraryData>({
    libraryName: '',
    ownerName: '',
    ownerEmail: '',
    ownerContact: '',
  })
  const [errors, setErrors] = useState<Partial<CreateLibraryData>>({})
  const [formError, setFormError] = useState<string | null>(null)
  const [formSuccess, setFormSuccess] = useState<string | null>(null)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsLoading(true)
    setErrors({})
    setFormError(null)

    try {
      const result = createLibrarySchema.safeParse(formData)
      if (!result.success) {
        const fieldErrors: Partial<CreateLibraryData> = {}
        for (const issue of result.error.issues) {
          const path = issue.path[0] as keyof CreateLibraryData
          fieldErrors[path] = issue.message
        }
        setErrors(fieldErrors)
        return
      }

      const response = await createLibrary(result.data)
      if (response.status === 'success') {
        setFormSuccess('Library created successfully')
        setFormData({
          libraryName: '',
          ownerName: '',
          ownerEmail: '',
          ownerContact: '',
        })
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      } else {
        setFormError(response.payload || 'Failed to create library')
      }
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target

    setFormData((prev: CreateLibraryData) => ({ ...prev, [name]: value }))
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Create Library</h1>
      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.formGroup}>
          <label htmlFor='library_name' className={styles.label}>
            Library Name
          </label>
          <input
            id='library_name'
            name='libraryName'
            type='text'
            className={styles.input}
            value={formData.libraryName}
            onChange={handleChange}
            disabled={isLoading}
            required
          />
          {errors.libraryName && (
            <div className={styles.error}>{errors.libraryName}</div>
          )}
        </div>

        <div className={styles.formGroup}>
          <label htmlFor='name' className={styles.label}>
            Owner Name
          </label>
          <input
            id='name'
            name='ownerName'
            type='text'
            className={styles.input}
            value={formData.ownerName}
            onChange={handleChange}
            disabled={isLoading}
            required
          />
          {errors.ownerName && (
            <div className={styles.error}>{errors.ownerName}</div>
          )}
        </div>

        <div className={styles.formGroup}>
          <label htmlFor='email' className={styles.label}>
            Owner Email
          </label>
          <input
            id='email'
            name='ownerEmail'
            type='email'
            className={styles.input}
            value={formData.ownerEmail}
            onChange={handleChange}
            disabled={isLoading}
            required
          />
          {errors.ownerEmail && (
            <div className={styles.error}>{errors.ownerEmail}</div>
          )}
        </div>

        <div className={styles.formGroup}>
          <label htmlFor='contact' className={styles.label}>
            Owner Contact Number
          </label>
          <input
            id='contact'
            name='ownerContact'
            type='tel'
            className={styles.input}
            value={formData.ownerContact}
            onChange={handleChange}
            disabled={isLoading}
            required
          />
          {errors.ownerContact && (
            <div className={styles.error}>{errors.ownerContact}</div>
          )}
        </div>

        {formError && (
          <div className={`${styles.formMessage} ${styles.error}`}>
            {formError}
          </div>
        )}
        {formSuccess && (
          <div className={`${styles.formMessage} ${styles.success}`}>
            {formSuccess}
          </div>
        )}

        <button type='submit' className={styles.button} disabled={isLoading}>
          {isLoading ? 'Creating...' : 'Create Library'}
        </button>
      </form>
    </div>
  )
}
