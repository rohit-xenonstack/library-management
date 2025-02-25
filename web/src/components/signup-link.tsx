import { useState } from 'react'
import styles from '../styles/modules/signup-link.module.scss'
import { useAuth } from '../hook/use-auth'

export function SignupLink() {
  const { user } = useAuth()
  const [copied, setCopied] = useState(false)

  const signupLink = `${window.location.origin}/register/${user?.library_id}`

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(signupLink)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.linkBox}>
        <div className={styles.label}>Reader Signup Link:</div>
        <div className={styles.linkWrapper}>
          <code className={styles.link}>{signupLink}</code>
          <button
            className={styles.copyButton}
            onClick={handleCopy}
          >
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>
      </div>
    </div>
  )
}
