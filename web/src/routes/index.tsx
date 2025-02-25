import { createFileRoute, Link } from '@tanstack/react-router'
import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/dashboard.module.scss'
import { SignupLink } from '../components/signup-link'
import { OwnerDashboard } from '../components/owner-dashboard'
import { AdminDashboard } from '../components/admin-dashboard'
import { ReaderDashboard } from '../components/reader-dashboard'

const quotes = [
  {
    text: 'Life, although it may only be an accumulation of anguish, is dear to me, and I will defend it.',
    author: 'Frankenstein by Mary Shelley',
  },
  {
    text: 'He was always late on principle, his principle being that punctuality is the thief of time.',
    author: 'The Picture of Dorian Gray by Oscar Wilde',
  },
  {
    text: '“He who controls the past controls the future. He who controls the present controls the past.',
    author: '1984 by George Orwell',
  },
]

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  const auth = useAuth()
  const randomQuote = quotes[Math.floor(Math.random() * quotes.length)]

  if (auth.user) {
    return (
      <>
        <div className={styles.greetingBanner}>
          <h1 className={styles.greeting}>
            Welcome back,{' '}
            <span className={styles.userName}>
              {auth.user?.name || 'Guest'}
            </span>
          </h1>
          <p className={styles.role}>
            {auth.user?.role && `Logged in as ${auth.user.role}`}
          </p>
          {auth.user?.role === 'admin' && <SignupLink />}
        </div>

        {auth.user?.role === 'owner' && <OwnerDashboard />}
        {auth.user?.role === 'admin' && <AdminDashboard />}
        {auth.user?.role === 'reader' && <ReaderDashboard />}
      </>
    )
  }

  return (
    <main className={styles.homeContainer}>
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <h1 className={styles.title}>
            Welcome to <span className={styles.highlight}>Library Manager</span>
          </h1>
          <p className={styles.subtitle}>
            Your modern solution for efficient library management
          </p>
          <div className={styles.quote}>
            <blockquote>"{randomQuote.text}"</blockquote>
            <cite>— {randomQuote.author}</cite>
          </div>
          <div className={styles.actions}>
            <Link to='/sign-in' className={styles.primaryButton}>
              Sign In
            </Link>
          </div>
        </div>
      </section>
    </main>
  )
}
