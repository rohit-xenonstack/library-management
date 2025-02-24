import { createFileRoute } from '@tanstack/react-router'

import { AdminDashboard } from '../components/admin-dashboard'
import { OwnerDashboard } from '../components/owner-dashboard'
import { ReaderDashboard } from '../components/reader-dashboard'
import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/dashboard.module.scss'

export const Route = createFileRoute('/')({
  component: Dashboard,
})

function Dashboard() {
  const auth = useAuth()
  return (
    <>
      <div className={styles.greetingBanner}>
        <h1 className={styles.greeting}>
          Welcome back,{' '}
          <span className={styles.userName}>{auth.user?.name || 'Guest'}</span>
        </h1>
        <p className={styles.role}>
          {auth.user?.role && `Logged in as ${auth.user.role}`}
        </p>
      </div>
      {auth.user?.role === 'owner' && <OwnerDashboard />}
      {auth.user?.role === 'admin' && <AdminDashboard />}
      {auth.user?.role === 'reader' && <ReaderDashboard />}
    </>
  )
}
