import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'
import type { FormEvent } from 'react'

import { createLibrary } from '../../api/owner'
import { fallback } from '../../lib/constants'
import styles from '../../styles/form.module.scss'
import {
  CreateLibraryWithOwnerData,
  createLibraryWithOwnerSchema,
} from '../../types/data'

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
  const [formData, setFormData] = useState<CreateLibraryWithOwnerData>({
    library_name: '',
    name: '',
    email: '',
    contact: '',
  })
  const [errors, setErrors] = useState<Partial<CreateLibraryWithOwnerData>>({})
  const [formError, setFormError] = useState<string | null>(null)
  const [formSuccess, setFormSuccess] = useState<string | null>(null)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsLoading(true)
    setErrors({})
    setFormError(null)

    try {
      const result = createLibraryWithOwnerSchema.safeParse(formData)
      if (!result.success) {
        const fieldErrors: Partial<CreateLibraryWithOwnerData> = {}
        for (const issue of result.error.issues) {
          const path = issue.path[0] as keyof CreateLibraryWithOwnerData
          fieldErrors[path] = issue.message
        }
        setErrors(fieldErrors)
        return
      }

      const response = await createLibrary(result.data)
      if (response.status === 'success') {
        setFormSuccess('Library created successfully')
        setFormData({
          library_name: '',
          name: '',
          email: '',
          contact: '',
        })
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      } else {
        setFormError('Failed: ' + response.message)
      }
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target

    setFormData((prev: CreateLibraryWithOwnerData) => ({
      ...prev,
      [name]: value,
    }))
  }

  return (
    <main>
      <div className={styles.container}>
        <h1 className={styles.title}>Create Library</h1>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor='library_name' className={styles.label}>
              Library Name
            </label>
            <input
              id='library_name'
              name='library_name'
              type='text'
              className={styles.input}
              value={formData.library_name}
              onChange={handleChange}
              disabled={isLoading}
              required
            />
            {errors.library_name && (
              <div className={styles.error}>{errors.library_name}</div>
            )}
          </div>

          <div className={styles.formGroup}>
            <label htmlFor='name' className={styles.label}>
              Owner Name
            </label>
            <input
              id='name'
              name='name'
              type='text'
              className={styles.input}
              value={formData.name}
              onChange={handleChange}
              disabled={isLoading}
              required
            />
            {errors.name && <div className={styles.error}>{errors.name}</div>}
          </div>

          <div className={styles.formGroup}>
            <label htmlFor='email' className={styles.label}>
              Owner Email
            </label>
            <input
              id='email'
              name='email'
              type='email'
              className={styles.input}
              value={formData.email}
              onChange={handleChange}
              disabled={isLoading}
              required
            />
            {errors.email && <div className={styles.error}>{errors.email}</div>}
          </div>

          <div className={styles.formGroup}>
            <label htmlFor='contact' className={styles.label}>
              Owner Contact Number
            </label>
            <input
              id='contact'
              name='contact'
              type='tel'
              className={styles.input}
              value={formData.contact}
              onChange={handleChange}
              disabled={isLoading}
              required
            />
            {errors.contact && (
              <div className={styles.error}>{errors.contact}</div>
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
    </main>
  )
}
