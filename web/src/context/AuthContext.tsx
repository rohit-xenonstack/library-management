import {
  useState,
  useEffect,
  ReactNode,
  useCallback,
  createContext,
  useContext,
} from "react";
import { useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { authAPI } from "../utils/api";
// import { AuthContext } from "./authContextDefinition";

interface AuthProviderProps {
  children: ReactNode;
}

export interface User {
  id: string;
  name: string;
  email: string;
  contact: string;
  role: "owner" | "admin" | "reader";
  libraryId?: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (data: { email: string }) => Promise<void>;
  logout: () => Promise<void>;
}

// Create context with undefined default value
export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const navigate = useNavigate();
  const [token, setTokenState] = useState<string | null>(() =>
    localStorage.getItem("token"),
  );

  const { data: user, refetch: refetchUser } = useQuery({
    queryKey: ["currentUser"],
    queryFn: authAPI.getCurrentUser,
    enabled: !!token,
    retry: false,
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
  });

  const login = async ({ email }: { email: string }) => {
    const response = await authAPI.login({ email });
    localStorage.setItem("token", response.data.token);
    setTokenState(response.data.token);
    refetchUser(); // Fetch user details when token is set
  };

  const logout = useCallback(async () => {
    localStorage.removeItem("token");
    setTokenState(null);
    navigate({ to: "/login" });
  }, [navigate]);

  // Handle token expiration or invalid token
  useEffect(() => {
    if (token) {
      logout();
    }
  }, [token, logout]);

  return (
    <AuthContext.Provider
      value={{
        user: user?.data || null,
        token,
        isAuthenticated: !!token && !!user?.data,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
};
