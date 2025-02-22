import { createRouter, NotFoundRoute } from "@tanstack/react-router";

import { Route as rootRoute } from "../routes/__root";
import { QueryClient } from "@tanstack/react-query";
import { routeTree } from "../routeTree.gen";

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
