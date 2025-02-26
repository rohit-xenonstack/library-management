import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

import { onboardAdmin } from '../../api/owner'
import { useAuth } from '../../hook/use-auth'
import { fallback } from '../../lib/constants'
import { onboardAdminSchema } from '../../types/data'
import styles from '../../styles/form.module.scss'

export const Route = createFileRoute('/(owner)/onboard-admin')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({
        to: '/sign-in',
      })
    }
    if (context.auth.user && context.auth.user.role !== 'owner') {
      throw redirect({
        to: fallback,
      })
    }
  },
  component: OnboardAdmin,
})

function OnboardAdmin() {
  const { user } = useAuth()
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    contact: '',
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
      const validatedData = onboardAdminSchema.parse({
        ...formData,
        library_id: user?.library_id as string,
      })

      const response = await onboardAdmin(validatedData)

      if (response.status === 'success') {
        setSuccess('Admin onboarded successfully!')
        setFormData({
          name: '',
          email: '',
          contact: '',
        })
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      } else {
        setError('Failed: ' + response.message)
      }
    } catch (err) {
      setError('Something went wrong: ' + err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }))
  }

  return (
    <main>
      <div className={styles.container}>
        <h1 className={styles.title}>Onboard New Admin</h1>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor='name' className={styles.label}>
              Admin Name
            </label>
            <input
              type='text'
              id='name'
              name='name'
              value={formData.name}
              onChange={handleChange}
              className={styles.input}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor='email' className={styles.label}>
              Admin Email
            </label>
            <input
              type='email'
              id='email'
              name='email'
              value={formData.email}
              onChange={handleChange}
              className={styles.input}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor='contact' className={styles.label}>
              Contact Number
            </label>
            <input
              type='tel'
              id='contact'
              name='contact'
              value={formData.contact}
              onChange={handleChange}
              className={styles.input}
              required
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

          <button type='submit' className={styles.button} disabled={isLoading}>
            {isLoading ? 'Onboarding...' : 'Onboard Admin'}
          </button>
        </form>
      </div>
    </main>
  )
}
