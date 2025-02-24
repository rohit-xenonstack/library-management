import { useState } from 'react'

import styles from '../styles/modules/search-bar.module.scss'

type SearchField = 'title' | 'authors' | 'publisher'

interface SearchBarProps {
  onSearch: (searchString: string, field: SearchField) => void
  isLoading?: boolean
}

export function SearchBar({ onSearch, isLoading = false }: SearchBarProps) {
  const [searchString, setSearchString] = useState('')
  const [searchField, setSearchField] = useState<SearchField>('title')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchString.trim()) {
      onSearch(searchString.trim(), searchField)
    }
  }

  return (
    <form className={styles.searchContainer} onSubmit={handleSubmit}>
      <div className={styles.searchWrapper}>
        <select
          className={styles.searchSelect}
          value={searchField}
          onChange={(e) => setSearchField(e.target.value as SearchField)}
          disabled={isLoading}
        >
          <option value='title'>Title</option>
          <option value='authors'>Author</option>
          <option value='publisher'>Publisher</option>
        </select>
        <input
          type='search'
          className={styles.searchInput}
          placeholder='Search books...'
          value={searchString}
          onChange={(e) => setSearchString(e.target.value)}
          disabled={isLoading}
        />
        <button
          type='submit'
          className={styles.searchButton}
          disabled={isLoading}
        >
          {isLoading ? (
            <span className={styles.loader} />
          ) : (
            <SearchIcon className={styles.searchIcon} />
          )}
        </button>
      </div>
    </form>
  )
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns='http://www.w3.org/2000/svg'
      width='20'
      height='20'
      viewBox='0 0 24 24'
      fill='none'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
      className={className}
    >
      <circle cx='11' cy='11' r='8' />
      <line x1='21' y1='21' x2='16.65' y2='16.65' />
    </svg>
  )
}
