import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

import { registerReader } from '../../api/auth'
import { fallback } from '../../lib/constants'
import styles from '../../styles/modules/register.module.scss'

export const Route = createFileRoute('/(reader)/register/$libID')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context, search }) => {
    if (context.auth.user) {
      throw redirect({
        to: search.redirect || fallback,
      })
    }
  },
  component: RegisterReader,
})

function RegisterReader() {
  const { libID } = Route.useParams()
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    contact: '',
  })

  if (!libID) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>
          Library ID is required. Please use a valid registration link.
        </div>
      </div>
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    setSuccess('')

    try {
      const response = await registerReader({
        ...formData,
        library_id: libID,
      })

      if (response.status === 'success') {
        setSuccess('Registration successful! Redirecting to login...')
        setTimeout(() => {
          navigate({ to: '/sign-in' })
        }, 2000)
      } else {
        setError(response.payload || 'Registration failed')
      }
    } catch (err) {
      console.error(err)
      setError('An error occurred during registration')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Reader Registration</h1>
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.formGroup}>
          <label htmlFor='libraryId'>Library ID</label>
          <input type='text' id='libraryId' value={libID} disabled />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='name'>Full Name</label>
          <input
            type='text'
            id='name'
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='email'>Email</label>
          <input
            type='email'
            id='email'
            value={formData.email}
            onChange={(e) =>
              setFormData({ ...formData, email: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>
        <div className={styles.formGroup}>
          <label htmlFor='contact'>Contact Number</label>
          <input
            type='tel'
            id='contact'
            value={formData.contact}
            onChange={(e) =>
              setFormData({ ...formData, contact: e.target.value })
            }
            disabled={isLoading}
            required
          />
        </div>

        {error && <div className={styles.error}>{error}</div>}
        {success && <div className={styles.success}>{success}</div>}

        <button
          type='submit'
          className={styles.submitButton}
          disabled={isLoading}
        >
          {isLoading ? 'Registering...' : 'Register'}
        </button>
      </form>
    </div>
  )
}
