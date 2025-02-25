import { createFileRoute, redirect } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { z } from 'zod'

import {
  approveRequest,
  getIssueRequests,
  rejectRequest,
} from '../../api/admin'
import { useAuth } from '../../hook/use-auth'
import styles from '../../styles/modules/issue-requests.module.scss'
import type { IssueRequest } from '../../api/admin'

export const Route = createFileRoute('/(admin)/issue-requests')({
  validateSearch: z.object({
    redirect: z.string().optional().catch(''),
  }),
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({
        to: '/sign-in',
      })
    }
    if (context.auth.user.role !== 'admin') {
      throw redirect({
        to: '/',
      })
    }
  },
  component: IssueRequests,
})

function IssueRequests() {
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(true)
  const [requests, setRequests] = useState<IssueRequest[]>([])
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  useEffect(() => {
    fetchRequests()
  }, [])

  const fetchRequests = async () => {
    try {
      const response = await getIssueRequests()
      if (response.status === 'success') {
        setRequests(response.payload)
      } else {
        setError('Error: ' + response.payload)
      }
    } catch (err) {
      setError('Error: ' + err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleApprove = async (requestId: string) => {
    if (!user?.user_id) return
    try {
      const response = await approveRequest(requestId, user.user_id)
      if (response.status === 'success') {
        setSuccess('Request approved successfully')
        setRequests((prev) =>
          prev.filter((req) => req.request_id !== requestId),
        )
        setTimeout(() => setSuccess(''), 3000)
        fetchRequests()
      } else {
        setError('Failed to approve request')
      }
    } catch (err) {
      setError('An error occurred while approving request: ' + err)
    }
  }

  const handleReject = async (requestId: string) => {
    if (!user?.user_id) return
    try {
      const response = await rejectRequest(requestId, user.user_id)
      if (response.status === 'success') {
        setSuccess('Request rejected successfully')
        setRequests((prev) =>
          prev.filter((req) => req.request_id !== requestId),
        )
        setTimeout(() => setSuccess(''), 3000)
      } else {
        setError('Failed to reject request')
      }
    } catch (err) {
      setError('An error occurred while rejecting request: ' + err)
    }
  }

  if (isLoading) {
    return <div className={styles.loading}>Loading requests...</div>
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Issue Requests</h1>

      {error && <div className={styles.error}>{error}</div>}
      {success && <div className={styles.success}>{success}</div>}

      {requests.length === 0 ? (
        <div className={styles.noRequests}>No pending requests</div>
      ) : (
        <div className={styles.requestsGrid}>
          {requests.map((request) => (
            <RequestCard
              key={request.request_id}
              request={request}
              onApprove={() => handleApprove(request.request_id)}
              onReject={() => handleReject(request.request_id)}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface RequestCardProps {
  request: IssueRequest
  onApprove: () => void
  onReject: () => void
}

function RequestCard({ request, onApprove, onReject }: RequestCardProps) {
  return (
    <div className={styles.requestCard}>
      <div className={styles.requestInfo}>
        <h3 className={styles.bookTitle}>{request.book_title}</h3>
        <p className={styles.requestDetails}>
          <span>ISBN:</span> {request.isbn}
        </p>
        <p className={styles.requestDetails}>
          <span>Request Date:</span>{' '}
          {new Date(request.request_date).toLocaleDateString()}
        </p>
        <p className={styles.requestDetails}>
          <span>Reader ID:</span> {request.reader_id}
        </p>
        <p className={styles.availableCopies}>
          <span>Available Copies:</span> {request.available_copies}
        </p>
      </div>
      <div className={styles.requestActions}>
        <button
          className={styles.approveButton}
          onClick={onApprove}
          disabled={request.available_copies === 0}
        >
          Approve
        </button>
        <button className={styles.rejectButton} onClick={onReject}>
          Reject
        </button>
      </div>
    </div>
  )
}
