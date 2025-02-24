import '@testing-library/jest-dom'

import { fireEvent, render, screen } from '@testing-library/react'

import { useAuth } from '../../hook/use-auth'
import { router } from '../../lib/router'

import { Navigation } from '../navbar'

// Mock the dependencies
jest.mock('~/hook/use-auth')
jest.mock('~/lib/router')
jest.mock('@tanstack/react-router', () => ({
  Link: ({ children, ...props }: React.PropsWithChildren<React.AnchorHTMLAttributes<HTMLAnchorElement>>) => <a {...props}>{children}</a>,
}))

describe('Navigation Component', () => {
  const mockedUseAuth = useAuth as jest.Mock
  const mockNavigate = jest.fn()

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks()

    // Mock router navigate
    router.navigate = mockNavigate

    // Default auth mock
    mockedUseAuth.mockReturnValue({
      user: null,
      logout: jest.fn(),
    })
  })

  test('renders navigation for unauthenticated user', () => {
    render(<Navigation />)

    // Check if logo is present
    expect(screen.getByText('Library App')).toBeInTheDocument()

    // Check if Sign In button is present
    expect(screen.getByText('Sign In')).toBeInTheDocument()

    // Check that admin-only links are not present
    expect(screen.queryByText('Add Book')).not.toBeInTheDocument()
    expect(screen.queryByText('Issue Requests')).not.toBeInTheDocument()
  })

  test('renders navigation for admin user', () => {
    mockedUseAuth.mockReturnValue({
      user: { role: 'admin' },
      logout: jest.fn(),
    })

    render(<Navigation />)

    // Check admin-specific links
    expect(screen.getByText('Add Book')).toBeInTheDocument()
    expect(screen.getByText('Issue Requests')).toBeInTheDocument()
    expect(screen.getByText('Sign Out')).toBeInTheDocument()
  })

  test('renders navigation for owner user', () => {
    mockedUseAuth.mockReturnValue({
      user: { role: 'owner' },
      logout: jest.fn(),
    })

    render(<Navigation />)

    // Check owner-specific links
    expect(screen.getByText('Onboard Admin')).toBeInTheDocument()
    expect(screen.getByText('Create Library')).toBeInTheDocument()
  })

  test('handles sign out correctly', () => {
    const mockLogout = jest.fn()
    mockedUseAuth.mockReturnValue({
      user: { role: 'admin' },
      logout: mockLogout,
    })

    render(<Navigation />)

    const signOutButton = screen.getByText('Sign Out')
    fireEvent.click(signOutButton)

    expect(mockLogout).toHaveBeenCalled()
    expect(mockNavigate).toHaveBeenCalledWith({ to: '/' })
  })

  test('handles sign in navigation', () => {
    render(<Navigation />)

    const signInButton = screen.getByText('Sign In')
    fireEvent.click(signInButton)

    expect(mockNavigate).toHaveBeenCalledWith({ to: '/sign-in' })
  })

  test('handles mobile menu toggle', () => {
    render(<Navigation />)

    const hamburgerButton = screen.getByLabelText('Toggle navigation')
    const nav = screen.getByRole('navigation')

    // Initial state
    expect(nav).not.toHaveClass('open')

    // Open menu
    fireEvent.click(hamburgerButton)
    expect(nav).toHaveClass('open')

    // Close menu
    fireEvent.click(hamburgerButton)
    expect(nav).not.toHaveClass('open')
  })

  test('closes menu when clicking backdrop', () => {
    render(<Navigation />)

    const hamburgerButton = screen.getByLabelText('Toggle navigation')
    const backdrop = screen.getByLabelText('Close menu')
    const nav = screen.getByRole('navigation')

    // Open menu
    fireEvent.click(hamburgerButton)
    expect(nav).toHaveClass('open')

    // Click backdrop
    fireEvent.click(backdrop)
    expect(nav).not.toHaveClass('open')
  })
})
