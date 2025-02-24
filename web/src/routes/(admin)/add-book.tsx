import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

import { addBook } from '../../api/admin'
import { useAuth } from '../../hook/use-auth'
import { addBookSchema } from '../../lib/schema'
import styles from '../../styles/modules/add-book.module.scss'

export const Route = createFileRoute('/(admin)/add-book')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({
        to: '/sign-in',
      })
    }
    if (context.auth.user && context.auth.user.role !== 'admin') {
      throw redirect({
        to: '/',
      })
    }
  },
  component: AddBook,
})

function AddBook() {
  const { user } = useAuth()
  const [formData, setFormData] = useState({
    isbn: '',
    title: '',
    authors: '',
    publisher: '',
    version: '',
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    setIsLoading(true)

    try {
      const validatedData = addBookSchema.parse({
        ...formData,
        email: user?.email as string,
      })

      const response = await addBook(validatedData)

      if (response.status === 'success') {
        setSuccess('Book added successfully!')
        setFormData({
          isbn: '',
          title: '',
          authors: '',
          publisher: '',
          version: '',
        })
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      } else {
        setError(response.payload || 'Failed to add book')
      }
    } catch (err) {
      console.error(err)
      setError('Please check the form fields')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Add New Book</h1>
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.formGroup}>
          <label htmlFor='isbn'>ISBN</label>
          <input
            type='text'
            id='isbn'
            value={formData.isbn}
            onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='title'>Title</label>
          <input
            type='text'
            id='title'
            value={formData.title}
            onChange={(e) =>
              setFormData({ ...formData, title: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='authors'>Authors</label>
          <input
            type='text'
            id='authors'
            value={formData.authors}
            onChange={(e) =>
              setFormData({ ...formData, authors: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='publisher'>Publisher</label>
          <input
            type='text'
            id='publisher'
            value={formData.publisher}
            onChange={(e) =>
              setFormData({ ...formData, publisher: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='version'>Version</label>
          <input
            type='text'
            id='version'
            value={formData.version}
            onChange={(e) =>
              setFormData({ ...formData, version: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>
        {error && (
          <div className={`${styles.formMessage} ${styles.error}`}>{error}</div>
        )}

        {success && (
          <div className={`${styles.formMessage} ${styles.success}`}>
            {success}
          </div>
        )}
        <button type='submit' className={styles.button} disabled={isLoading}>
          {isLoading ? 'Adding Book...' : 'Add Book'}
        </button>
      </form>
    </div>
  )
}
