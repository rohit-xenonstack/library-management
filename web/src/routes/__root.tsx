import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

import { Navigation } from '../components/navbar'
import type { RouterContext } from '../lib/router'

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
})

function RootComponent() {
  return (
    <>
      <Navigation />

      <Outlet />

      <TanStackRouterDevtools position='bottom-right' />
    </>
  )
}
