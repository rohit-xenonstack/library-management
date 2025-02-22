import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { bookAPI } from "../../utils/api";
import styles from "./Admin.module.scss";

export const AddBook = () => {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({
    isbn: "",
    title: "",
    authors: "",
    publisher: "",
    version: "",
  });

  const addBookMutation = useMutation({
    mutationFn: bookAPI.addBook,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books"] });
      setFormData({
        isbn: "",
        title: "",
        authors: "",
        publisher: "",
        version: "",
      });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    addBookMutation.mutate(formData);
  };

  return (
    <div className={styles.adminContainer}>
      <form onSubmit={handleSubmit} className={styles.adminForm}>
        <h2>Add New Book</h2>
        {addBookMutation.error && (
          <div className={styles.error}>
            {(addBookMutation.error as Error).message}
          </div>
        )}
        <input
          type="text"
          placeholder="ISBN"
          value={formData.isbn}
          onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
          required
        />
        <input
          type="text"
          placeholder="Title"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          required
        />
        <input
          type="text"
          placeholder="Authors"
          value={formData.authors}
          onChange={(e) =>
            setFormData({ ...formData, authors: e.target.value })
          }
          required
        />
        <input
          type="text"
          placeholder="Publisher"
          value={formData.authors}
          onChange={(e) =>
            setFormData({ ...formData, publisher: e.target.value })
          }
          required
        />
        <input
          type="text"
          placeholder="Version"
          value={formData.authors}
          onChange={(e) =>
            setFormData({ ...formData, version: e.target.value })
          }
          required
        />
        <button type="submit" disabled={addBookMutation.isPending}>
          {addBookMutation.isPending ? "Adding..." : "Add Book"}
        </button>
      </form>
    </div>
  );
};
