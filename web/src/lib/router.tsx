import { createRouter } from '@tanstack/react-router'

import { routeTree } from '../route-tree.gen'
import type { AuthContextType } from '../context/auth-context'


type RouterContext = {
  auth: AuthContextType
}

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
  context: {
    auth: null as unknown as AuthContextType,
  },
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

export { router }
export type { RouterContext }
