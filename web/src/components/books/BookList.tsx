import { useQuery, useQueryClient } from "@tanstack/react-query";
import { bookAPI } from "../../utils/api";
import styles from "./Books.module.scss";
import { useState } from "react";

interface SearchParams {
  searchString: string;
  field: "title" | "authors" | "publisher";
}

export const BookList = () => {
  const queryClient = useQueryClient();
  const [searchParams, setSearchParams] = useState<SearchParams>({
    searchString: "",
    field: "title",
  });

  const { data: books, isLoading } = useQuery({
    queryKey: ["books", searchParams],
    queryFn: () =>
      searchParams.searchString
        ? bookAPI.searchBooks(searchParams)
        : bookAPI.getAllBooks(),
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // Trigger query refetch
    queryClient.invalidateQueries({ queryKey: ["books"] });
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.loader}></div>
        <p>Loading books...</p>
      </div>
    );
  }

  return (
    <div className={styles.bookListContainer}>
      <form onSubmit={handleSearch} className={styles.searchSection}>
        <div className={styles.searchInputs}>
          <input
            type="text"
            placeholder="Search books..."
            value={searchParams.searchString}
            onChange={(e) =>
              setSearchParams((prev) => ({
                ...prev,
                searchString: e.target.value,
              }))
            }
          />
          <select
            value={searchParams.field}
            onChange={(e) =>
              setSearchParams((prev) => ({
                ...prev,
                field: e.target.value as SearchParams["field"],
              }))
            }
          >
            <option value="title">Title</option>
            <option value="authors">Author</option>
            <option value="publisher">Publisher</option>
          </select>
          <button type="submit">Search</button>
        </div>
      </form>

      <div className={styles.bookGrid}>
        {books?.data.map((book: any) => (
          <div key={book.id} className={styles.bookCard}>
            <div className={styles.bookInfo}>
              <h3>{book.title}</h3>
              <p>
                <strong>Author:</strong> {book.authors}
              </p>
              <p>
                <strong>Publisher:</strong> {book.publisher}
              </p>
              <p>
                <strong>ISBN:</strong> {book.isbn}
              </p>
              <p className={styles.quantity}>
                <strong>Available:</strong>
                <span
                  className={
                    book.quantity > 0 ? styles.inStock : styles.outOfStock
                  }
                >
                  {book.quantity > 0
                    ? `${book.quantity} copies`
                    : "Out of stock"}
                </span>
              </p>
            </div>
            <button
              className={styles.issueButton}
              onClick={() => handleIssueRequest(book.id)}
              disabled={book.quantity === 0}
            >
              {book.quantity === 0 ? "Not Available" : "Request Book"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
