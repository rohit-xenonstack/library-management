import { useState, useEffect } from 'react'

import { checkAvailability, requestBook } from '../api/reader'
import { SearchBar } from '../components/search-bar'
import { useAuth } from '../hook/use-auth'
import styles from '../styles/modules/reader-dashboard.module.scss'
import { BookData } from '../types/data'
import { getBooks, searchBooks } from '../api/shared'
import { SearchBookRequest } from '../types/request'
import { HTTPError } from 'ky'

export function ReaderDashboard() {
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(false)
  const [books, setBooks] = useState<BookData[]>([])
  const [latestBooks, setLatestBooks] = useState<BookData[]>([])
  const [error, setError] = useState('')
  const [latestBooksError, setLatestBooksError] = useState('')

  useEffect(() => {
    const fetchLatestBooks = async () => {
      try {
        const response = await getBooks()
        if (response.status === 'success') {
          setLatestBooks(response.books || [])
        }
      } catch (err) {
        setLatestBooksError(
          err instanceof HTTPError
            ? 'Failed: ' + (await err.response.json()).message
            : 'something went wrong, please again try later',
        )
      }
    }

    fetchLatestBooks()
  }, [])

  const handleSearch = async (data: SearchBookRequest) => {
    setIsLoading(true)
    setError('')
    try {
      const response = await searchBooks(data)
      if (response.status === 'success') {
        setBooks(response.books || [])
      }
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

      <section className={styles.latestBooksSection}>
        <h2>Latest Added Books</h2>
        {latestBooksError && (
          <div className={styles.error}>{latestBooksError}</div>
        )}
        <div className={styles.booksGrid}>
          {latestBooks.map((book) => (
            <BookCard
              key={book.isbn}
              book={book}
              userEmail={user?.email || ''}
            />
          ))}
        </div>
      </section>
    </div>
  )
}

interface BookCardProps {
  book: BookData
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
        setAvailableDate(response.date || '')
      }
    } catch (err) {
      setRequestStatus({
        message:
          err instanceof HTTPError
            ? 'Failed: ' + (await err.response.json()).message
            : 'something went wrong, please again try later',
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
      const response = await requestBook({ isbn: book.isbn, email: userEmail })
      if (response.status === 'success') {
        setRequestStatus({
          message: 'Book request submitted successfully!',
          type: 'success',
        })
      }
    } catch (err) {
      setRequestStatus({
        message:
          err instanceof HTTPError
            ? 'Failed: ' + (await err.response.json()).message
            : 'something went wrong, please again try later',
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
