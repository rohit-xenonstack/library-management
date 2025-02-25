import { describe, it } from 'vitest'
import { App } from './app'
import { render, screen } from '@testing-library/react'

describe('Main entry point', () => {
  it('should render the app', () => {
    render(<App />)
    screen.debug()
  })
})
