import { RouterProvider } from '@tanstack/react-router'

import { AuthProvider } from './context/auth-provider'
import { useAuth } from './hook/use-auth'
import { router } from './lib/router'

function App() {
  return (
    <AuthProvider>
      <RouterProviderWithContext />
    </AuthProvider>
  )
}

export const RouterProviderForTest = () => {
  return <RouterProvider router={router} />
}

function RouterProviderWithContext() {
  const auth = useAuth()

  return <RouterProvider router={router} context={{ auth }} />
}

export { App }
