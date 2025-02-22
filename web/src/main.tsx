import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import {
  NotFoundRoute,
  RouterProvider,
  createRouter,
} from "@tanstack/react-router";
import { Route as rootRoute } from "./routes/__root.tsx";

// Import the generated route tree
import { routeTree } from "./routeTree.gen";
import { AuthProvider } from "./context/AuthContext.tsx";
import { useAuth } from "./hooks/useAuth.ts";
import { QueryClient } from "@tanstack/react-query";

const notFoundRoute = new NotFoundRoute({
  getParentRoute: () => rootRoute,
  component: () => "404 Not Found",
});
const queryClient = new QueryClient();
// Create a new router instance
const router = createRouter({
  routeTree,
  notFoundRoute,
  defaultPreload: "intent",
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
  context: {
    auth: undefined!,
    queryClient,
  },
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  // const auth = useAuth();
  const queryClient = new QueryClient();
  return <RouterProvider router={router} context={{ queryClient }} />;
}

function App() {
  return (
    <AuthProvider>
      <InnerApp />
    </AuthProvider>
  );
}

// Render the app
const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <App />
    </StrictMode>,
  );
}
