import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

import { readerRegister } from '../../api/auth'
import styles from '../../styles/form.module.scss'
import { DASHBOARD } from '../../lib/constants'
import { HTTPError } from 'ky'

export const Route = createFileRoute('/(reader)/register/$libID')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context, search }) => {
    if (context.auth.user) {
      throw redirect({
        to: search.redirect || DASHBOARD,
      })
    }
  },
  component: RegisterReader,
})

function RegisterReader() {
  const { libID } = Route.useParams()
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [formError, setFormError] = useState('')
  const [formSuccess, setFormSuccess] = useState('')
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
    setFormError('')
    setFormSuccess('')

    try {
      const response = await readerRegister({
        ...formData,
        library_id: libID,
      })

      if (response.status === 'success') {
        setFormSuccess('Registration successful! Redirecting to login...')
        setTimeout(() => {
          navigate({ to: '/sign-in' })
        }, 2000)
      }
    } catch (err) {
      setFormError(
        err instanceof HTTPError
          ? 'Failed: ' + (await err.response.json()).message
          : 'something went wrong, please again try later',
      )
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <main>
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
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
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

          <button className={styles.button} type='submit' disabled={isLoading}>
            {isLoading ? <span className={styles.loader} /> : 'Register'}
          </button>
        </form>
      </div>
    </main>
  )
}
