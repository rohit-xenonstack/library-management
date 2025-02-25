import { useState, useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import { decreaseBookCount, searchBooks, getLatestBooks } from '../api/admin'
import { SearchBar } from '../components/search-bar'
import styles from '../styles/modules/admin-dashboard.module.scss'
import type { Book } from '../api/admin'

export function AdminDashboard() {
  const [isLoading, setIsLoading] = useState(false)
  const [books, setBooks] = useState<Book[]>([])
  const [latestBooks, setLatestBooks] = useState<Book[]>([])
  const [error, setError] = useState('')
  const [latestBooksError, setLatestBooksError] = useState('')

  useEffect(() => {
    const fetchLatestBooks = async () => {
      try {
        const response = await getLatestBooks()
        if (response.status === 'success' && response.payload) {
          setLatestBooks(response.payload)
        } else {
          setLatestBooksError('Failed to fetch latest books')
        }
      } catch (err) {
        console.error(err)
        setLatestBooksError('An error occurred while fetching latest books')
      }
    }

    fetchLatestBooks()
  }, [])

  const handleSearch = async (
    searchString: string,
    field: 'title' | 'authors' | 'publisher',
  ) => {
    setIsLoading(true)
    setError('')
    try {
      console.log(searchString, field)
      const response = await searchBooks(searchString, field)
      if (response.status === 'success') {
        setBooks(response.payload)
      } else {
        setError('Failed to fetch books')
      }
    } catch (err) {
      console.error(err)
      setError('An error occurred while searching books')
    } finally {
      setIsLoading(false)
    }
  }

  const handleDecreaseCount = async (isbn: string) => {
    try {
      const response = await decreaseBookCount(isbn)
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
      } else {
        setError('Failed to decrease book count')
      }
    } catch (err) {
      console.error(err)
      setError('An error occurred while updating book count')
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
            />
          ))}
        </div>
      </section>
    </div>
  )
}

interface BookCardProps {
  book: Book
  onDecreaseCount: () => void
}

function BookCard({ book, onDecreaseCount }: BookCardProps) {
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
