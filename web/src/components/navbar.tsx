import { Link, useNavigate } from '@tanstack/react-router'
import clsx from 'clsx'
import { useState } from 'react'

import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/navbar.module.scss'
import type { FileRouteTypes } from '../routeTree.gen'

type NavigationItem = {
  to: FileRouteTypes['to']
  name: string
  roles?: string[]
}

const navigation: NavigationItem[] = [
  {
    to: '/',
    name: 'Dashboard',
    roles: ['reader', 'admin', 'owner'],
  },
  {
    to: '/add-book',
    name: `Add Book`,
    roles: ['admin'],
  },
  {
    to: '/issue-requests',
    name: 'Issue Requests',
    roles: ['admin'],
  },
  {
    to: '/onboard-admin',
    name: `Onboard Admin`,
    roles: ['owner'],
  },
  {
    to: '/create-library',
    name: 'Create Library',
    roles: ['owner'],
  },
]

export function Navigation() {
  const [isOpen, setIsOpen] = useState(false)
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const filteredNavigation = navigation.filter(
    (item) => !item.roles || (user && item.roles.includes(user.role)),
  )

  const handleSignOut = () => {
    logout()
    setIsOpen(false)
    navigate({ to: '/' })
  }
  const handleSignIn = () => {
    setIsOpen(false)
    navigate({ to: '/sign-in' })
  }

  return (
    <header className={styles.navbar}>
      <div className={styles.container}>
        <Link to='/' className={styles.logo}>
          Library App
        </Link>

        <button
          className={styles.hamburger}
          onClick={() => setIsOpen(!isOpen)}
          aria-label='Toggle navigation'
        >
          <span className={clsx(styles.hamburgerIcon, isOpen && styles.open)} />
        </button>

        <nav className={clsx(styles.nav, isOpen && styles.open)}>
          {filteredNavigation.map(({ to, name }) => (
            <Link
              key={to}
              to={to}
              className={styles.link}
              activeProps={{ className: styles.active }}
              activeOptions={{ exact: to === '/' }}
              onClick={() => setIsOpen(false)}
            >
              {name}
            </Link>
          ))}
          {user ? (
            <button className={styles.button} onClick={handleSignOut}>
              Sign Out
            </button>
          ) : (
            <button className={styles.button} onClick={handleSignIn}>
              Sign In
            </button>
          )}
        </nav>
        <button
          className={styles.backdrop}
          onClick={() => setIsOpen(false)}
          onKeyDown={(e) => e.key === 'Enter' && setIsOpen(false)}
          aria-label='Close menu'
        />
      </div>
    </header>
  )
}
