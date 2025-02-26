import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { z } from 'zod'

import { getBook, updateBook } from '../../api/admin'
import styles from '../../styles/modules/update-book.module.scss'
import type { BookData } from '../../types/data'
import { DASHBOARD, LOGIN_PAGE } from '../../lib/constants'
import { HTTPError } from 'ky'

export const Route = createFileRoute('/(admin)/$isbn/update')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({
        to: LOGIN_PAGE,
      })
    }
    if (context.auth.user.role !== 'admin') {
      throw redirect({
        to: DASHBOARD,
      })
    }
  },
  component: UpdateBook,
})

function UpdateBook() {
  const { isbn } = Route.useParams()
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [book, setBook] = useState<BookData | null>(null)
  const [formData, setFormData] = useState({
    title: '',
    authors: '',
    publisher: '',
    version: '',
  })

  useEffect(() => {
    const fetchBook = async () => {
      try {
        const response = await getBook(isbn)
        if (response.status === 'success' && response.book) {
          setBook(response.book)
          setFormData({
            title: response.book.title,
            authors: response.book.authors,
            publisher: response.book.publisher,
            version: response.book.version,
          })
        }
      } catch (err) {
        if (err instanceof HTTPError) {
          const res = await err.response.json()
          setError('Failed: ' + res.message)
        } else {
          setError('something went wrong, please again try later')
        }
      }
    }
    fetchBook()
  }, [isbn])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    setSuccess('')

    try {
      const response = await updateBook({
        isbn,
        ...formData,
      })
      if (response.status === 'success') {
        setSuccess('Book updated successfully!')
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      }
    } catch (err) {
      setError(
        err instanceof HTTPError
          ? 'Failed: ' + (await err.response.json()).message
          : 'something went wrong, please again try later',
      )
    } finally {
      setIsLoading(false)
    }
  }

  if (!book) {
    return <div className={styles.loading}>Loading...</div>
  }

  return (
    <main>
      <div className={styles.container}>
        <h1 className={styles.title}>Update Book</h1>
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label htmlFor='isbn'>ISBN</label>
            <input type='text' id='isbn' value={isbn} disabled />
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
              min={1}
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
              min={3}
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
              min={3}
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
              min={1}
            />
          </div>

          {error && (
            <div className={`${styles.formMessage} ${styles.error}`}>
              {error}
            </div>
          )}

          {success && (
            <div className={`${styles.formMessage} ${styles.success}`}>
              {success}
            </div>
          )}

          <div className={styles.actions}>
            <button
              type='button'
              className={styles.cancelButton}
              onClick={() => navigate({ to: '/' })}
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type='submit'
              className={styles.submitButton}
              disabled={isLoading}
            >
              {isLoading ? 'Updating...' : 'Update Book'}
            </button>
          </div>
        </form>
      </div>
    </main>
  )
}
