import { useState } from 'react'

import { checkAvailability, requestBook, searchBooks } from '../api/reader'
import { SearchBar } from '../components/search-bar'
import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/reader-dashboard.module.scss'
import type { Book } from '../api/admin'

type SearchField = 'title' | 'authors' | 'publisher'

export function ReaderDashboard() {
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(false)
  const [books, setBooks] = useState<Book[]>([])
  const [error, setError] = useState('')

  const handleSearch = async (searchString: string, field: SearchField) => {
    setIsLoading(true)
    setError('')
    try {
      const response = await searchBooks(searchString, field)
      if (response.status === 'success') {
        setBooks(response.payload)
      } else {
        setError('Failed: ' + response.payload)
      }
    } catch (err) {
      setError('An error occurred while searching books: ' + err)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={styles.container}>
      <SearchBar onSearch={handleSearch} isLoading={isLoading} />

      {error && <div className={styles.error}>{error}</div>}

      {!isLoading && books.length === 0 ? (
        <div className={styles.noResults}>
          <p>No books found. Try adjusting your search criteria.</p>
        </div>
      ) : (
        <div className={styles.booksGrid}>
          {books.map((book) => (
            <BookCard
              key={book.isbn}
              book={book}
              userEmail={user?.email || ''}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface BookCardProps {
  book: Book
  userEmail: string
}

function BookCard({ book, userEmail }: BookCardProps) {
  const [requestStatus, setRequestStatus] = useState<{
    message: string
    type: 'success' | 'error' | null
  }>({ message: '', type: null })
  const [isRequesting, setIsRequesting] = useState(false)
  const [availableDate, setAvailableDate] = useState<string | null>(null)
  const [isCheckingAvailability, setIsCheckingAvailability] = useState(false)

  const handleCheckAvailability = async () => {
    setIsCheckingAvailability(true)
    try {
      const response = await checkAvailability(book.isbn)
      if (response.status === 'success') {
        setAvailableDate(response.payload)
      }
    } catch (err) {
      setRequestStatus({
        message: 'Failed to check availability: ' + err,
        type: 'error',
      })
    } finally {
      setIsCheckingAvailability(false)
    }
  }

  const handleRequestBook = async () => {
    if (!userEmail) return
    setIsRequesting(true)
    try {
      const response = await requestBook(book.isbn, userEmail)
      if (response.status === 'success') {
        setRequestStatus({
          message: 'Book request submitted successfully!',
          type: 'success',
        })
      } else {
        setRequestStatus({
          message: 'Failed: ' + response.payload,
          type: 'error',
        })
      }
    } catch (err) {
      setRequestStatus({
        message: 'An error occurred while requesting the book: ' + err,
        type: 'error',
      })
    } finally {
      setIsRequesting(false)
    }
  }

  return (
    <div className={styles.bookCard}>
      <div className={styles.bookInfo}>
        <h3 className={styles.bookTitle}>{book.title}</h3>
        <p className={styles.bookDetails}>
          <span>ISBN:</span> {book.isbn}
        </p>
        <p className={styles.bookDetails}>
          <span>Authors:</span> {book.authors}
        </p>
        <p className={styles.bookDetails}>
          <span>Publisher:</span> {book.publisher}
        </p>
        <p className={styles.bookDetails}>
          <span>Version:</span> {book.version}
        </p>
        <p className={styles.bookCopies}>
          <span>Available:</span> {book.available_copies} of {book.total_copies}
        </p>
        {availableDate && (
          <p className={styles.availabilityInfo}>
            <span>Expected availability:</span>{' '}
            {new Date(availableDate).toLocaleDateString()}
          </p>
        )}
      </div>

      <div className={styles.bookActions}>
        {requestStatus.type && (
          <div
            className={`${styles.requestStatus} ${
              requestStatus.type === 'success' ? styles.success : styles.error
            }`}
          >
            {requestStatus.message}
          </div>
        )}
        <button
          className={styles.requestButton}
          onClick={handleRequestBook}
          disabled={book.available_copies === 0 || isRequesting}
        >
          {isRequesting ? 'Requesting...' : 'Request Book'}
        </button>
        {book.available_copies === 0 && (
          <button
            className={styles.checkButton}
            onClick={handleCheckAvailability}
            disabled={isCheckingAvailability}
          >
            {isCheckingAvailability ? 'Checking...' : 'Check Availability'}
          </button>
        )}
      </div>
    </div>
  )
}
