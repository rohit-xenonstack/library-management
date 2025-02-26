import { useEffect, useMemo, useState } from 'react'

import { api } from '../api/config'
import Card from '../components/card'
import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/owner-dashboard.module.scss'
import { ROLE } from '../lib/constants'

interface Library {
  library_id: string
  name: string
  owner_name: string
  owner_email: string
  total_books: number
}

interface Admin {
  user_id: string
  name: string
  email: string
  contact: string
  role: ROLE.ADMIN
  library_id: string
}

export function OwnerDashboard() {
  const [libraries, setLibraries] = useState<Library[]>([])
  const [admins, setAdmins] = useState<Admin[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const auth = useAuth()

  const sortedLibraries = useMemo(() => {
    return [...libraries].sort((a, b) => {
      if (a.library_id === auth.user?.library_id) return -1
      if (b.library_id === auth.user?.library_id) return 1
      return 0
    })
  }, [libraries, auth.user?.library_id])

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [librariesRes, adminsRes] = await Promise.all([
          api.get('protected/owner/libraries').json<{ libraries: Library[] }>(),
          api
            .post('protected/owner/admins', {
              json: {
                library_id: auth.user?.library_id,
              },
            })
            .json<{ admins: Admin[] }>(),
        ])

        setLibraries(librariesRes.libraries)
        setAdmins(adminsRes.admins)
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

    fetchData()
  }, [auth.user?.library_id])

  if (isLoading) {
    return <div className={styles.loading}>Loading...</div>
  }

  if (error) {
    return <div className={styles.error}>{error}</div>
  }

  return (
    <div className={styles.container}>
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Libraries</h2>
        <div className={styles.grid}>
          {sortedLibraries.map((library) => (
            <Card
              key={library.library_id}
              className={`${styles.card} ${
                library.library_id === auth.user?.library_id
                  ? styles.currentLibrary
                  : ''
              }`}
            >
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>
                  {library.name}
                  {library.library_id === auth.user?.library_id && (
                    <span className={styles.currentBadge}>Owned Library</span>
                  )}
                </h3>
                <p className={styles.cardSubTitle}>{library.library_id}</p>
              </div>
              <div className={styles.cardContent}>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Owner</span>
                  <span className={styles.statValue}>
                    {library.owner_name || auth.user?.name}
                  </span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Total Books</span>
                  <span className={styles.statValue}>
                    {library.total_books || 0}
                  </span>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>My Library Admins</h2>
        <div className={styles.grid}>
          {admins.map((admin) => (
            <Card key={admin.user_id} className={styles.card}>
              <div className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>{admin.name}</h3>
                <h1 className={styles.cardSubTitle}>{admin.user_id}</h1>
              </div>
              <div className={styles.cardContent}>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Email</span>
                  <span className={styles.statValue}>{admin.email}</span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Contact</span>
                  <span className={styles.statValue}>{admin.contact}</span>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </section>
    </div>
  )
}
