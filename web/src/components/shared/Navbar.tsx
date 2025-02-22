import { useState } from "react";
import { Link } from "@tanstack/react-router";
import { useAuth } from "../../hooks/useAuth";
import styles from "./Navbar.module.scss";

export const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { logout } = useAuth();
  const user = JSON.parse(localStorage.getItem("user") || "{}");

  const menuItems = {
    reader: [
      { to: "/", label: "Browse Books" },
      { to: "/reader/my-requests", label: "My Requests" },
    ],
    admin: [
      { to: "/", label: "Books" },
      { to: "/admin/add-book", label: "Add Book" },
      { to: "/admin/issue-requests", label: "Issue Requests" },
    ],
    owner: [
      { to: "/", label: "Books" },
      { to: "/owner/create-library", label: "Create Library" },
      { to: "/owner/create-admin", label: "Create Admin" },
    ],
  };

  const handleLogout: React.MouseEventHandler<HTMLButtonElement> = (event) => {
    event.preventDefault();
    logout();
  };

  const roleBasedMenu = user.role
    ? menuItems[user.role as keyof typeof menuItems]
    : [];

  return (
    <nav className={styles.navbar}>
      <div className={styles.navContent}>
        <Link to="/" className={styles.logo}>
          Library Management
        </Link>

        <button
          className={`${styles.menuButton} ${isOpen ? styles.active : ""}`}
          onClick={() => setIsOpen(!isOpen)}
          aria-label="Toggle menu"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>

        <div className={`${styles.navLinks} ${isOpen ? styles.active : ""}`}>
          {!user.id ? (
            <>
              <Link to="/login">Login</Link>
              {/* <Link to="/signup">Register</Link> */}
            </>
          ) : (
            <>
              {roleBasedMenu.map((item) => (
                <Link key={item.to} to={item.to}>
                  {item.label}
                </Link>
              ))}
              <div className={styles.userInfo}>
                <span>{user.name}</span>
                <button onClick={handleLogout} className={styles.logoutButton}>
                  Logout
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};
