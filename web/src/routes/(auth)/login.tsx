import { createFileRoute, redirect, useRouter } from "@tanstack/react-router";
import { useAuth } from "../../hooks/useAuth";
import styles from "./auth.module.scss";
import { set, z } from "zod";
import { useState } from "react";

// eslint-disable-next-line @typescript-eslint/no-unnecessary-type-assertion
const fallback = "/" as const;

export const Route = createFileRoute("/(auth)/login")({
  validateSearch: z.object({
    redirect: z.string().optional().catch(""),
  }),
  beforeLoad: ({ context, search }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: search.redirect || fallback });
    }
  },
  component: LoginComponent,
});

function LoginComponent() {
  const { login, isLoading } = useAuth();
  const router = useRouter();
  const navigate = Route.useNavigate();
  const search = Route.useSearch();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    setLoading(true);
    try {
      e.preventDefault();
      const data = new FormData(e.currentTarget);
      const fieldValue = data.get("username");

      if (!fieldValue) return;
      const email = fieldValue.toString();
      await login({ email: email });
      await router.invalidate();
      await navigate({ to: search.redirect || fallback });
    } catch (error) {
      console.error("Error logging in: ", error);
    }
    setLoading(false);
  };

  return (
    <div className={styles.authContainer}>
      <form onSubmit={handleSubmit} className={styles.authForm}>
        <h2>Login</h2>
        {error && <div className={styles.error}>{error.message}</div>}
        <input
          id="email"
          name="email"
          type="email"
          placeholder="Enter your email"
          required
        />
        <button type="submit" disabled={isLoading}>
          {isLoading ? "Loading..." : "Login"}
        </button>
      </form>
    </div>
  );
}
