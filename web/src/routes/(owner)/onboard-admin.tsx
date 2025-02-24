import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

import { onboardAdmin } from '../../api/owner'
import { useAuth } from '../../hook/use-auth'
import { fallback } from '../../lib/constants'
import { onboardAdminSchema } from '../../lib/schema'
import styles from '../../styles/modules/onboard-admin.module.scss'

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
    adminName: '',
    adminEmail: '',
    adminContact: '',
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
        libraryId: user?.library_id as string,
      })

      const response = await onboardAdmin(validatedData)

      if (response.status === 'success') {
        setSuccess('Admin onboarded successfully!')
        setFormData({
          adminName: '',
          adminEmail: '',
          adminContact: '',
        })
        setTimeout(() => {
          navigate({ to: '/' })
        }, 2000)
      } else {
        setError(response.payload || 'Failed to onboard admin')
      }
    } catch (err) {
      if (err instanceof z.ZodError) {
        setError(err.errors[0].message)
      } else {
        setError('Something went wrong')
      }
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
    <div className={styles.container}>
      <h1 className={styles.title}>Onboard New Admin</h1>

      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.formGroup}>
          <label htmlFor='adminName' className={styles.label}>
            Admin Name
          </label>
          <input
            type='text'
            id='adminName'
            name='adminName'
            value={formData.adminName}
            onChange={handleChange}
            className={styles.input}
            required
          />
        </div>

        <div className={styles.formGroup}>
          <label htmlFor='adminEmail' className={styles.label}>
            Admin Email
          </label>
          <input
            type='email'
            id='adminEmail'
            name='adminEmail'
            value={formData.adminEmail}
            onChange={handleChange}
            className={styles.input}
            required
          />
        </div>

        <div className={styles.formGroup}>
          <label htmlFor='adminContact' className={styles.label}>
            Contact Number
          </label>
          <input
            type='tel'
            id='adminContact'
            name='adminContact'
            value={formData.adminContact}
            onChange={handleChange}
            className={styles.input}
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
          {isLoading ? 'Onboarding...' : 'Onboard Admin'}
        </button>
      </form>
    </div>
  )
}
