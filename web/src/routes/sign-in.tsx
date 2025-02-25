import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'
import type { FormEvent } from 'react'

import { signIn } from '../api/auth'
import { useAuth } from '../hook/use-auth'
import { fallback } from '../lib/constants'
import { signInFormSchema } from '../lib/schema'
import styles from '../styles/form.module.scss'

export const Route = createFileRoute('/sign-in')({
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
  component: SignInComponent,
})

function SignInComponent() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({ email: '' })
  const [fieldError, setFieldError] = useState<string | null>(null)
  const [formError, setFormError] = useState<string | null>(null)
  const [formSuccess, setFormSuccess] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setIsLoading(true)
    setFieldError(null)
    setFormError(null)

    try {
      const result = signInFormSchema.safeParse(formData)
      if (!result.success) {
        setFieldError(
          result.error.flatten().fieldErrors.email?.[0] || 'Invalid email',
        )
        return
      }

      const response = await signIn(result.data)
      if (response.status === 'success' && response.payload.user) {
        setFormData({
          email: '',
        })
        setFormSuccess('Signed in successfully. Redirecting to dashboard...')
        login(response.payload.user)
        setTimeout(() => {
          navigate({ to: fallback })
        }, 2000)
      } else {
        setFormError("Couldn't sign in: " + response.payload.message)
      }
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'An error occurred')
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
        <h1 className={styles.title}>Sign In</h1>
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label htmlFor='email'>Email</label>
            <input
              className={styles.input}
              id='email'
              type='email'
              name='email'
              placeholder='Enter your email'
              required
              disabled={isLoading}
              onChange={handleChange}
            />
            {fieldError && <div className={styles.error}>{fieldError}</div>}
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
            {isLoading ? <span className={styles.loader} /> : 'Sign In'}
          </button>
        </form>
      </div>
    </main>
  )
}
