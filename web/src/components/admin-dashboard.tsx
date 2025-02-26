import { useState, useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import { decreaseBookCount } from '../api/admin'
import { SearchBar } from '../components/search-bar'
import styles from '../styles/modules/admin-dashboard.module.scss'
import type { BookData } from '../types/data'
import { getBooks, searchBooks } from '../api/shared'
import { SearchBookRequest } from '../types/request'
import { HTTPError } from 'ky'

export function AdminDashboard() {
  const [isLoading, setIsLoading] = useState(false)
  const [books, setBooks] = useState<BookData[]>([])
  const [latestBooks, setLatestBooks] = useState<BookData[]>([])
  const [error, setError] = useState('')
  const [bookCardError, setBookCardError] = useState('')
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

  const handleDecreaseCount = async (isbn: string) => {
    try {
      const response = await decreaseBookCount({ isbn: isbn })
      if (response.status === 'success') {
        setBooks((prevBooks) =>
          prevBooks.map((book) =>
            book.isbn === isbn
              ? {
                  ...book,
                  available_copies: book.available_copies - 1,
                  total_copies: book.total_copies - 1,
                }
              : book,
          ),
        )
        setLatestBooks((prevBooks) =>
          prevBooks.map((latestBook) =>
            latestBook.isbn === isbn
              ? {
                  ...latestBook,
                  available_copies: latestBook.available_copies - 1,
                  total_copies: latestBook.total_copies - 1,
                }
              : latestBook,
          ),
        )
      }
    } catch (err) {
      setBookCardError(
        err instanceof HTTPError
          ? 'Failed: ' + (await err.response.json()).message
          : 'something went wrong, please again try later',
      )
      setTimeout(() => setBookCardError(''), 3000)
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
              onDecreaseCount={() => handleDecreaseCount(book.isbn)}
              bookCardError={bookCardError}
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
              onDecreaseCount={() => handleDecreaseCount(book.isbn)}
              bookCardError={bookCardError}
            />
          ))}
        </div>
      </section>
    </div>
  )
}

interface BookCardProps {
  book: BookData
  onDecreaseCount: () => void
  bookCardError: string
}

function BookCard({ book, onDecreaseCount, bookCardError }: BookCardProps) {
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
      </div>
      {bookCardError && <div className={styles.error}>{bookCardError}</div>}
      <div className={styles.bookActions}>
        <Link
          to='/$isbn/update'
          params={{ isbn: book.isbn }}
          className={styles.editButton}
        >
          Edit Book
        </Link>
        <button
          className={styles.decreaseButton}
          onClick={onDecreaseCount}
          disabled={book.available_copies === 0}
        >
          Decrease Count
        </button>
      </div>
    </div>
  )
}
