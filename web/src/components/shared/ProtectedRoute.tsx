import { Navigate, useNavigate } from "@tanstack/react-router";
import { useAuth } from "../../hooks/useAuth";
import { ReactNode, useEffect } from "react";

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRole?: "admin" | "owner" | "reader";
}

export const ProtectedRoute = ({
  children,
  requiredRole,
}: ProtectedRouteProps) => {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const { isAuthenticated, user } = useAuth();
  // const user = JSON.parse(localStorage.getItem('user') || '{}')

  useEffect(() => {
    if (!isAuthenticated) {
      navigate({ to: "/login" });
    }

    if (requiredRole && user && user.role !== requiredRole) {
      navigate({ to: "/" });
    }
  }, [token, user, requiredRole, navigate]);

  if (!token) {
    return <Navigate to="/login" />;
  }

  if (requiredRole && user && user.role !== requiredRole) {
    return <Navigate to="/" />;
  }

  return <>{children}</>;
};
